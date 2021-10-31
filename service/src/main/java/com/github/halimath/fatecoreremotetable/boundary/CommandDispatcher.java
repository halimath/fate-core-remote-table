package com.github.halimath.fatecoreremotetable.boundary;

import java.util.Optional;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

import com.github.halimath.fatecoreremotetable.boundary.commands.AddAspect;
import com.github.halimath.fatecoreremotetable.boundary.commands.JoinTable;
import com.github.halimath.fatecoreremotetable.boundary.commands.NewTable;
import com.github.halimath.fatecoreremotetable.boundary.commands.RemoveAspect;
import com.github.halimath.fatecoreremotetable.boundary.commands.SpendFatePoint;
import com.github.halimath.fatecoreremotetable.boundary.commands.UpdateFatePoints;
import com.github.halimath.fatecoreremotetable.control.TableController;
import com.github.halimath.fatecoreremotetable.control.TableController.TableControllerException;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;

@ApplicationScoped
@AllArgsConstructor
@Slf4j
public class CommandDispatcher {
    @Inject
    private final TableController tableController;

    public Table dispatchCommand(@NonNull final User user, @NonNull final Command command)
            throws TableControllerException {

        log.info("Dispatching {}", command);

        if (command instanceof NewTable c) {
            return tableController.startNew(user, c.getTitle());
        }

        if (command instanceof JoinTable c) {
            return tableController.join(user, c.getTableId(), c.getName());
        }

        if (command instanceof UpdateFatePoints c) {
            return tableController.updateFatePoints(user, c.getTableId(), c.getPlayerId(), c.getFatePoints());
        }

        if (command instanceof SpendFatePoint c) {
            return tableController.spendFatePoint(user, c.getTableId());
        }

        if (command instanceof AddAspect c) {
            return tableController.addAspect(user, c.getTableId(), c.getName(), 
            c.getOptionalPlayerId().map(id -> new User(id)));
        }

        if (command instanceof RemoveAspect c) {
            return tableController.removeAspect(user, c.getTableId(), c.getId());
        }

        throw new IllegalArgumentException("Unknown command: " + command);
    }
}
