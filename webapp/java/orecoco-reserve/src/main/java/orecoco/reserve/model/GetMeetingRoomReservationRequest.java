package orecoco.reserve.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record GetMeetingRoomReservationRequest(
        @JsonProperty("schedule_id") String scheduleId
) {
}
