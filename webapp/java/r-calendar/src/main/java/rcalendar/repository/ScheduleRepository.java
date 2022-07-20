package rcalendar.repository;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.mysqlclient.MySQLPool;
import io.vertx.mutiny.sqlclient.Row;
import io.vertx.mutiny.sqlclient.RowIterator;
import io.vertx.mutiny.sqlclient.Tuple;
import org.jboss.logging.Logger;
import rcalendar.model.Schedule;

import javax.inject.Singleton;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

@Singleton
public class ScheduleRepository {
    private static final Logger LOG = Logger.getLogger(ScheduleRepository.class);

    private final MySQLPool client;
    private final ScheduleMapper scheduleMapper;

    public ScheduleRepository(MySQLPool client, ScheduleMapper scheduleMapper) {
        this.client = client;
        this.scheduleMapper = scheduleMapper;
    }

    public Uni<Void> create(Schedule schedule) {
        return client.preparedQuery("INSERT INTO schedule (id, title, description, schedule_attendee, start_at, end_at) VALUES (?, ?, ?, ?, ?, ?)")
                .execute(Tuple.of(schedule.id(), schedule.title(), schedule.description(),
                        schedule.scheduleAttendee(),
                        schedule.startAt(),
                        schedule.endAt()
                ))
                .replaceWithVoid();
    }

    public Uni<Void> update(Schedule schedule) {
        return client.preparedQuery("UPDATE schedule SET title = ?, description = ?, schedule_attendee = ?, start_at = ?, end_at = ? WHERE id = ?")
                .execute(Tuple.of(schedule.title(), schedule.description(),
                        schedule.scheduleAttendee(),
                        schedule.startAt(),
                        schedule.endAt(),
                        schedule.id()
                ))
                .replaceWithVoid();
    }

    public Uni<List<Schedule>> findByAttendee(String userId) {
        return client.preparedQuery("SELECT * FROM schedule WHERE schedule_attendee LIKE ?")
                .execute(Tuple.of("%" + userId + "%"))
                .map(rows -> StreamSupport.stream(rows.spliterator(), false)
                        .map(scheduleMapper::rowToModel)
                        .collect(Collectors.toList()));
    }

    public Uni<Schedule> findById(String scheduleId) {
        return client.preparedQuery("SELECT * FROM schedule WHERE id = ?")
                .execute(Tuple.of(scheduleId))
                .map(rows -> {
                    RowIterator<Row> iterator = rows.iterator();
                    return iterator.hasNext() ?
                            scheduleMapper.rowToModel(iterator.next())
                            :
                            null;
                });
    }

    public Uni<List<Schedule>> findByIds(List<String> ids) {
        return client.preparedQuery("SELECT * FROM schedule WHERE id IN ("+ "?,".repeat(ids.size()).substring(0,ids.size()*2-1) +") ORDER BY start_at ASC, end_at DESC, id ASC")
                .execute(Tuple.wrap(ids))
                .map(rows -> StreamSupport.stream(rows.spliterator(), false)
                        .map(scheduleMapper::rowToModel)
                        .collect(Collectors.toList()));
    }
}
