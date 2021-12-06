package com.github.halimath.fatecoreremotetable.control;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.smallrye.mutiny.Uni;
import lombok.NonNull;

@ApplicationScoped
class TableRepository {
    private final Map<String, Table> tables = new HashMap<>();
    private final ExecutorService executor = Executors.newSingleThreadExecutor();

    Uni<Table> findByGamemaster(@NonNull final User user) {
        return Uni.createFrom().completionStage(CompletableFuture.supplyAsync(() -> tables.values().stream().filter(t -> t.getGamemaster().equals(user)).findFirst().orElse(null), executor));
    }

    Uni<Table> findByPlayer(@NonNull final User user) {
        return Uni.createFrom().completionStage(CompletableFuture.supplyAsync(() -> tables.values().stream().filter(t -> t.findPlayer(user.getId()).isPresent()).findFirst().orElse(null), executor));
    }
    
    Uni<Table> findById(@NonNull final String id) {
        return Uni.createFrom().completionStage(CompletableFuture.supplyAsync(() -> tables.get(id), executor));
    }

    Uni<Table> save (@NonNull final Table table) {
        return Uni.createFrom().completionStage(CompletableFuture.supplyAsync(() -> {
            tables.put(table.getId(), table);
            return table;
        }, executor));
    }

    Uni<Void> delete(@NonNull final Table table) {
        return Uni.createFrom().completionStage(CompletableFuture.supplyAsync(() -> {
            tables.remove(table.getId());
            return null;
        }, executor));
    }
}
