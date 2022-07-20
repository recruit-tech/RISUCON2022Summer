package rcalendar.resource;

public class UserNotFoundException extends RuntimeException {
    public UserNotFoundException(String attendeeId) {
        super(attendeeId);
    }
}
