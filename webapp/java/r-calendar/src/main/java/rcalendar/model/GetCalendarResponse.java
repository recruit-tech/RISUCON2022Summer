package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record GetCalendarResponse(
    @JsonProperty("date") Long date,
    @JsonProperty("schedules") List<GetScheduleResponse> schedules
) {
}
