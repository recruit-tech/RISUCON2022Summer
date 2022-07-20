package rcalendar.resource;

public class PasswordMismatchException extends RuntimeException {
    public PasswordMismatchException(String email) {
        super(email);
    }
}
