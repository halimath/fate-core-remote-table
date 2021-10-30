package com.github.halimath.fatecoreremotetable.control;

import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

import javax.enterprise.context.ApplicationScoped;

import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;

@ApplicationScoped
public class TableController {
    private final Map<String, Table> tables = new ConcurrentHashMap<>();
    
    public Table startNew (@NonNull User user, @NonNull final String title) {
        final var table = new Table(UUID.randomUUID().toString(), title, user);
        tables.put(table.getId(), table);
        return table;
    }

    public Table join(@NonNull final User user, @NonNull final String tableId, @NonNull final String name) throws TableControllerException {
        if (!tables.containsKey(tableId)) {
            throw new TableNotFoundException();
        }
        final var table = tables.get(tableId);
        table.join(new Player(user, name));
        return table;
    }

    public Optional<Table> disconnect(@NonNull final User user) {
        return tables.values().stream()
            .filter(t -> t.findPlayer(user.getId()).isPresent())
            .findFirst()
            .map(t -> {
                t.removePlayer(user);
                return t;
            });
    }

    public Table updateFatePoints(@NonNull final User user, @NonNull final String tableId, @NonNull final String playerId, @NonNull final Integer fatePoints) throws TableControllerException {
        if (!tables.containsKey(tableId)) {
            throw new TableNotFoundException();
        }

        final var table = tables.get(tableId);

        if (!table.getGameMaster().equals(user)) {
            throw new OperationForbiddenException();
        }


        table.findPlayer(playerId)
            .orElseThrow(() -> new PlayerNotFoundException())
            .setFatePoints(fatePoints);

        return table;
    }

    public Table spendFatePoint (@NonNull final User user, @NonNull final String tableId) throws TableControllerException {
        if (!tables.containsKey(tableId)) {
            throw new TableNotFoundException();
        }

        final var table = tables.get(tableId);

        final var player = table.findPlayer(user.getId())
            .orElseThrow(() -> new PlayerNotFoundException());
        if (player.getFatePoints() == 0) {
            throw new OperationForbiddenException();       
        }
        
        player.setFatePoints(player.getFatePoints() - 1);

        return table;
    }

    public static class TableControllerException extends Exception {}
    public static class TableNotFoundException extends TableControllerException {}
    public static class PlayerNotFoundException extends TableControllerException {}
    public static class OperationForbiddenException extends TableControllerException {}
}
