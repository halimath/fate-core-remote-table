package com.github.halimath.fatecoreremotetable.boundary;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.boundary.Request.Command;
import com.github.halimath.fatecoreremotetable.control.TableCommand;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.vertx.mutiny.core.eventbus.EventBus;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ApplicationScoped
@RequiredArgsConstructor
@Slf4j
public class CommandDispatcher {
    private final EventBus eventBus;

    public void dispatchCommand(@NonNull final User user, @NonNull final Request request) {
        log.info("Dispatching {}", request);

        if (request.command() instanceof Command.Create c) {
            eventBus.publish(TableCommand.Create.LISTENER_NAME,
                    new TableCommand.Create(user, request.tableId(), c.title()));
            return;
        }

        if (request.command() instanceof Command.Join c) {
            eventBus.publish(TableCommand.Join.LISTENER_NAME, new TableCommand.Join(user, request.tableId(), c.name()));
            return;
        }

        if (request.command() instanceof Command.UpdateFatePoints c) {
            eventBus.publish(TableCommand.UpdateFatePoints.LISTENER_NAME,
                    new TableCommand.UpdateFatePoints(user, request.tableId(), c.playerId(), c.fatePoints()));
            return;
        }

        if (request.command() instanceof Command.SpendFatePoint c) {
            eventBus.publish(TableCommand.SpendFatePoint.LISTENER_NAME,
                    new TableCommand.SpendFatePoint(user, request.tableId()));
            return;
        }

        if (request.command() instanceof Command.AddAspect c) {
            eventBus.publish(TableCommand.AddAspect.LISTENER_NAME,
                    new TableCommand.AddAspect(user, request.tableId(), c.name(), c.optionalPlayerId().orElse(null)));
            return;
        }

        if (request.command() instanceof Command.RemoveAspect c) {
            eventBus.publish(TableCommand.RemoveAspect.LISTENER_NAME,
                    new TableCommand.RemoveAspect(user, request.tableId(), c.id()));
            return;
        }

        throw new IllegalArgumentException("Unknown command: " + request.command());
    }
}
