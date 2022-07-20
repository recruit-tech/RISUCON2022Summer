package orecoco.reserve.resource;

import orecoco.reserve.model.OrecocoToken;

import javax.inject.Singleton;

@Singleton
public class TokenHolder {
    private OrecocoToken orecocoToken;

    OrecocoToken set(OrecocoToken orecocoToken) {
        this.orecocoToken = orecocoToken;
        return orecocoToken;
    }

    public OrecocoToken get() {
        return orecocoToken;
    }
}
