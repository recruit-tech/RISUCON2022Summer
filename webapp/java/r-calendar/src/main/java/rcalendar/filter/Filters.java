package rcalendar.filter;

import io.smallrye.common.annotation.NonBlocking;
import io.smallrye.jwt.auth.principal.ParseException;
import io.smallrye.mutiny.Uni;
import org.jboss.logging.Logger;
import org.jboss.resteasy.reactive.RestResponse;
import org.jboss.resteasy.reactive.server.ServerRequestFilter;
import rcalendar.middleware.JWTAuthenticator;
import rcalendar.middleware.SessionManager;

import javax.inject.Singleton;
import javax.ws.rs.container.ContainerRequestContext;
import javax.ws.rs.core.Cookie;
import javax.ws.rs.core.Response;
import java.util.Objects;

import static org.jboss.resteasy.reactive.RestResponse.Status.BAD_REQUEST;
import static org.jboss.resteasy.reactive.RestResponse.Status.UNAUTHORIZED;

@Singleton
public class Filters {
    private static final Logger LOG = Logger.getLogger(Filters.class);
    private final JWTAuthenticator jwtAuthenticator;
    private final SessionManager sessionManager;

    public Filters(JWTAuthenticator jwtAuthenticator, SessionManager sessionManager) {
        this.jwtAuthenticator = jwtAuthenticator;
        this.sessionManager = sessionManager;
    }


    @ServerRequestFilter
    @NonBlocking
    public Uni<RestResponse<String>> getFilter(ContainerRequestContext ctx) {
        String path = ctx.getUriInfo().getPath();
        String method = ctx.getMethod();
        Cookie jwtCookie = ctx.getCookies().get("jwt");
        if (Objects.equals(path, "/login")
                || Objects.equals(path, "/initialize")
                || (Objects.equals(path, "/user") && Objects.equals(method, "POST"))
        ) {
            return Uni.createFrom().nullItem();
        }
		if (jwtCookie == null) {
            return Uni.createFrom()
                    .item(RestResponse.status(Response.Status.UNAUTHORIZED, "ログインしてください"));
		}

        return Uni.createFrom().voidItem()
                .chain(() -> {
                    try {
                        return Uni.createFrom().item(jwtAuthenticator.verify(jwtCookie.getValue()));
                    } catch (ParseException e) {
                        return Uni.createFrom().failure(e);
                    }
                })
                .chain(jwt -> sessionManager.existSession(jwt.getClaim("jti")) ?
                        Uni.createFrom().item(jwt)
                        :
                        Uni.createFrom().failure(new SessionNotFoundException())
                )
                .<RestResponse<String>>chain(jwt -> {
                    if (jwt != null && jwt.containsClaim("user_id")) {
                        ctx.setSecurityContext(new JwtSecurityContext(
                                jwt.getClaim("user_id"),
                                jwt.getClaim("jti")));
                        return Uni.createFrom().nullItem();
                    } else {
                        return Uni.createFrom().failure(new InvalidTokenException());
                    }
                })
                .onFailure(InvalidTokenException.class).recoverWithItem(() ->
                        RestResponse.status(BAD_REQUEST, "invalid token"))
                .onFailure(SessionNotFoundException.class).recoverWithItem(() ->
                        RestResponse.status(UNAUTHORIZED, "ログインしてください"));
    }
}
