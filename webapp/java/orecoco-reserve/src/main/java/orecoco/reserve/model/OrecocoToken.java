package orecoco.reserve.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.util.Objects;

public class OrecocoToken {
    private final String value;

    public OrecocoToken() {
        try {
            ProcessBuilder pb = new ProcessBuilder("/bin/sh", "-c", "cat /dev/urandom | base64 | fold -w 32 | head -n 1");
            pb.redirectErrorStream(true);
            Process process = pb.start();
            try (BufferedReader reader = process.inputReader()) {
                this.value = reader.readLine().trim();
            }
            process.waitFor();
        } catch (IOException ex) {
            throw new UncheckedIOException(ex);
        } catch (InterruptedException ex) {
            throw new IllegalStateException(ex);
        }
    }

    @JsonProperty("token")
    public String asString() {
        return value;
    }

    public boolean matches(String tokenString) {
        return Objects.equals(value, tokenString);
    }

    @Override
    public String toString() {
        return asString();
    }
}
