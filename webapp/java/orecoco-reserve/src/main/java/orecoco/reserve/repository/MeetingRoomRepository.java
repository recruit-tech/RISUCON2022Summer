package orecoco.reserve.repository;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.mysqlclient.MySQLPool;
import io.vertx.mutiny.sqlclient.Row;
import io.vertx.mutiny.sqlclient.RowIterator;
import io.vertx.mutiny.sqlclient.Tuple;
import orecoco.reserve.model.MeetingRoom;

import javax.inject.Singleton;
import java.time.LocalDateTime;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

@Singleton
public class MeetingRoomRepository {
    private final MySQLPool client;
    private final MeetingRoomMapper meetingRoomMapper;

    public MeetingRoomRepository(MySQLPool client, MeetingRoomMapper meetingRoomMapper) {
        this.client = client;
        this.meetingRoomMapper = meetingRoomMapper;
    }

    public Uni<List<MeetingRoom>> find(String roomId, LocalDateTime startAt, LocalDateTime endAt) {
        return client.preparedQuery("SELECT * FROM meeting_room WHERE room_id = ? AND NOT (start_at >= ? OR end_at <= ?)")
                .execute(Tuple.of(roomId, endAt, startAt))
                .map(rows -> StreamSupport.stream(rows.spliterator(), false)
                        .map(meetingRoomMapper::rowToModel)
                        .collect(Collectors.toList()));
    }

    public Uni<List<MeetingRoom>> findForUpdate(String id, String roomId, LocalDateTime startAt, LocalDateTime endAt) {
        return client.preparedQuery("SELECT * FROM meeting_room WHERE room_id = ? AND NOT (start_at >= ? OR end_at <= ?) AND id != ? FOR UPDATE")
                .execute(Tuple.of(roomId, endAt, startAt, id))
                .map(rows -> StreamSupport.stream(rows.spliterator(), false)
                        .map(meetingRoomMapper::rowToModel)
                        .collect(Collectors.toList()));
    }

    public Uni<MeetingRoom> findById(String id) {
        return client.preparedQuery("SELECT * FROM meeting_room WHERE id = ?")
                .execute(Tuple.of(id))
                .map(rows -> {
                    RowIterator<Row> iterator = rows.iterator();
                    return iterator.hasNext() ?
                            meetingRoomMapper.rowToModel(iterator.next())
                            :
                            null;
                });
    }

    public Uni<Void> create(MeetingRoom room) {
        return client.preparedQuery("INSERT INTO meeting_room (id, room_id, start_at, end_at) VALUES (?, ?, ?, ?)")
                .execute(Tuple.of(room.id(), room.roomId(), room.startAt(), room.endAt()))
                .chain(rows -> Uni.createFrom().voidItem());
    }

    public Uni<Void> delete(String id) {
        return client.preparedQuery("DELETE FROM meeting_room WHERE id = ?")
                .execute(Tuple.of(id))
                .chain(rows -> Uni.createFrom().voidItem());
    }
}
