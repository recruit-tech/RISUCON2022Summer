package rcalendar.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record OrecocoReserveInitializeResponse(
        @JsonProperty("token") String token
) {
}
