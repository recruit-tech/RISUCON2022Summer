package rcalendar.middleware;

import java.io.Serializable;

public interface SessionManager {
    <T extends Serializable> String setNewSession(T sessionData);

    boolean existSession(String sessionId);

    <T extends Serializable> T get(String sessionId, Class<T> clazz);

    void delete(String sessionId);
}
