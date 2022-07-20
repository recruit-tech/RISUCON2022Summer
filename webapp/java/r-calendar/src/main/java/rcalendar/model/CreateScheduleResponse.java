package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record CreateScheduleResponse(
        @JsonProperty("id") String id
) {
}
