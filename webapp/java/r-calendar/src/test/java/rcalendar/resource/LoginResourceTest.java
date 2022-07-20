package rcalendar.resource;

import io.smallrye.jwt.algorithm.SignatureAlgorithm;
import io.smallrye.jwt.auth.principal.DefaultJWTParser;
import io.smallrye.jwt.auth.principal.ParseException;
import io.smallrye.jwt.build.Jwt;
import org.jose4j.keys.HmacKey;
import org.junit.jupiter.api.Test;

import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;

class LoginResourceTest {
    @Test
    void tokenGenerate() throws ParseException {
        String secret = "secret_keysecret_keysecret_keysecret_key";
        SecretKeySpec key = new HmacKey(secret.getBytes(StandardCharsets.UTF_8));

        String jwt = Jwt.claim("jti", "1234")
                .subject("jojo")
                .jws()
                .algorithm(SignatureAlgorithm.HS256)
                .sign(key);

        System.out.println(new DefaultJWTParser()
                        .verify(jwt, key));
    }
}