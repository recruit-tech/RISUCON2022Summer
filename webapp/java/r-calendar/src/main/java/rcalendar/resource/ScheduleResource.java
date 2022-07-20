package rcalendar.resource;

import am.ik.yavi.core.ConstraintViolations;
import am.ik.yavi.core.ConstraintViolationsException;
import com.github.f4b6a3.ulid.UlidCreator;
import io.quarkus.narayana.jta.QuarkusTransaction;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RestClient;
import org.jboss.logging.Logger;

import org.jboss.resteasy.reactive.RestPath;
import rcalendar.client.OrecocoReserveClient;
import rcalendar.model.*;
import rcalendar.repository.ScheduleRepository;
import rcalendar.repository.UserRepository;

import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.time.LocalDateTime;
import java.time.ZoneOffset;

import static javax.ws.rs.core.Response.Status.BAD_REQUEST;
import static javax.ws.rs.core.Response.Status.CREATED;

@Produces(MediaType.APPLICATION_JSON)
@Path("/schedule")
public class ScheduleResource {
    private static final Logger LOG = Logger.getLogger(ScheduleResource.class);

    private final UserRepository userRepository;
    private final ScheduleRepository scheduleRepository;
    private final OrecocoReserveClient orecocoReserveClient;

    public ScheduleResource(UserRepository userRepository, ScheduleRepository scheduleRepository,
                            @RestClient OrecocoReserveClient orecocoReserveClient) {
        this.userRepository = userRepository;
        this.scheduleRepository = scheduleRepository;
        this.orecocoReserveClient = orecocoReserveClient;
    }

    @Consumes(MediaType.APPLICATION_JSON)
    @POST
    public Uni<Response> createNewSchedule(CreateScheduleRequest request) {
        return Uni.createFrom().item(request)
                .chain(r -> {
                    ConstraintViolations violations = r.validate();
                    return violations.isValid() ?
                            Uni.createFrom().item(r)
                            :
                            Uni.createFrom().failure(new ConstraintViolationsException(violations.violations()));
                })
                .chain(r -> Multi.createFrom().items(r.attendees().stream())
                        .onItem()
                        .transformToUni(attendeeId -> userRepository.findById(attendeeId)
                                .chain(user -> user == null ?
                                        Uni.createFrom().failure(new UserNotFoundException(attendeeId))
                                        :
                                        Uni.createFrom().item(user)))
                        .concatenate()
                        .onFailure(UserNotFoundException.class).invoke(t -> Response.status(BAD_REQUEST)
                                .build())
                        .collect()
                        .asMultiMap(User::id)
                        .chain(userMap ->
                                userMap.entrySet().stream().anyMatch(entry -> entry.getValue().size() > 1) ?
                                        Uni.createFrom().failure(new DuplicateAttendeeException(""))
                                        :
                                        Uni.createFrom().item(userMap.keySet())
                        )
                )
                .map(attendees -> new Schedule(
                        UlidCreator.getUlid().toString(),
                        request.title(),
                        request.description(),
                        String.join(",", attendees),
                        LocalDateTime.ofEpochSecond(request.startAt(), 0, ZoneOffset.UTC),
                        LocalDateTime.ofEpochSecond(request.endAt(), 0, ZoneOffset.UTC)
                ))
                .chain(schedule-> QuarkusTransaction.call(() -> scheduleRepository.create(schedule))
                        .replaceWith(schedule))
                .chain(schedule-> {
                    if (request.meetingRoom() != null && !request.meetingRoom().isEmpty()) {
                        return orecocoReserveClient.reserveRoom(new ReserveMeetingRoomRequest(
                                        schedule.id(),
                                        request.meetingRoom(),
                                        request.startAt(),
                                        request.endAt()))
                                .replaceWith(schedule);
                    } else {
                        return Uni.createFrom().item(schedule);
                    }
                })
                .map(schedule -> Response.status(CREATED)
                        .entity(new CreateScheduleResponse(schedule.id()))
                        .build())
                .onFailure(ConstraintViolationsException.class).recoverWithItem(t -> Response.status(BAD_REQUEST)
                        .entity(t.getMessage())
                        .build());
    }

    @Path("/{scheduleId}")
    @GET
    public Uni<Response> getSchedule(@RestPath("scheduleId") String scheduleId) {
        var scheduleResponseComposer = new ScheduleResponseComposer(userRepository, orecocoReserveClient);
        return scheduleRepository.findById(scheduleId)
                .chain(schedule -> schedule != null ?
                        Uni.createFrom().item(schedule)
                        :
                        Uni.createFrom().failure(new ScheduleNotFoundException(scheduleId)))
                .chain(scheduleResponseComposer::compose)
                .map(scheduleResponse -> Response.ok()
                        .entity(scheduleResponse)
                        .build());
    }

    @PUT
    @Consumes(MediaType.APPLICATION_JSON)
    public Uni<Response> updateScheduleWithoutScheduleId(UpdateScheduleRequest request) {
        return Uni.createFrom().item(Response.status(BAD_REQUEST).build());
    }

    @Path("/{scheduleId}")
    @Consumes(MediaType.APPLICATION_JSON)
    @PUT
    public Uni<Response> updateSchedule(@RestPath("scheduleId") String scheduleId,
                                        UpdateScheduleRequest request) {
        return Uni.createFrom().item(request)
                .chain(r -> {
                    ConstraintViolations violations = r.validate();
                    return violations.isValid() ?
                            Uni.createFrom().item(r)
                            :
                            Uni.createFrom().failure(new ConstraintViolationsException(violations.violations()));
                })
                .chain(r -> Multi.createFrom().items(r.attendees().stream())
                        .onItem()
                        .transformToUni(attendeeId -> userRepository.findById(attendeeId)
                                .chain(user -> user == null ?
                                        Uni.createFrom().failure(new UserNotFoundException(attendeeId))
                                        :
                                        Uni.createFrom().item(user)))
                        .concatenate()
                        .onFailure(UserNotFoundException.class).invoke(t -> Response.status(BAD_REQUEST).build())
                        .collect()
                        .asMultiMap(User::id)
                        .chain(userMap ->
                                userMap.entrySet().stream().anyMatch(entry -> entry.getValue().size() > 1) ?
                                        Uni.createFrom().failure(new DuplicateAttendeeException(""))
                                        :
                                        Uni.createFrom().item(userMap.keySet())
                        )
                )
                .map(attendees -> new Schedule(
                        scheduleId,
                        request.title(),
                        request.description(),
                        String.join(",", attendees),
                        LocalDateTime.ofEpochSecond(request.startAt(), 0, ZoneOffset.UTC),
                        LocalDateTime.ofEpochSecond(request.endAt(), 0, ZoneOffset.UTC)
                ))
                .chain(schedule-> scheduleRepository.update(schedule)
                        .replaceWith(schedule))
                .chain(schedule-> orecocoReserveClient.updateMeetingRoom(new UpdateMeetingRoomRequest(
                            schedule.id(),
                            request.meetingRoom(),
                            request.startAt(),
                            request.endAt()))
                        .replaceWith(schedule)
                )
                .map(schedule -> Response.ok()
                        .build())
                .onFailure(ConstraintViolationsException.class).invoke(t -> Response.status(BAD_REQUEST)
                        .build());
    }
}
