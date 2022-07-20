package rcalendar.middleware;

import javax.inject.Singleton;
import java.io.Serializable;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

@Singleton
public class InMemorySessionManager implements SessionManager {
    private final ConcurrentMap<String, Serializable> sessions;
    private final SessionIdGenerator sessionIdGenerator;

    public InMemorySessionManager(SessionIdGenerator sessionIdGenerator) {
        this.sessionIdGenerator = sessionIdGenerator;
        this.sessions = new ConcurrentHashMap<>();
    }

    @Override
    public <T extends Serializable> String setNewSession(T sessionData) {
        String sessionId = sessionIdGenerator.generate();
        sessions.put(sessionId, sessionData);
        return sessionId;
    }

    @Override
    public boolean existSession(String sessionId) {
        return sessions.containsKey(sessionId);
    }

    @Override
    public <T extends Serializable> T get(String sessionId, Class<T> clazz) {
        return clazz.cast(sessions.get(sessionId));
    }

    @Override
    public void delete(String sessionId) {
        sessions.remove(sessionId);
    }
}
