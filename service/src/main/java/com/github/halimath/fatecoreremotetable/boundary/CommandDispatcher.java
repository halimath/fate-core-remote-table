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

    Uni<Table> dispatchCommand(@NonNull final User user, @NonNull final Command command) {
        try {
            log.info("Dispatching {}", command);
            return tableController.applyCommand(convertCommand(user, command));
        } catch (final IllegalArgumentException e) {
            log.warn("Received unexpected command {}", command, e);
            return Uni.createFrom().failure(e);
        }
    }

    private TableController.Command convertCommand(final User user, final Request.Command command) {
        if (command instanceof Command.Create c) {
            return new TableController.Command.Create(user, c.title());
        }

        if (command instanceof Command.Join c) {
            return new TableController.Command.Join(user, c.tableId(), c.name());
        }

        if (command instanceof Command.UpdateFatePoints c) {
            return new TableController.Command.UpdateFatePoints(user, c.playerId(), c.fatePoints());
        }

        if (command instanceof Command.SpendFatePoint c) {
            return new TableController.Command.SpendFatePoint(user);
        }

        if (command instanceof Command.AddAspect c) {
            return new TableController.Command.AddAspect(user, c.name(), c.optionalPlayerId().orElse(null));
        }

        if (command instanceof Command.RemoveAspect c) {
            return new TableController.Command.RemoveAspect(user, c.id());
        }

        throw new IllegalArgumentException("Unknown command: " + command);
    }
}
