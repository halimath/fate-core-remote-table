package com.github.halimath.fatecoreremotetable.control;

import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Aspect;
import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;

@ApplicationScoped
public class TableController {
    /** Contains a mapping from gamemaster id == table id -> Table */
    private final Map<String, Table> tables = new ConcurrentHashMap<>();

    public Table create(@NonNull final User user, final String title) throws TableControllerException {
        if (findByGamemaster(user).isPresent()) {
            throw new OperationForbiddenException();
        }

        if (findByPlayer(user).isPresent()) {
            throw new OperationForbiddenException();
        }

        final var table = new Table(user.getId(), title, user);
        tables.put(table.getId(), table);

        return table;
    }

    public Table join(@NonNull User user, @NonNull final String tableId, final String name)
            throws TableControllerException {
        if (!tables.containsKey(tableId)) {
            throw new TableNotFoundException();
        }

        if (findByPlayer(user).isPresent()) {
            throw new OperationForbiddenException();
        }

        final var table = tables.get(tableId);

        if (table.getGamemaster().equals(user)) {
            throw new OperationForbiddenException();
        }

        table.join(new Player(user, name));

        return table;
    }

    public Optional<Table> leave(@NonNull final User user) {
        // TODO: What if the gamemaster leaves?
        
        return findByPlayer(user).map(t -> {
            t.removePlayer(user);
            return t;
        });
    }

    public Table updateFatePoints(@NonNull final User user, @NonNull final String playerId,
            @NonNull final Integer fatePoints) throws TableControllerException {
        final var table = findByGamemaster(user).orElseThrow(TableNotFoundException::new);

        table.findPlayer(playerId).orElseThrow(() -> new PlayerNotFoundException()).setFatePoints(fatePoints);

        return table;
    }

    public Table spendFatePoint(@NonNull final User user) throws TableControllerException {

        final var table = findByPlayer(user).orElseThrow(PlayerNotFoundException::new);

        final var player = table.findPlayer(user.getId()).orElseThrow(() -> new PlayerNotFoundException());
        if (player.getFatePoints() == 0) {
            throw new OperationForbiddenException();
        }

        player.setFatePoints(player.getFatePoints() - 1);

        return table;
    }

    public Table addAspect(@NonNull final User user, @NonNull final String name, @NonNull Optional<User> targetPlayer)
            throws TableControllerException {
        final var table = findByGamemaster(user).orElseThrow(TableNotFoundException::new);

        final var aspect = new Aspect(UUID.randomUUID().toString(), name);

        if (targetPlayer.isEmpty()) {
            table.addAspect(aspect);
        } else {
            final var player = targetPlayer.flatMap(p -> table.findPlayer(p.getId()))
                    .orElseThrow(() -> new PlayerNotFoundException());
            player.addAspect(aspect);
        }

        return table;
    }

    public Table removeAspect(@NonNull final User user, @NonNull final String id) throws TableControllerException {
        final var table = findByGamemaster(user).orElseThrow(TableNotFoundException::new);

        table.removeAspect(id);

        return table;
    }

    private Optional<Table> findByGamemaster(final User user) {
        return Optional.ofNullable(tables.get(user.getId()));
    }

    private Optional<Table> findByPlayer(final User user) {
        return tables.values().stream().filter(t -> t.findPlayer(user.getId()).isPresent()).findFirst();
    }

    public static class TableControllerException extends Exception {
    }

    public static class TableNotFoundException extends TableControllerException {
    }

    public static class PlayerNotFoundException extends TableControllerException {
    }

    public static class OperationForbiddenException extends TableControllerException {
    }
}
