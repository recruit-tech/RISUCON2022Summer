package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record GetScheduleResponse(
        @JsonProperty("id") String id,
        @JsonProperty("title") String title,
        @JsonProperty("description") String description,
        @JsonProperty("start_at") Long startAt,
        @JsonProperty("end_at") Long endAt,
        @JsonProperty("attendees") List<GetUserResponse> attendees,
        @JsonProperty("meeting_room") String meetingRoom
) {
}
