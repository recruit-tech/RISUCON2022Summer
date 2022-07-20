package rcalendar.resource;

import io.smallrye.mutiny.Uni;
import rcalendar.filter.JwtSecurityContext;
import rcalendar.middleware.JWTAuthenticator;
import rcalendar.middleware.SHA256;
import rcalendar.middleware.SessionManager;
import rcalendar.model.LoginRequest;
import rcalendar.repository.UserRepository;

import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.*;
import java.util.HashMap;
import java.util.Objects;
import java.util.Optional;

import static javax.ws.rs.core.Response.Status.BAD_REQUEST;

@Path("")
public class LoginResource {
    private final UserRepository userRepository;
    private final SessionManager sessionManager;
    private final JWTAuthenticator jwtAuthenticator;

    public LoginResource(UserRepository userRepository, SessionManager sessionManager, JWTAuthenticator jwtAuthenticator) {
        this.userRepository = userRepository;
        this.sessionManager = sessionManager;
        this.jwtAuthenticator = jwtAuthenticator;
    }

    @Produces(MediaType.APPLICATION_JSON)
    @Consumes(MediaType.APPLICATION_JSON)
    @POST
    @Path("/login")
    public Uni<Response> postLogin(LoginRequest request) {
        return userRepository.findByEmail(request.email())
                .chain(user -> user != null ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new UserNotFoundException(request.email()))
                )
                .chain(user -> Objects.equals(user.password(), SHA256.hash(request.password())) ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new PasswordMismatchException(request.email()))
                )
                .map(user -> {
                    String sessionId = sessionManager.setNewSession(new HashMap<>());
                    return jwtAuthenticator.generateTokenString(user.id(), sessionId);
                })
                .map(token -> Response.status(201)
                        .cookie(new NewCookie("jwt", token, "/", null,
                                null,
                                24 * 60 * 60,
                                false,
                                true))
                        //.header("Set-Cookie", "jwt="+token+";Path=/;Max-Age=86400;HttpOnly;Secure;SameSite=None")
                        .build()
                )
                .onFailure(PasswordMismatchException.class).recoverWithItem(() -> Response.status(BAD_REQUEST)
                        .entity("ユーザー名またはパスワードが不正です")
                        .build())
                .onFailure(UserNotFoundException.class).recoverWithItem(() -> Response.status(BAD_REQUEST)
                        .entity("ユーザー名またはパスワードが不正です")
                        .build());

    }

    @POST
    @Path("/logout")
    public Uni<Response> postLogout(@Context SecurityContext ctx) {
        String sessionKey = Optional.of(ctx)
                .map(JwtSecurityContext.class::cast)
                .map(JwtSecurityContext::getSessionId)
                .orElseThrow(() -> new IllegalArgumentException("Invalid session id"));
        sessionManager.delete(sessionKey);
        return Uni.createFrom().item(() -> Response.status(Response.Status.CREATED).build());
    }
}
