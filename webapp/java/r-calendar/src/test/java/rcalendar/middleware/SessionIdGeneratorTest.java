package rcalendar.middleware;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class SessionIdGeneratorTest {
    DefaultSessionIdGenerator sut;

    @BeforeEach
    void setup() {
        sut = new DefaultSessionIdGenerator();
    }
    @Test
    void test() {
        System.out.println(sut.generate());
    }
}