package rcalendar.resource;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.Vertx;
import io.vertx.mutiny.core.buffer.Buffer;
import io.vertx.mutiny.mysqlclient.MySQLPool;
import org.eclipse.microprofile.rest.client.inject.RestClient;
import rcalendar.client.OrecocoInitializeClient;
import rcalendar.client.OrecocoReserveTokenHolder;
import rcalendar.model.InitializeResponse;

import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

@Produces(MediaType.APPLICATION_JSON)
@Path("/initialize")
public class InitializeResource {
    private final MySQLPool client;
    private final Vertx vertx;
    private final OrecocoInitializeClient orecocoInitializeClient;
    private final OrecocoReserveTokenHolder orecocoReserveTokenHolder;

    public InitializeResource(MySQLPool client, Vertx vertx,
                              @RestClient OrecocoInitializeClient orecocoInitializeClient,
                              OrecocoReserveTokenHolder orecocoReserveTokenHolder) {
        this.client = client;
        this.vertx = vertx;
        this.orecocoInitializeClient = orecocoInitializeClient;
        this.orecocoReserveTokenHolder = orecocoReserveTokenHolder;
    }

    @POST
    @Path("")
    public Uni<Response> initialize() {
        Uni<String> ddlSql = vertx.fileSystem().readFile("../../sql/r-calendar-0_Schema.sql")
                .map(Buffer::toString);
        Uni<String> dataSql = vertx.fileSystem().readFile("../../sql/r-calendar-1_DummyUserData.sql")
                .map(Buffer::toString);

        return ddlSql
                .chain(sql -> client.query(sql).execute())
                .replaceWith(dataSql)
                .chain(sql -> client.query(sql).execute())
                .chain($ -> orecocoInitializeClient.initialize())
                .map(response -> orecocoReserveTokenHolder.set(response.token()))
                .replaceWith(() -> Response.ok().entity(new InitializeResponse("java")).build())
                .onFailure()
                .recoverWithItem(ex -> {
                    ex.printStackTrace();
                    return Response.status(500).entity("サーバー側のエラーです").build();
                });
    }

}
