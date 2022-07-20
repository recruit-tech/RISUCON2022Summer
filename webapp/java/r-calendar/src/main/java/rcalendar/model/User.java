package rcalendar.model;

import java.time.LocalDateTime;

public record User(
    String id,
    String email,
    String name,
    String password,
    byte[] imageBinary,
    LocalDateTime createdAt
) {
}
