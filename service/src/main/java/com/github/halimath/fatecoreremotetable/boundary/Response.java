package com.github.halimath.fatecoreremotetable.boundary;

import java.util.List;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.NonNull;

@JsonInclude(JsonInclude.Include.NON_NULL)
record Response(String id, String self, Type type, Table table, Error error) {
    enum Type {
        @JsonProperty("table")
        TABLE,

        @JsonProperty("error")
        ERROR

        // TODO: NOTIFICATION
    }

    static Response error(@NonNull final String self, final String requestId, final int code,
            @NonNull final String reason) {
        return new Response(UUID.randomUUID().toString(), self, Type.ERROR, null, new Error(requestId, code, reason));
    }

    static Response error(@NonNull final String self, final int code, final String reason) {
        return new Response(UUID.randomUUID().toString(), self, Type.ERROR, null, new Error(null, code, reason));
    }

    static Response table(@NonNull final String self,
            @NonNull final com.github.halimath.fatecoreremotetable.entity.Table table) {
        return new Response(UUID.randomUUID().toString(), self, Type.TABLE, Table.fromEntity(table), null);
    }

    record Error(String requestId, int code, String reason) {
    }

    record Table(String id, String title, String gamemaster, List<Player> players, List<Aspect> aspects) {

        static Table fromEntity(@NonNull final com.github.halimath.fatecoreremotetable.entity.Table table) {
            return new Table(table.getId(), table.getTitle(), table.getGamemaster().getId(), table.getPlayers().stream()
                    .map(p -> new Table.Player(p.getUser().getId(), p.getName(), p.getFatePoints(),
                            p.getAspects().stream().map(a -> new Table.Aspect(a.getId(), a.getName())).toList()))
                    .toList(), table.getAspects().stream().map(a -> new Table.Aspect(a.getId(), a.getName())).toList());
        }
        record Player(String id, String name, Integer fatePoints, List<Aspect> aspects) {
        }

        record Aspect(String id, String name) {
        }
    }
}
