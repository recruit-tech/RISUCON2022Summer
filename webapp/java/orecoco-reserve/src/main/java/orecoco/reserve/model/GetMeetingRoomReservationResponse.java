package orecoco.reserve.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record GetMeetingRoomReservationResponse(
        @JsonProperty("room_id") String roomId
) {

}
