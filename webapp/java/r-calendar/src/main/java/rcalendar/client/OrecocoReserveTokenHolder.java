package rcalendar.client;

import javax.inject.Singleton;

@Singleton
public class OrecocoReserveTokenHolder {
    private String token;

    public String set(String token) {
        this.token = token;
        return token;
    }

    public String get() {
        return token;
    }
}
