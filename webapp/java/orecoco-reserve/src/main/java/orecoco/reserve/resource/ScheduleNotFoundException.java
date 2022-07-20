package orecoco.reserve.resource;

public class ScheduleNotFoundException extends RuntimeException{
    public ScheduleNotFoundException(String scheduleId) {
        super(scheduleId);
    }
}
