package rcalendar.middleware;

import javax.inject.Singleton;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;

@Singleton
public class DefaultSessionIdGenerator implements SessionIdGenerator {
    private static final char[] LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".toCharArray();
    private final SecureRandom random;

    public DefaultSessionIdGenerator() {
        try {
            random = SecureRandom.getInstance("NativePRNG");
        } catch (NoSuchAlgorithmException e) {
            throw new IllegalStateException(e);
        }
    }

    @Override
    public String generate() {
        String[] possibleSessionIdList = new String[10000];

        for (int i = 0; i < 10000; i++) {
            byte[] bytes = new byte[128];
            random.nextBytes(bytes);
            StringBuilder sb = new StringBuilder(128);
            for (byte b : bytes) {
                sb.append(LETTERS[((int) b & 0xff) % LETTERS.length]);
            }
            possibleSessionIdList[i] = sb.toString();
        }
        return possibleSessionIdList[random.nextInt(LETTERS.length)];
    }
}
