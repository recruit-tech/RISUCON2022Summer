package rcalendar.resource;

import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RestClient;
import org.jboss.logging.Logger;
import org.jboss.resteasy.reactive.RestPath;
import org.jboss.resteasy.reactive.RestQuery;
import rcalendar.client.OrecocoReserveClient;
import rcalendar.model.*;
import rcalendar.repository.ScheduleRepository;
import rcalendar.repository.UserRepository;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import static javax.ws.rs.core.Response.Status.BAD_REQUEST;
import static javax.ws.rs.core.Response.Status.NOT_FOUND;

@Produces(MediaType.APPLICATION_JSON)
@Path("/calendar")
public class CalendarResource {
    private static final Logger LOG = Logger.getLogger(CalendarResource.class);

    private final UserRepository userRepository;
    private final ScheduleRepository scheduleRepository;
    private final OrecocoReserveClient orecocoReserveClient;

    public CalendarResource(UserRepository userRepository, ScheduleRepository scheduleRepository,
                            @RestClient OrecocoReserveClient orecocoReserveClient) {
        this.userRepository = userRepository;
        this.scheduleRepository = scheduleRepository;
        this.orecocoReserveClient = orecocoReserveClient;
    }

    @Path("/{userId}")
    @GET
    public Uni<Response> getCalendar(@RestPath("userId") String userId,
                                     @RestQuery("date") String dateFromOriginStr) {
        int dateFromOrigin;
        try {
            dateFromOrigin = Integer.parseInt(dateFromOriginStr);
        } catch (NumberFormatException e) {
            return Uni.createFrom().item(Response.status(BAD_REQUEST)
                    .entity("指定された日時は不正です")
                    .build());
        }
        var scheduleResponseComposer = new ScheduleResponseComposer(userRepository, orecocoReserveClient);
        return userRepository.findById(userId)
                .chain(user -> user == null ?
                        Uni.createFrom().failure(new UserNotFoundException(userId))
                        :
                        Uni.createFrom().item(user)
                )
                .chain(user -> scheduleRepository.findByAttendee(user.id())
                        .map(schedules -> filterParticipantSchedule(dateFromOrigin, schedules)))
                .chain(scheduleIds -> scheduleIds.size() == 0 ?
                        Uni.createFrom().failure(new NoParticipationScheduleException())
                        :
                        Uni.createFrom().item(scheduleIds))
                .chain(scheduleRepository::findByIds)
                .chain(schedules -> Multi.createFrom().items(schedules.stream())
                        .onItem()
                        .transformToUni(scheduleResponseComposer::compose)
                        .concatenate()
                        .collect()
                        .asList())
                .map(scheduleResponses -> Response.ok()
                        .entity(new GetCalendarResponse(
                                (long) dateFromOrigin,
                                scheduleResponses
                        ))
                        .build())
                .onFailure(UserNotFoundException.class).recoverWithItem(() ->
                        Response.status(NOT_FOUND)
                                .entity("ユーザーが見つかりませんでした")
                                .build())
                .onFailure(NoParticipationScheduleException.class).recoverWithItem(() ->
                        Response.ok()
                                .entity(new GetCalendarResponse(
                                        (long)dateFromOrigin,
                                        Collections.emptyList()
                                ))
                                .build());
    }

    private List<String> filterParticipantSchedule(int dateFromOrigin, List<Schedule> participantSchedules) {
        LocalDateTime dateRaw = LocalDateTime.ofEpochSecond((long) dateFromOrigin *24*60*60, 0, ZoneOffset.UTC);
        LocalDateTime startOfTheDay = LocalDateTime.of(dateRaw.getYear(), dateRaw.getMonth(), dateRaw.getDayOfMonth(), 0, 0, 0);
        LocalDateTime endOfTheDay = startOfTheDay.plusDays(1L);

        return participantSchedules.stream()
                .filter(schedule -> timeIn(schedule.startAt(), schedule.endAt(), startOfTheDay, endOfTheDay))
                .map(Schedule::id)
                .collect(Collectors.toList());
    }

    static boolean timeIn(LocalDateTime start, LocalDateTime end,
                           LocalDateTime targetStartAt, LocalDateTime targetEndAt) {
        return !(targetEndAt.isEqual(start) || targetEndAt.isBefore(start) || targetStartAt.isAfter(end));
    }
}
