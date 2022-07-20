package rcalendar.repository;

import io.vertx.mutiny.sqlclient.Row;
import rcalendar.model.Schedule;

import javax.inject.Singleton;

@Singleton
public class ScheduleMapper {
    public Schedule rowToModel(Row row) {
        return new Schedule(row.getString("id"),
                row.getString("title"),
                row.getString("description"),
                row.getString("schedule_attendee"),
                row.getLocalDateTime("start_at"),
                row.getLocalDateTime("end_at"));
    }}
