package com.github.halimath.fatecoreremotetable.control;

import java.util.Optional;

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
    public interface Command<R> {
        User user();

        public record Create(@NonNull User user, @NonNull String tableId, @NonNull String title)
                implements Command<Table> {
        }

        public record Join( //
                @NonNull User user, //
                @NonNull String tableId, //
                @NonNull String name) implements Command<Table> {
        }

        public record AddAspect(//
                @NonNull User user, //
                @NonNull String tableId, //
                @NonNull String name, //
                String playerId) implements Command<Table> {
            public Optional<String> optionalPlayerId() {
                return Optional.ofNullable(playerId);
            }
        }

        public record RemoveAspect(//
                @NonNull User user, //
                @NonNull String tableId, //
                @NonNull String id) implements Command<Table> {
        }

        public record SpendFatePoint(@NonNull User user, @NonNull String tableId) implements Command<Table> {
        }

        public record UpdateFatePoints(//
                @NonNull User user, //
                @NonNull String tableId, //
                @NonNull String playerId, //
                @NonNull Integer fatePoints) implements Command<Table> {
        }

        public record Leave(@NonNull User user) implements Command<TableOrPlayers> {
        }
    }

    private final TableDomainService domainService;

    @SuppressWarnings("unchecked")
    public <T> Uni<? extends T> applyCommand(@NonNull final Command<T> command) {
        log.info("Processing command {}", command);

        if (command instanceof Command.Create c) {
            return (Uni<? extends T>) domainService.create(c.user(), c.tableId(), c.title());
        }

        if (command instanceof Command.Join c) {
            return (Uni<? extends T>) domainService.join(c.user(), c.tableId(), c.name());
        }

        if (command instanceof Command.UpdateFatePoints c) {
            return (Uni<? extends T>) domainService.updateFatePoints(c.user(), c.user().getId(), c.playerId(),
                    c.fatePoints());
        }

        if (command instanceof Command.SpendFatePoint c) {
            return (Uni<? extends T>) domainService.spendFatePoint(c.user(), c.tableId());
        }

        if (command instanceof Command.AddAspect c) {
            return (Uni<? extends T>) c.optionalPlayerId()
                    .map(playerId -> domainService.addAspect(c.user(), c.tableId(), c.name(), new User(playerId)))
                    .orElseGet(() -> domainService.addAspect(c.user(), c.tableId(), c.name()));
        }

        if (command instanceof Command.RemoveAspect c) {
            return (Uni<? extends T>) domainService.removeAspect(c.user(), c.tableId(), c.id());
        }

        if (command instanceof Command.Leave t) {
            return (Uni<? extends T>) domainService.leave(t.user());
        }

        log.error("Received unhandled command: {}", command);
        return Uni.createFrom().failure(new IllegalArgumentException("Unknown command " + command));
    }
}
