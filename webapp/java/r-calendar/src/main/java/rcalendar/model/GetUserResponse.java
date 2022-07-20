package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record GetUserResponse(
        @JsonProperty("id") String id,
        @JsonProperty("email") String email,
        @JsonProperty("name") String name,
        @JsonProperty("icon") String icon
) {
}
