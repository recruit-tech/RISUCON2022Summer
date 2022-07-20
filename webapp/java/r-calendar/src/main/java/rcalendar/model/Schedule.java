package rcalendar.model;

import java.time.LocalDateTime;

public record Schedule(
        String id,
        String title,
        String description,
        String scheduleAttendee,
        LocalDateTime startAt,
        LocalDateTime endAt
) {
}
