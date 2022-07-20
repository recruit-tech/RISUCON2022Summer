package rcalendar.repository;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.buffer.Buffer;
import io.vertx.mutiny.mysqlclient.MySQLPool;
import io.vertx.mutiny.sqlclient.Row;
import io.vertx.mutiny.sqlclient.RowIterator;
import io.vertx.mutiny.sqlclient.Tuple;
import org.jboss.logging.Logger;
import org.jose4j.base64url.Base64;
import rcalendar.model.User;

import javax.inject.Singleton;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

@Singleton
public class UserRepository {
    private static final Logger LOG = Logger.getLogger(UserRepository.class);

    private final MySQLPool client;
    private final UserMapper userMapper;

    public UserRepository(MySQLPool client, UserMapper userMapper) {
        this.client = client;
        this.userMapper = userMapper;
    }

    public Uni<User> findById(String id) {
        return client.preparedQuery("SELECT * FROM user WHERE id = ?")
                .execute(Tuple.of(id))
                .map(rows -> {
                    RowIterator<Row> iterator = rows.iterator();
                    return iterator.hasNext() ?
                            userMapper.rowToModel(iterator.next())
                            :
                            null;
                });
    }

    public Uni<User> findByEmail(String email) {
        return client.preparedQuery("SELECT * FROM user WHERE email = ?")
                .execute(Tuple.of(email))
                .map(rows -> {
                    RowIterator<Row> iterator = rows.iterator();
                    return iterator.hasNext() ?
                            userMapper.rowToModel(iterator.next())
                            :
                            null;
                });
    }

    public Uni<List<User>> search(String query) {
        String q = query + "%";
        // ユーザーが作成された順にソートして返す
        return client.preparedQuery("SELECT * FROM user WHERE email LIKE ? OR name LIKE ? ORDER BY id")
                .execute(Tuple.of(q, q))
                .map(rows -> StreamSupport.stream(rows.spliterator(), false)
                        .map(userMapper::rowToModel)
                        .collect(Collectors.toList()));
    }

    public Uni<Void> create(User user) {
        return client.preparedQuery("INSERT INTO user(id, email, name, password) VALUES (?, ?, ?, ?)")
                .execute(Tuple.of(user.id(), user.email(), user.name(), user.password()))
                .replaceWithVoid();
    }

    public Uni<Void> update(String id, String name, String email, String password) {
        return client.preparedQuery("UPDATE user SET name = ?, email = ?, password = ? WHERE id = ?")
                .execute(Tuple.of(name, email, password, id))
                .replaceWithVoid();
    }

    public Uni<Void> updateIcon(String id, Buffer imageBinary) {
        return client.preparedQuery("UPDATE user SET image_binary = ? WHERE id = ?")
                .execute(Tuple.of(Base64.encode(imageBinary.getBytes()), id))
                .replaceWithVoid();
    }
}
