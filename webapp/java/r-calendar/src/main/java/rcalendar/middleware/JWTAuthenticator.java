package rcalendar.middleware;

import io.smallrye.jwt.auth.principal.JWTParser;
import io.smallrye.jwt.auth.principal.ParseException;
import io.smallrye.jwt.build.Jwt;
import org.eclipse.microprofile.jwt.JsonWebToken;
import org.jose4j.keys.HmacKey;

import javax.crypto.SecretKey;
import javax.inject.Singleton;
import java.nio.charset.StandardCharsets;
import java.time.Duration;

@Singleton
public class JWTAuthenticator {
    private final JWTParser jwtParser;

    private final SecretKey secretKey;
    public JWTAuthenticator(JWTParser jwtParser) {
        this.jwtParser = jwtParser;
        byte[] key = new byte[256];
        System.arraycopy("secret_key".getBytes(StandardCharsets.UTF_8), 0, key, 0, 10);
        secretKey = new HmacKey(key);
    }

    public String generateTokenString(String userId, String jwtId) {
        return Jwt.claim("user_id", userId)
                .claim("jti", jwtId)
                .subject(jwtId)
                .expiresIn(Duration.ofDays(1))
                .jws()
                .sign(secretKey);
    }

    public JsonWebToken verify(String token) throws ParseException {
        return jwtParser.verify(token, secretKey);
    }
}
