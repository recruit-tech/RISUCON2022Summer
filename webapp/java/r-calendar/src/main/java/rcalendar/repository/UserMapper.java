package rcalendar.repository;

import io.vertx.core.buffer.Buffer;
import io.vertx.mutiny.sqlclient.Row;
import org.jboss.logging.Logger;
import org.jose4j.base64url.Base64;
import rcalendar.model.User;

import javax.inject.Singleton;
import java.util.Optional;

@Singleton
public class UserMapper {
    private static final Logger LOG = Logger.getLogger(UserMapper.class);

    public User rowToModel(Row row) {
        return new User(row.getString("id"),
                row.getString("email"),
                row.getString("name"),
                row.getString("password"),
                Optional.ofNullable(row.getBuffer("image_binary"))
                        .map(buf -> buf.getString(0, buf.length()))
                        .map(Base64::decode)
                        .orElse(null),
                row.getLocalDateTime("created_at"));
    }
}
