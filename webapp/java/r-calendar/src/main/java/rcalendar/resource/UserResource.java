package rcalendar.resource;

import am.ik.yavi.core.ConstraintViolationsException;
import com.github.f4b6a3.ulid.UlidCreator;
import io.quarkus.narayana.jta.QuarkusTransaction;
import io.smallrye.mutiny.Uni;
import io.smallrye.mutiny.tuples.Tuple;
import io.smallrye.mutiny.tuples.Tuple2;
import io.vertx.mutiny.core.Vertx;
import io.vertx.mutiny.core.buffer.Buffer;
import org.jboss.logging.Logger;
import org.jboss.resteasy.reactive.RestPath;
import org.jboss.resteasy.reactive.RestQuery;
import rcalendar.middleware.JWTAuthenticator;
import rcalendar.middleware.SHA256;
import rcalendar.middleware.SessionManager;
import rcalendar.model.CreateUserRequest;
import rcalendar.model.GetUserResponse;
import rcalendar.model.SearchUserResponse;
import rcalendar.model.User;
import rcalendar.repository.UserRepository;

import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.NewCookie;
import javax.ws.rs.core.Response;
import java.util.HashMap;
import java.util.Objects;
import java.util.Optional;
import java.util.stream.Collectors;

import static javax.ws.rs.core.Response.Status.*;
import static rcalendar.resource.MimeTypeDetector.detectMimeType;

@Produces(MediaType.APPLICATION_JSON)
@Path("/user")
public class UserResource {
    private static final Logger LOG = Logger.getLogger(UserResource.class);

    private final UserRepository userRepository;
    private final SessionManager sessionManager;
    private final JWTAuthenticator jwtAuthenticator;
    private final Vertx vertx;

    public UserResource(UserRepository userRepository, SessionManager sessionManager, JWTAuthenticator jwtAuthenticator, Vertx vertx) {
        this.userRepository = userRepository;
        this.sessionManager = sessionManager;
        this.jwtAuthenticator = jwtAuthenticator;
        this.vertx = vertx;
    }

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    public Uni<Response> createUser(CreateUserRequest request) {
        return Uni.createFrom().item(request.validate())
                .chain(violations -> violations.isValid() ?
                        Uni.createFrom().item(request)
                        :
                        Uni.createFrom().failure(new ConstraintViolationsException(violations)))
                .chain(r -> userRepository.findByEmail(r.email()))
                .chain(user -> Objects.nonNull(user) ?
                        Uni.createFrom().failure(new ConflictEmailException(request.email()))
                        :
                        Uni.createFrom().voidItem())
                .replaceWith(() -> new User(UlidCreator.getUlid().toString(),
                        request.email(),
                        request.name(),
                        SHA256.hash(request.password()),
                        null,
                        null
                        ))
                .chain(user -> QuarkusTransaction.call(() -> userRepository.create(user)
                        .replaceWith(user)))
                .map(user -> {
                    String sessionId = sessionManager.setNewSession(new HashMap<>());
                    return jwtAuthenticator.generateTokenString(user.id(), sessionId);
                }).map(token-> Response.status(CREATED)
                        .cookie(new NewCookie("jwt", token, "/", null, null,
                                24 * 60 * 60,
                                false,
                                true))
                        .build()
                )
                .onFailure(ConstraintViolationsException.class).recoverWithItem(ex -> Response.status(BAD_REQUEST)
                        .entity(ex.getMessage())
                        .build())
                .onFailure(ConflictEmailException.class).recoverWithItem(() -> Response.status(BAD_REQUEST)
                        .entity("指定されたメールアドレスはすでに利用されています")
                        .build())
                .onFailure().recoverWithItem(() -> Response.status(INTERNAL_SERVER_ERROR)
                    .entity("サーバー側のエラーです")
                    .build());
    }

    @GET
    public Uni<Response> searchUser(@RestQuery("query") String query) {
        if (query == null || query.isEmpty()) {
            return Uni.createFrom().item(Response.status(BAD_REQUEST)
                    .entity("検索条件を指定してください")
                    .build());
        }
        return userRepository.search(query)
                .chain(users -> users.isEmpty() ?
                        Uni.createFrom().failure(new UserNotFoundException(""))
                        :
                        Uni.createFrom().item(users))
                .map(users -> Response.ok()
                        .entity(new SearchUserResponse(users.stream().map(u -> new GetUserResponse(
                                u.id(),
                                u.email(),
                                u.name(),
                                Optional.ofNullable(u.imageBinary())
                                        .map($-> "/icon/" + u.id())
                                        .orElse("")
                                ))
                                .collect(Collectors.toList())))
                        .build())
                .onFailure(UserNotFoundException.class).recoverWithItem(() ->
                        Response.noContent()
                                .entity("ユーザーが見つかりませんでした")
                                .build())
                .onFailure().recoverWithItem(() ->
                        Response.status(INTERNAL_SERVER_ERROR)
                                .entity("サーバー側のエラーです")
                                .build());
    }

    @GET
    @Path("/{userId}")
    public Uni<Response> getUser(@RestPath("userId") String userId) {
        if (userId == null || userId.isEmpty()) {
            return Uni.createFrom().item(Response.status(BAD_REQUEST)
                    .entity("検索条件を指定してください")
                    .build());
        }
        return userRepository.findById(userId)
                .chain(user -> user != null ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new UserNotFoundException(userId))
                )
                .map(user -> Response.ok()
                        .entity(new GetUserResponse(user.id(), user.email(), user.name(),
                                Optional.ofNullable(user.imageBinary())
                                        .map($ -> "/icon/"+userId)
                                        .orElse("")))
                        .build()
                )
                .onFailure(UserNotFoundException.class).recoverWithItem(() -> Response.status(NOT_FOUND)
                        .entity("ユーザーが見つかりませんでした")
                        .build());
    }

    @GET
    @Path("/icon/{userId}")
    public Uni<Response> getUserIcon(@RestPath("userId") String userId) {
        return userRepository.findById(userId)
                .chain(user -> user != null ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new UserNotFoundException(userId))
                )
                .chain(user -> user.imageBinary() != null ?
                        Uni.createFrom().item(user.imageBinary())
                        :
                        Uni.createFrom().failure(new UserNotFoundException(userId))
                )
                .chain(buf -> vertx.fileSystem().createTempFile("icon", "")
                        .chain(fname -> vertx.fileSystem().writeFile(fname, Buffer.buffer(buf))
                                .replaceWith(fname))
                        .chain(fname -> {
                            String mimeType = detectMimeType(fname);
                            if (mimeType == null) throw new IllegalStateException("不明なMimeTypeです");
                            return vertx.fileSystem().delete(fname).replaceWith(mimeType);
                        })
                        .map(mimeType -> Tuple2.of(mimeType, buf))
                )
                .map(tpl -> Response.ok()
                        .header("Content-Type", tpl.getItem1())
                        .entity(tpl.getItem2())
                        .build()
                )
                .onFailure(UserNotFoundException.class).recoverWithItem(() -> Response.status(NOT_FOUND)
                        .entity("ユーザーが見つかりませんでした")
                        .build());
    }
}
