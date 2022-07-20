package rcalendar.resource;

import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import rcalendar.client.OrecocoReserveClient;
import rcalendar.model.GetMeetingRoomReservationResponse;
import rcalendar.model.GetScheduleResponse;
import rcalendar.model.GetUserResponse;
import rcalendar.model.Schedule;
import rcalendar.repository.UserRepository;

import java.time.ZoneOffset;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

public class ScheduleResponseComposer {
    private final UserRepository userRepository;
    private final OrecocoReserveClient orecocoReserveClient;

    public ScheduleResponseComposer(UserRepository userRepository, OrecocoReserveClient orecocoReserveClient) {
        this.userRepository = userRepository;
        this.orecocoReserveClient = orecocoReserveClient;
    }

    public Uni<GetScheduleResponse> compose(Schedule schedule) {
        String[] scheduleAttendees = schedule.scheduleAttendee().split(",");
        Uni<List<GetUserResponse>> userResponseUni = Multi.createFrom().items(scheduleAttendees)
                .onItem()
                .transformToUni(userRepository::findById)
                .concatenate()
                .map(attendeeUser -> new GetUserResponse(
                        attendeeUser.id(),
                        attendeeUser.email(),
                        attendeeUser.name(),
                        attendeeUser.imageBinary() != null ? "/icon/" + attendeeUser.id() : "")
                )
                .collect()
                .asList()
                .map(response -> response.stream().sorted(Comparator.comparing(GetUserResponse::id)).collect(Collectors.toList()));

        Uni<GetMeetingRoomReservationResponse> reservationUni = orecocoReserveClient.getRooms(schedule.id())
                .onFailure()
                .recoverWithItem(new GetMeetingRoomReservationResponse("")); // FIXME 404のときだけ実行したい

        return Uni.combine().all()
                .unis(userResponseUni, reservationUni)
                .combinedWith((attendees, reservation) -> new GetScheduleResponse(
                        schedule.id(),
                        schedule.title(),
                        schedule.description(),
                        schedule.startAt().toEpochSecond(ZoneOffset.UTC),
                        schedule.endAt().toEpochSecond(ZoneOffset.UTC),
                        attendees,
                        reservation.roomId()
                ));
    }
}
