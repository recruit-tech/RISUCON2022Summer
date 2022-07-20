package orecoco.reserve.resource;

public class MeetingRoomNotFoundException extends RuntimeException {
    public MeetingRoomNotFoundException(String meetingRoomId) {
        super(meetingRoomId);
    }
}
