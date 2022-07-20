package rcalendar.resource;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ConflictEmailException extends RuntimeException{
    public ConflictEmailException(@JsonProperty("email") String email) {
        super(email);
    }
}
