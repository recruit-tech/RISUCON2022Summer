package rcalendar.resource;

public class DuplicateAttendeeException extends RuntimeException{
    public DuplicateAttendeeException(String attendeeId) {
        super(attendeeId);
    }
}
