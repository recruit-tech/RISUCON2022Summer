package orecoco.reserve.repository;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.Vertx;

import javax.enterprise.context.ApplicationScoped;
import java.util.Arrays;

@ApplicationScoped
public class MeetingRoomsRepository {
    private final Vertx vertx;

    public MeetingRoomsRepository(Vertx vertx) {
        this.vertx = vertx;
    }

    public Uni<Boolean> meetingRoomExist(String name) {
        return vertx.fileSystem()
                .readFile("./meeting_room.txt")
                .map(buffer -> Arrays.asList(buffer.toString().split("\\n")).contains(name));
    }
}
