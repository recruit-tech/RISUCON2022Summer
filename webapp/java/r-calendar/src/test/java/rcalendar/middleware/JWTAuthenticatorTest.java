package rcalendar.middleware;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import javax.inject.Inject;

@QuarkusTest
class JWTAuthenticatorTest {
    @Inject
    JWTAuthenticator jwtAuthenticator;

    @Test
    void testJWT() {
        String tokenString = jwtAuthenticator.generateTokenString("hoge", "1234");
        System.out.println(tokenString);
    }
}
