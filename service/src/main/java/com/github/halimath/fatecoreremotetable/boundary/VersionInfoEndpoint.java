package com.github.halimath.fatecoreremotetable.boundary;

import javax.enterprise.context.ApplicationScoped;
import javax.ws.rs.GET;
import javax.ws.rs.Path;

import org.eclipse.microprofile.config.inject.ConfigProperty;

import lombok.NonNull;

@ApplicationScoped
@Path("/version-info")
public class VersionInfoEndpoint {
    private final VersionInfo info;

    VersionInfoEndpoint(
        @ConfigProperty(name = "app.version") @NonNull final String version, 
        @ConfigProperty(name = "app.commit") @NonNull final String commit) {
            this.info = new VersionInfo(version, commit);
        }

    @GET    
    public VersionInfo versionInfo() {
        return info;
    }

    static record VersionInfo (String version, String commit) {}
}
