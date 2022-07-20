package rcalendar.client;

import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.annotation.RegisterClientHeaders;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;
import org.jboss.resteasy.annotations.jaxrs.QueryParam;
import rcalendar.model.GetMeetingRoomReservationResponse;
import rcalendar.model.ReserveMeetingRoomRequest;
import rcalendar.model.UpdateMeetingRoomRequest;

import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;

@Path("/room")
@RegisterRestClient
@RegisterClientHeaders(OrecocoTokenAuthorizationFactory.class)
public interface OrecocoReserveClient {
    @GET
    Uni<GetMeetingRoomReservationResponse> getRooms(@QueryParam("scheduleId") String scheduleId);

    @POST
    Uni<Void> reserveRoom(ReserveMeetingRoomRequest request);

    @PUT
    Uni<Void> updateMeetingRoom(UpdateMeetingRoomRequest request);
}
