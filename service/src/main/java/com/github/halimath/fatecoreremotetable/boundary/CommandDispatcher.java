package com.github.halimath.fatecoreremotetable.boundary;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.boundary.Request.Command;
import com.github.halimath.fatecoreremotetable.control.TableController;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.smallrye.mutiny.Uni;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ApplicationScoped
@RequiredArgsConstructor
@Slf4j
class CommandDispatcher {
    private final TableController tableController;

    Uni<? extends Table> dispatchCommand(@NonNull final User user, @NonNull final Request request) {
        try {
            log.info("Dispatching {}", request);
            return tableController.applyCommand(convertCommand(user, request));
        } catch (final IllegalArgumentException e) {
            log.warn("Received unexpected command {}", request.command(), e);
            return Uni.createFrom().failure(e);
        }
    }

    private TableController.Command<Table> convertCommand(final User user, final Request request) {
        if (request.command() instanceof Command.Create c) {
            return new TableController.Command.Create(user, request.tableId(), c.title());
        }

        if (request.command() instanceof Command.Join c) {
            return new TableController.Command.Join(user, request.tableId(), c.name());
        }

        if (request.command() instanceof Command.UpdateFatePoints c) {
            return new TableController.Command.UpdateFatePoints(user, request.tableId(), c.playerId(), c.fatePoints());
        }

        if (request.command() instanceof Command.SpendFatePoint c) {
            return new TableController.Command.SpendFatePoint(user, request.tableId());
        }

        if (request.command() instanceof Command.AddAspect c) {
            return new TableController.Command.AddAspect(user, request.tableId(), c.name(), c.optionalPlayerId().orElse(null));
        }

        if (request.command() instanceof Command.RemoveAspect c) {
            return new TableController.Command.RemoveAspect(user, request.tableId(), c.id());
        }

        throw new IllegalArgumentException("Unknown command: " + request.command());
    }
}
