package com.github.halimath.fatecoreremotetable.control;

import java.util.UUID;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Aspect;
import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.smallrye.mutiny.Uni;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

/**
 * {@link TableService} implements domain services {@link Table}s. This
 * class provides method to carry
 * out the operations supported on {@link Table}s.
 */
@ApplicationScoped
@RequiredArgsConstructor
class TableService {
    private final TableRepository repository;

    Uni<Table> create(@NonNull final User user, @NonNull String tableId, @NonNull final String title) {
        return repository.findById(tableId)
                .onItem().ifNotNull().failWith(TableException.Conflict::new) //

                // If a table with the given user as a gamemaster already exists, refuse
                // to create a new one
                .chain(() -> repository.findByGamemaster(user)) //
                .onItem().ifNotNull().failWith(TableException.OperationForbidden::new) //

                // If a table with the given user as a player exists, refuse to create one
                .chain(() -> repository.findByPlayer(user)) //
                .onItem().ifNotNull().failWith(TableException.OperationForbidden::new) //

                // Create a new table and save it
                .chain(() -> {
                    return repository.save(new Table(tableId, title, user));
                });
    }

    Uni<Table> join(@NonNull final User user, @NonNull final String tableId, final String name) throws TableException {
        return repository.findByPlayer(user)
                .onItem().ifNotNull().failWith(TableException.OperationForbidden::new)

                .chain(() -> repository.findByGamemaster(user))
                .onItem().ifNotNull().failWith(TableException.OperationForbidden::new)

                .chain(() -> repository.findById(tableId))
                .onItem().ifNull().failWith(TableException.TableNotFound::new)

                .flatMap(table -> {
                    if (table.getGamemaster().equals(user)) {
                        throw new TableException.OperationForbidden();
                    }
                    table.join(new Player(user, name));
                    return repository.save(table);
                });
    }

    Uni<TableOrPlayers> leave(@NonNull final User user) {
        return repository.findByGamemaster(user)
                .onItem().ifNotNull().transformToUni(table -> repository.delete(table)
                        .map(ignored -> new TableOrPlayers(null, table.getPlayers())))
                .onItem().ifNull().switchTo(() -> repository.findByPlayer(user)
                        .onItem().ifNull().failWith(TableException.PlayerNotFound::new)
                        .flatMap(table -> {
                            table.removePlayer(user);
                            return repository.save(table).map(savedTable -> new TableOrPlayers(savedTable, null));
                        }));
    }

    Uni<Table> updateFatePoints(@NonNull final User user, @NonNull final String tableId, @NonNull final String playerId,
            @NonNull final Integer fatePoints) {
        return repository.findById(tableId)
                .onItem().ifNull().failWith(TableException.TableNotFound::new)

                .flatMap(table -> {
                    if (!table.getGamemaster().equals(user)) {
                        throw new TableException.OperationForbidden();
                    }

                    table.findPlayer(playerId).orElseThrow(TableException.PlayerNotFound::new)
                            .setFatePoints(fatePoints);

                    return repository.save(table);
                });
    }

    Uni<Table> spendFatePoint(@NonNull final User user, @NonNull final String tableId) {
        return repository.findById(tableId)
                .onItem().ifNull().failWith(TableException.TableNotFound::new)

                .flatMap(table -> {
                    final var player = table.findPlayer(user.getId())
                            .orElseThrow(TableException.PlayerNotFound::new);
                    if (player.getFatePoints() == 0) {
                        throw new TableException.OperationForbidden();
                    }

                    player.setFatePoints(player.getFatePoints() - 1);

                    return repository.save(table);
                });
    }

    Uni<Table> addAspect(@NonNull final User user, @NonNull final String tableId, @NonNull final String name) {
        return addAspect(user, tableId, name, null);
    }

    Uni<Table> addAspect(@NonNull final User user, @NonNull final String tableId, @NonNull final String name,
            final User targetPlayer) {
        return repository.findById(tableId)
                .onItem().ifNull().failWith(TableException.TableNotFound::new)
                .flatMap(table -> {
                    if (!table.getGamemaster().equals(user)) {
                        throw new TableException.OperationForbidden();
                    }

                    final var aspect = new Aspect(UUID.randomUUID().toString(), name);

                    if (targetPlayer == null) {
                        table.addAspect(aspect);
                    } else {
                        final var player = table.findPlayer(targetPlayer.getId())
                                .orElseThrow(TableException.PlayerNotFound::new);
                        player.addAspect(aspect);
                    }

                    return repository.save(table);
                });
    }

    Uni<Table> removeAspect(@NonNull final User user, @NonNull final String tableId, @NonNull final String aspectId) {
        return repository.findById(tableId)
                .onItem().ifNull().failWith(TableException.TableNotFound::new)
                .flatMap(table -> {
                    if (!table.getGamemaster().equals(user)) {
                        throw new TableException.OperationForbidden();
                    }
                    table.removeAspect(aspectId);
                    return repository.save(table);
                });
    }
}
