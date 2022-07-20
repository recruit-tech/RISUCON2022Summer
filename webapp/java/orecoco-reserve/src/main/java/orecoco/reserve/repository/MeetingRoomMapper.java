package orecoco.reserve.repository;

import io.vertx.mutiny.sqlclient.Row;
import orecoco.reserve.model.MeetingRoom;

import javax.inject.Singleton;

@Singleton
public class MeetingRoomMapper {
    public MeetingRoom rowToModel(Row row) {
        return new MeetingRoom(row.getString("id"),
                            row.getString("room_id"),
                            row.getLocalDateTime("start_at"),
                            row.getLocalDateTime("end_at"));
    }
}
