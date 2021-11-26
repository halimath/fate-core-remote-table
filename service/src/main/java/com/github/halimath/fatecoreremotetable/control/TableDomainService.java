package com.github.halimath.fatecoreremotetable.control;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Aspect;
import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;

/**
 * {@link TableDomainService} contains the business logic for tables. This class
 * acts as a facade for all business operations.
 * <p>
 * All tables are stored in an internal Map and not concurrency synchronisation
 * is done. It is up to the caller to ensure thread safety. For a thread safe
 * async variant of this class, see {@link AsyncTableController}.
 */
@ApplicationScoped
class TableDomainService {
    /** Contains a mapping from gamemaster (id == table id) -> Table */
    private final Map<String, Table> tables = new HashMap<>();

    Table create(@NonNull final User user, final String title) throws TableException {
        if (findByGamemaster(user).isPresent()) {
            throw new TableException.OperationForbidden();
        }

        if (findByPlayer(user).isPresent()) {
            throw new TableException.OperationForbidden();
        }

        final var table = new Table(user.getId(), title, user);
        tables.put(table.getId(), table);

        return table;
    }

    Table join(@NonNull final User user, @NonNull final String tableId, final String name) throws TableException {
        if (!tables.containsKey(tableId)) {
            throw new TableException.TableNotFound();
        }

        if (findByPlayer(user).isPresent()) {
            throw new TableException.OperationForbidden();
        }

        final var table = tables.get(tableId);

        if (table.getGamemaster().equals(user)) {
            throw new TableException.OperationForbidden();
        }

        table.join(new Player(user, name));

        return table;
    }

    Table leave(@NonNull final User user) {
        if (tables.containsKey(user.getId())) {
            // The game master leaves the table.
            // TODO: What to do here?
            return tables.get(user.getId());
        }

        return findByPlayer(user).map(t -> {
            t.removePlayer(user);
            return t;
        }).orElseThrow(() -> new TableException.PlayerNotFound());
    }

    Table updateFatePoints(@NonNull final User user, @NonNull final String playerId, @NonNull final Integer fatePoints)
            throws TableException {
        final var table = findByGamemaster(user).orElseThrow(TableException.TableNotFound::new);

        table.findPlayer(playerId).orElseThrow(() -> new TableException.PlayerNotFound()).setFatePoints(fatePoints);

        return table;
    }

    Table spendFatePoint(@NonNull final User user) throws TableException {

        final var table = findByPlayer(user).orElseThrow(TableException.PlayerNotFound::new);

        final var player = table.findPlayer(user.getId()).orElseThrow(() -> new TableException.PlayerNotFound());
        if (player.getFatePoints() == 0) {
            throw new TableException.OperationForbidden();
        }

        player.setFatePoints(player.getFatePoints() - 1);

        return table;
    }

    Table addAspect(@NonNull final User user, @NonNull final String name, @NonNull Optional<User> targetPlayer)
            throws TableException {
        final var table = findByGamemaster(user).orElseThrow(TableException.TableNotFound::new);

        final var aspect = new Aspect(UUID.randomUUID().toString(), name);

        if (targetPlayer.isEmpty()) {
            table.addAspect(aspect);
        } else {
            final var player = targetPlayer.flatMap(p -> table.findPlayer(p.getId()))
                    .orElseThrow(() -> new TableException.PlayerNotFound());
            player.addAspect(aspect);
        }

        return table;
    }

    Table removeAspect(@NonNull final User user, @NonNull final String id) throws TableException {
        final var table = findByGamemaster(user).orElseThrow(TableException.TableNotFound::new);

        table.removeAspect(id);

        return table;
    }

    private Optional<Table> findByGamemaster(final User user) {
        return Optional.ofNullable(tables.get(user.getId()));
    }

    private Optional<Table> findByPlayer(final User user) {
        return tables.values().stream().filter(t -> t.findPlayer(user.getId()).isPresent()).findFirst();
    }
}
