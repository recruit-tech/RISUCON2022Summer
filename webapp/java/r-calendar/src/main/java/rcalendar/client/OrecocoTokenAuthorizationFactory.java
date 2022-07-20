package rcalendar.client;

import org.eclipse.microprofile.rest.client.ext.ClientHeadersFactory;

import javax.inject.Singleton;
import javax.ws.rs.core.MultivaluedHashMap;
import javax.ws.rs.core.MultivaluedMap;
import java.util.Optional;

@Singleton
public class OrecocoTokenAuthorizationFactory implements ClientHeadersFactory {
    private final OrecocoReserveTokenHolder orecocoReserveTokenHolder;

    public OrecocoTokenAuthorizationFactory(OrecocoReserveTokenHolder orecocoReserveTokenHolder) {
        this.orecocoReserveTokenHolder = orecocoReserveTokenHolder;
    }

    @Override
    public MultivaluedMap<String, String> update(MultivaluedMap<String, String> incoming, MultivaluedMap<String, String> outgoing) {
        MultivaluedMap<String, String> result = new MultivaluedHashMap<>();
        Optional.ofNullable(orecocoReserveTokenHolder.get())
                .ifPresent(token -> result.add("Authorization", "Bearer " + token));
        return result;
    }
}
