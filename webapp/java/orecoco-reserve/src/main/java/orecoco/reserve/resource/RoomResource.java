package orecoco.reserve.resource;

import io.quarkus.narayana.jta.QuarkusTransaction;
import io.smallrye.mutiny.Uni;
import orecoco.reserve.model.GetMeetingRoomReservationResponse;
import orecoco.reserve.model.MeetingRoom;
import orecoco.reserve.model.ReserveMeetingRoomRequest;
import orecoco.reserve.model.UpdateMeetingRoomRequest;
import orecoco.reserve.repository.MeetingRoomRepository;
import orecoco.reserve.repository.MeetingRoomsRepository;
import org.jboss.logging.Logger;
import org.jboss.resteasy.reactive.RestQuery;

import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.time.LocalDateTime;
import java.time.ZoneOffset;

import static javax.ws.rs.core.Response.Status.*;

@Produces(MediaType.APPLICATION_JSON)
@Path("/room")
public class RoomResource {
    private static final Logger LOG = Logger.getLogger(RoomResource.class);

    private final MeetingRoomsRepository meetingRoomsRepository;
    private final MeetingRoomRepository meetingRoomRepository;

    public RoomResource(MeetingRoomsRepository meetingRoomsRepository, MeetingRoomRepository meetingRoomRepository) {
        this.meetingRoomsRepository = meetingRoomsRepository;
        this.meetingRoomRepository = meetingRoomRepository;
    }

    @Consumes(MediaType.APPLICATION_JSON)
    @POST
    public Uni<Response> reserveMeetingRoom(ReserveMeetingRoomRequest request) {
        return meetingRoomsRepository.meetingRoomExist(request.meetingRoomId())
                .chain(exists -> exists ?
                        Uni.createFrom().voidItem()
                        :
                        Uni.createFrom().failure(new MeetingRoomNotFoundException(request.meetingRoomId()))
                )
                .chain($ -> meetingRoomRepository.find(request.meetingRoomId(),
                        LocalDateTime.ofEpochSecond(request.startAt(), 0, ZoneOffset.UTC),
                        LocalDateTime.ofEpochSecond(request.endAt(), 0, ZoneOffset.UTC)))
                .chain(meetingRooms -> meetingRooms.size() > 0 ?
                        Uni.createFrom().failure(new AlreadyReservedException())
                        :
                        Uni.createFrom().item(new MeetingRoom(request.scheduleId(),
                                request.meetingRoomId(),
                                LocalDateTime.ofEpochSecond(request.startAt(), 0, ZoneOffset.UTC),
                                LocalDateTime.ofEpochSecond(request.endAt(), 0, ZoneOffset.UTC)))
                )
                .chain(meetingRoom -> QuarkusTransaction.call(() -> meetingRoomRepository.create(meetingRoom)))
                .map($ -> Response.status(CREATED).build())
                .onFailure(MeetingRoomNotFoundException.class).recoverWithItem(() -> Response.status(Response.Status.BAD_REQUEST)
                        .entity("会議室が見つかりません")
                        .build())
                .onFailure(AlreadyReservedException.class).recoverWithItem(ex -> Response.status(CONFLICT)
                        .entity("すでに予定が入っています")
                        .build())
                .onFailure().recoverWithItem(() -> Response
                        .status(INTERNAL_SERVER_ERROR)
                        .entity("ーバー側のエラーです")
                        .build());
    }

    @Consumes(MediaType.APPLICATION_JSON)
    @GET
    public Uni<Response> getMeetingRoomReservation(@RestQuery("scheduleId") String scheduleId) {
        return Uni.createFrom().item(scheduleId)
                .chain(id -> id == null ?
                        Uni.createFrom().failure(new MissingScheduleIdException())
                        :
                        Uni.createFrom().item(id))
                .chain(meetingRoomRepository::findById)
                .chain(room -> room == null ?
                        Uni.createFrom().failure(new ScheduleNotFoundException(scheduleId))
                        :
                        Uni.createFrom().item(new GetMeetingRoomReservationResponse(room.roomId())))
                .map(room -> Response.ok(room).build())
                .onFailure(MissingScheduleIdException.class).recoverWithItem(() -> Response
                        .status(BAD_REQUEST)
                        .entity("スケジュールを指定してください")
                        .build())
                .onFailure(ScheduleNotFoundException.class).recoverWithItem(() -> Response
                        .status(NOT_FOUND)
                        .entity("指定したスケジュールには部屋が予約されていません")
                        .build())
                .onFailure().recoverWithItem(() -> Response
                        .status(INTERNAL_SERVER_ERROR)
                        .entity("サーバー側のエラーです")
                        .build());
    }

    @Consumes(MediaType.APPLICATION_JSON)
    @PUT
    public Uni<Response> updateMeetingRoom(UpdateMeetingRoomRequest request) {
        if (request.meetingRoomId() == null || request.meetingRoomId().isEmpty()) {
            return QuarkusTransaction.call(() -> meetingRoomRepository.delete(request.scheduleId()))
                    .map($ -> Response.ok().build())
                    .onFailure().recoverWithItem(() -> Response
                            .status(INTERNAL_SERVER_ERROR)
                            .entity("サーバー側のエラーです")
                            .build());
        }
        return meetingRoomsRepository.meetingRoomExist(request.meetingRoomId())
                .chain(exists -> exists ?
                        Uni.createFrom().voidItem()
                        :
                        Uni.createFrom().failure(new MeetingRoomNotFoundException(request.meetingRoomId()))
                )
                .chain($ -> QuarkusTransaction.call(() -> meetingRoomRepository.findForUpdate(request.scheduleId(), request.meetingRoomId(),
                                LocalDateTime.ofEpochSecond(request.startAt(), 0, ZoneOffset.UTC),
                                LocalDateTime.ofEpochSecond(request.endAt(), 0, ZoneOffset.UTC)
                                )
                                .chain(meetingRooms -> meetingRooms.size() > 0 ?
                                        Uni.createFrom().failure(new AlreadyReservedException())
                                        :
                                        Uni.createFrom().item(new MeetingRoom(request.scheduleId(),
                                                request.meetingRoomId(),
                                                LocalDateTime.ofEpochSecond(request.startAt(), 0, ZoneOffset.UTC),
                                                LocalDateTime.ofEpochSecond(request.endAt(), 0, ZoneOffset.UTC)))
                                )
                                .chain(meetingRoom -> meetingRoomRepository.delete(request.scheduleId())
                                        .replaceWith(meetingRoom))
                                .chain(meetingRoomRepository::create)
                        )
                )
                .map($ -> Response.status(OK).build())
                .onFailure(MeetingRoomNotFoundException.class).recoverWithItem(() -> Response
                        .status(BAD_REQUEST)
                        .entity("会議室が見つかりません")
                        .build())
                .onFailure(AlreadyReservedException.class).recoverWithItem(ex -> Response
                        .status(CONFLICT)
                        .entity("すでに予定が入っています")
                        .build())
                .onFailure().recoverWithItem(() -> Response
                        .status(INTERNAL_SERVER_ERROR)
                        .entity("サーバー側のエラーです")
                        .build());
    }
}
