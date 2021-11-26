package com.github.halimath.fatecoreremotetable.control;

import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.smallrye.mutiny.Uni;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ApplicationScoped
@RequiredArgsConstructor
@Slf4j
public class TableController {
    public interface Command {
        User user();

        public record Create(@NonNull User user, @NonNull String title) implements Command {
        }

        public record Join( //
                        @NonNull User user, //
                        @NonNull String tableId, //
                        @NonNull String name) implements Command {
        }

        public record AddAspect(//
                        @NonNull User user, //
                        @NonNull String name, //
                        String playerId) implements Command {
                public Optional<String> optionalPlayerId() {
                        return Optional.ofNullable(playerId);
                }
        }

        public record RemoveAspect(//
                        @NonNull User user, //
                        @NonNull String id) implements Command {
        }

        public record SpendFatePoint(@NonNull User user) implements Command {
        }

        public record UpdateFatePoints(//
                        @NonNull User user, //
                        @NonNull String playerId, //
                        @NonNull Integer fatePoints) implements Command {
        }

        public record Leave(@NonNull User user) implements Command {}
}


    private final TableDomainService domainService;
    // TODO: This currently means that all mutations for all tables are performed on a single executor.
    private final ExecutorService executorService = Executors.newSingleThreadExecutor();

    public Uni<Table> applyCommand(@NonNull final Command command) {
        log.info("Processing command {}", command);

        if (command instanceof Command.Create c) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.create(c.user(), c.title()), executorService));
        }

        if (command instanceof Command.Join c) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.join(c.user(), c.tableId(), c.name())));
        }

        if (command instanceof Command.UpdateFatePoints c) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.updateFatePoints(c.user(), c.playerId(), c.fatePoints())));
        }

        if (command instanceof Command.SpendFatePoint c) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.spendFatePoint(c.user())));
        }

        if (command instanceof Command.AddAspect c) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.addAspect(c.user(), c.name(), c.optionalPlayerId().map(id -> new User(id)))));
        }

        if (command instanceof Command.RemoveAspect c) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.removeAspect(c.user(), c.id())));
        }

        if (command instanceof Command.Leave t) {
            return Uni.createFrom().completionStage(
                    CompletableFuture.supplyAsync(() -> domainService.leave(t.user())));
        }

        log.error("Received unhandled command: {}", command);
        return Uni.createFrom().failure(new IllegalArgumentException("Unknown command " + command));
    }
}
