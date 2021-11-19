package com.github.halimath.fatecoreremotetable.control;

import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@ApplicationScoped
@RequiredArgsConstructor
public class AsyncTableController {
    private final TableController tableController;
    private final ExecutorService executor = Executors.newSingleThreadExecutor();

    public CompletableFuture<Table> create(@NonNull final User user, final String title) {
        return CompletableFuture.supplyAsync(() -> tableController.create(user, title), executor);
    }

    public CompletableFuture<Table> join(@NonNull final User user, @NonNull final String tableId, final String name) {
        return CompletableFuture.supplyAsync(() -> tableController.join(user, tableId, name), executor);
    }

    public CompletableFuture<TableController.LeaveResult> leave(@NonNull final User user) {
        return CompletableFuture.supplyAsync(() -> tableController.leave(user), executor);
    }

    public CompletableFuture<Table> updateFatePoints(@NonNull final User user, @NonNull final String playerId,
            @NonNull final Integer fatePoints) {
                return CompletableFuture.supplyAsync(() -> tableController.updateFatePoints(user, playerId, fatePoints), executor);
    }

    public CompletableFuture<Table> spendFatePoint(@NonNull final User user) {
        return CompletableFuture.supplyAsync(() -> tableController.spendFatePoint(user), executor);
    }

    public CompletableFuture<Table> addAspect(@NonNull final User user, @NonNull final String name, @NonNull Optional<User> targetPlayer) {
        return CompletableFuture.supplyAsync(() -> tableController.addAspect(user, name, targetPlayer), executor);
    }

    public CompletableFuture<Table> removeAspect(@NonNull final User user, @NonNull final String id) {
        return CompletableFuture.supplyAsync(() -> tableController.removeAspect(user, id), executor);
    }    
}
