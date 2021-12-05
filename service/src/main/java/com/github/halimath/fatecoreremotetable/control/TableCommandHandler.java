package com.github.halimath.fatecoreremotetable.control;

import java.util.stream.Collectors;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.quarkus.vertx.ConsumeEvent;
import io.vertx.mutiny.core.eventbus.EventBus;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@ApplicationScoped
@RequiredArgsConstructor
class TableCommandHandler {
        private final TableService service;
        private final EventBus eventBus;

        @ConsumeEvent(TableCommand.Create.LISTENER_NAME)
        void handle(@NonNull final TableCommand.Create command) {
                service.create(command.user(), command.tableId(), command.title())
                                .subscribe().with(
                                                table -> eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                new TableEvent.Updated(table)),
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));
        }

        @ConsumeEvent(TableCommand.Join.LISTENER_NAME)
        void handle(@NonNull final TableCommand.Join command) {
                service.join(command.user(), command.tableId(), command.name())
                                .subscribe().with(
                                                table -> eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                new TableEvent.Updated(table)),
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));
        }

        @ConsumeEvent(TableCommand.AddAspect.LISTENER_NAME)
        void handle(@NonNull final TableCommand.AddAspect command) {
                service.addAspect(command.user(), command.tableId(), command.name(),
                                command.optionalPlayerId().map(User::new).orElse(null))
                                .subscribe().with(
                                                table -> eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                new TableEvent.Updated(table)),
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));

        }

        @ConsumeEvent(TableCommand.RemoveAspect.LISTENER_NAME)
        void handle(@NonNull final TableCommand.RemoveAspect command) {
                service.removeAspect(command.user(), command.tableId(), command.id())
                                .subscribe().with(
                                                table -> eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                new TableEvent.Updated(table)),
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));
        }

        @ConsumeEvent(TableCommand.SpendFatePoint.LISTENER_NAME)
        void handle(@NonNull final TableCommand.SpendFatePoint command) {
                service.spendFatePoint(command.user(), command.tableId())
                                .subscribe().with(
                                                table -> eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                new TableEvent.Updated(table)),
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));

        }

        @ConsumeEvent(TableCommand.UpdateFatePoints.LISTENER_NAME)
        void handle(@NonNull final TableCommand.UpdateFatePoints command) {
                service.updateFatePoints(command.user(), command.tableId(), command.playerId(), command.fatePoints())
                                .subscribe().with(
                                                table -> eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                new TableEvent.Updated(table)),
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));
        }

        @ConsumeEvent(TableCommand.Leave.LISTENER_NAME)
        void handle(@NonNull final TableCommand.Leave command) {
                service.leave(command.user())
                                .subscribe().with(
                                                tableOrPlayers -> {
                                                        if (tableOrPlayers.table() != null) {
                                                                eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                                new TableEvent.Updated(tableOrPlayers
                                                                                                .table()));
                                                        } else {
                                                                eventBus.publish(TableEvent.Updated.LISTENER_NAME,
                                                                                new TableEvent.Closed(tableOrPlayers
                                                                                                .players().stream()
                                                                                                .map(Player::getUser)
                                                                                                .collect(Collectors
                                                                                                                .toSet())));
                                                        }
                                                },
                                                t -> eventBus.publish(TableEvent.Error.LISTENER_NAME,
                                                                new TableEvent.Error(command.user(), t)));

        }
}
