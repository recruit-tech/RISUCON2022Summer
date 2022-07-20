package orecoco.reserve.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.LocalDateTime;

public record MeetingRoom(
    @JsonProperty("id") String id,
    @JsonProperty("room_id") String roomId,
    @JsonProperty("start_at") LocalDateTime startAt,
    @JsonProperty("end_at") LocalDateTime endAt
) {
}
