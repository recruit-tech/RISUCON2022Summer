package rcalendar.resource;

public class ScheduleNotFoundException extends RuntimeException{
    public ScheduleNotFoundException(String scheduleId) {
        super(scheduleId);
    }
}
