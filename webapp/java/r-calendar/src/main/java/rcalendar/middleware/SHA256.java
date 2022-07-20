package rcalendar.middleware;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.HexFormat;

public class SHA256 {
    public static String hash(String text) {
        try {
            MessageDigest md = MessageDigest.getInstance("SHA-256");
            byte[] b = md.digest(text.getBytes(StandardCharsets.UTF_8));
            return HexFormat.of().formatHex(b);
        } catch (NoSuchAlgorithmException e) {
            throw new IllegalStateException("発生しちゃダメなやつ");
        }
    }
}
