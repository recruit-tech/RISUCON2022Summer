package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record SearchUserResponse(
        @JsonProperty("users")List<GetUserResponse> users
) {
}
