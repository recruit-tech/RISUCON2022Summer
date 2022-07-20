package orecoco.reserve.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record UpdateMeetingRoomRequest(
    @JsonProperty("schedule_id") String scheduleId,
    @JsonProperty("meeting_room_id") String meetingRoomId,
    @JsonProperty("start_at") Long startAt,
    @JsonProperty("end_at") Long endAt
) {
}
