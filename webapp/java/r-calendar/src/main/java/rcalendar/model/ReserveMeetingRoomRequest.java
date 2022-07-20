package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record ReserveMeetingRoomRequest(
        @JsonProperty("schedule_id") String scheduleId,
        @JsonProperty("meeting_room_id") String meetingRoomId,
        @JsonProperty("start_at") Long startAt,
        @JsonProperty("end_at") Long endAt
) {
}
