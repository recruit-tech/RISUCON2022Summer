package rcalendar.client;

import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;
import rcalendar.model.OrecocoReserveInitializeResponse;

import javax.ws.rs.POST;
import javax.ws.rs.Path;

@Path("/initialize")
@RegisterRestClient
public interface OrecocoInitializeClient {
    @POST
    Uni<OrecocoReserveInitializeResponse> initialize();
}
