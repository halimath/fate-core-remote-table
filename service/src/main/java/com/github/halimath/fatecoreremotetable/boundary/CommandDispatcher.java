package com.github.halimath.fatecoreremotetable.boundary;

import java.util.concurrent.CompletableFuture;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.boundary.dto.Command;
import com.github.halimath.fatecoreremotetable.control.AsyncTableController;
import com.github.halimath.fatecoreremotetable.control.TableController.TableControllerException;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ApplicationScoped
@RequiredArgsConstructor
@Slf4j
public class CommandDispatcher {
    private final AsyncTableController
     tableController;

    public CompletableFuture<Table> dispatchCommand(@NonNull final User user, @NonNull final Command command) throws TableControllerException {

        log.info("Dispatching {}", command);

        if (command instanceof Command.Create c) {
            return tableController.create(user, c.title());
        }

        if (command instanceof Command.Join c) {
            return tableController.join(user, c.tableId(), c.name());
        }

        if (command instanceof Command.UpdateFatePoints c) {
            return tableController.updateFatePoints(user, c.playerId(), c.fatePoints());
        }

        if (command instanceof Command.SpendFatePoint c) {
            return tableController.spendFatePoint(user);
        }

        if (command instanceof Command.AddAspect c) {
            return tableController.addAspect(user, c.name(), c.optionalPlayerId().map(id -> new User(id)));
        }

        if (command instanceof Command.RemoveAspect c) {
            return tableController.removeAspect(user, c.id());
        }

        return CompletableFuture.failedFuture(new IllegalArgumentException("Unknown command: " + command));
    }
}
