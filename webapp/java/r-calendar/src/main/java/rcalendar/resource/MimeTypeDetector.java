package rcalendar.resource;

import org.jboss.logging.Logger;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.util.Map;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

public class MimeTypeDetector {
    private static final Logger LOG = Logger.getLogger(MimeTypeDetector.class);

    private static final Map<String, String> ALLOWED_FILE_TYPES = Map.of(
            "image/jpeg","JPEG image data",
            "image/png","PNG image data",
            "image/gif","GIF image data",
            "image/bmp","PB bitmap"
    );

    public static String detectMimeType(String fname) {
        return ALLOWED_FILE_TYPES.entrySet()
                .stream()
                .filter(entry -> {
                    Pattern reg = Pattern.compile(entry.getValue());
                    String out = execFileCommand(fname);
                    return reg.matcher(out).find();
                })
                .findAny()
                .map(Map.Entry::getKey)
                .orElse(null);
    }

    private static String execFileCommand(String fileName) {
        try {
            ProcessBuilder pb = new ProcessBuilder("file", fileName);
            pb.redirectErrorStream(true);
            Process process = pb.start();
            String out = process.inputReader().lines().collect(Collectors.joining(System.lineSeparator()));
            process.waitFor();
            return out;
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        } catch (InterruptedException e) {
            throw new IllegalStateException(e);
        }
    }
}
