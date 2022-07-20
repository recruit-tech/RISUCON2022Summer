package rcalendar.filter;

import javax.ws.rs.core.SecurityContext;
import java.security.Principal;

public class JwtSecurityContext implements SecurityContext {
    private final String sessionId;
    private final String userId;

    public JwtSecurityContext(String userId, String sessionId) {
        this.userId = userId;
        this.sessionId = sessionId;
    }

    public String getSessionId() {
        return sessionId;
    }

    @Override
    public Principal getUserPrincipal() {
        return () -> userId;
    }

    @Override
    public boolean isUserInRole(String role) {
        return false;
    }

    @Override
    public boolean isSecure() {
        return false;
    }

    @Override
    public String getAuthenticationScheme() {
        return null;
    }
}
