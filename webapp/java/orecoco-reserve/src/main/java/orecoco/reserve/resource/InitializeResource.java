package orecoco.reserve.resource;

import io.smallrye.context.api.CurrentThreadContext;
import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.Vertx;
import io.vertx.mutiny.core.buffer.Buffer;
import io.vertx.mutiny.mysqlclient.MySQLPool;
import orecoco.reserve.model.OrecocoToken;
import org.eclipse.microprofile.context.ThreadContext;

import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Produces(MediaType.APPLICATION_JSON)
@Path("/initialize")
public class InitializeResource {
    private final MySQLPool client;
    private final Vertx vertx;
    private final TokenHolder tokenHolder;

    public InitializeResource(MySQLPool client, Vertx vertx, TokenHolder tokenHolder) {
        this.client = client;
        this.vertx = vertx;
        this.tokenHolder = tokenHolder;
    }

    @POST
    @Path("")
    @CurrentThreadContext(propagated = {}, unchanged = ThreadContext.ALL_REMAINING)
    public Uni<OrecocoToken> initialize() {
        return vertx.fileSystem().readFile("../../sql/orecoco-reserve_0_Schema.sql")
                .map(Buffer::toString)
                .chain(initSql -> client.query(initSql).execute())
                .map(rows -> tokenHolder.set(new OrecocoToken()));
    }
}
