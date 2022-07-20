package rcalendar.resource;

import io.quarkus.narayana.jta.QuarkusTransaction;
import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.Vertx;
import org.jboss.logging.Logger;
import org.jboss.resteasy.reactive.MultipartForm;
import rcalendar.middleware.SHA256;
import rcalendar.model.GetUserResponse;
import rcalendar.model.UpdateIconRequest;
import rcalendar.model.UpdateUserRequest;
import rcalendar.repository.UserRepository;

import javax.ws.rs.*;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.SecurityContext;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Map;
import java.util.Optional;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static javax.ws.rs.core.Response.Status.*;
import static rcalendar.resource.MimeTypeDetector.detectMimeType;

@Produces(MediaType.APPLICATION_JSON)
@Path("/me")
public class MeResource {
    private final static Logger LOG = Logger.getLogger(MeResource.class);

    private final UserRepository userRepository;
    private final Vertx vertx;

    public MeResource(UserRepository userRepository, Vertx vertx) {
        this.userRepository = userRepository;
        this.vertx = vertx;
    }

    @GET
    public Uni<Response> getMe(@Context SecurityContext ctx) {
        String userId = ctx.getUserPrincipal().getName();

        return userRepository.findById(userId)
                .chain(user -> user != null ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new UserNotFoundException(userId)))
                .map(user -> Response.ok()
                        .entity(new GetUserResponse(
                                user.id(),
                                user.email(),
                                user.name(),
                                Optional.ofNullable(user.imageBinary())
                                        .map(b -> "/icon/" + user.id())
                                        .orElse(null)
                        ))
                        .build())
                .onFailure(UserNotFoundException.class)
                .recoverWithItem(() -> Response.status(400).build());
    }

    @PUT
    @Consumes(MediaType.APPLICATION_JSON)
    public Uni<Response> updateMe(@Context SecurityContext ctx,
                                  UpdateUserRequest request) {
        String userId = ctx.getUserPrincipal().getName();
        return userRepository.findById(userId)
                .chain(user -> user != null ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new UserNotFoundException(userId)))
                .chain(user -> Uni.createFrom().item(request))
                .chain(r ->
                    QuarkusTransaction.call(() -> userRepository.update(userId, r.name(), r.email(),
                            SHA256.hash(r.password())))
                            .replaceWithVoid()
                )
                .map(v -> Response.ok().build());
    }

    @PUT
    @Path("/icon")
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    public Uni<Response> updateMyIcon(@Context SecurityContext ctx,
                                      @MultipartForm UpdateIconRequest request) {
       String userId = ctx.getUserPrincipal().getName();
        return userRepository.findById(userId)
                .chain(user -> user != null ?
                        Uni.createFrom().item(user)
                        :
                        Uni.createFrom().failure(new UserNotFoundException(userId))
                )
                .chain(u -> vertx.fileSystem().createTempFile("icon", ""))
                .chain(fname -> vertx.fileSystem().readFile(request.getIcon().uploadedFile().toString())
                        .chain(buf -> vertx.fileSystem().writeFile(fname, buf))
                        .replaceWith(fname))
                .chain(fname -> detectMimeType(fname) != null ?
                        Uni.createFrom().item(fname)
                        :
                        Uni.createFrom().failure(new NotAllowedImageTypeException(fname)))
                .chain(fname -> vertx.fileSystem().readFile(fname)
                        .chain(buf -> vertx.fileSystem().delete(fname)
                                .replaceWith(buf))
                )
                .chain(buf -> QuarkusTransaction.call(() -> userRepository.updateIcon(userId, buf)))
                .map(v-> Response.ok().build())
                .onFailure(UserNotFoundException.class).recoverWithItem(() ->
                        Response.status(NOT_FOUND)
                                .entity("ユーザーが見つかりませんでした")
                                .build())
                .onFailure(NotAllowedImageTypeException.class).recoverWithItem(() ->
                        Response.status(BAD_REQUEST)
                                .entity("アイコンに指定できるのはjpeg, png, gifまたはbmpの画像ファイルのみです")
                                .build())
                .onFailure().recoverWithItem(ex -> {
                    LOG.error("Server Error", ex);
                    return Response.status(INTERNAL_SERVER_ERROR)
                            .entity("サーバ側のエラーです")
                            .build();
                });
    }

}
