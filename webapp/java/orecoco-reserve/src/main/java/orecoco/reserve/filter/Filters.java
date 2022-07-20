package orecoco.reserve.filter;

import orecoco.reserve.resource.TokenHolder;
import org.jboss.resteasy.reactive.RestResponse;
import org.jboss.resteasy.reactive.server.ServerRequestFilter;

import javax.inject.Singleton;
import javax.ws.rs.container.ContainerRequestContext;
import java.util.Objects;
import java.util.Optional;

import static javax.ws.rs.core.Response.Status.UNAUTHORIZED;

@Singleton
public class Filters {
    private final TokenHolder tokenHolder;

    public Filters(TokenHolder tokenHolder) {
        this.tokenHolder = tokenHolder;
    }

    @ServerRequestFilter
    public Optional<RestResponse<String>> getFilter(ContainerRequestContext ctx) {
        if (Objects.equals(ctx.getUriInfo().getPath(), "/initialize")) {
            return Optional.empty();
        }
		final String BEARER = "Bearer ";
        String tokenHeader = ctx.getHeaderString("Authorization");
		if (tokenHeader == null || !tokenHeader.startsWith(BEARER)) {
            return Optional.of(RestResponse.status(UNAUTHORIZED, "不正なアクセスです"));
		}
		String token = tokenHeader.substring(BEARER.length());
        if (!tokenHolder.get().matches(token)) {
            return Optional.of(RestResponse.status(UNAUTHORIZED, "不正なトークンです"));
        }
        return Optional.empty();
    }
}
