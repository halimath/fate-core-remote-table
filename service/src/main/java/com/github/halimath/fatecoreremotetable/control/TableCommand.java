package com.github.halimath.fatecoreremotetable.control;

import java.util.Optional;

import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;

public interface TableCommand {
    User user();

    record Create(@NonNull User user, @NonNull String tableId, @NonNull String title)
            implements TableCommand {
        public static final String LISTENER_NAME = "TableCommand.Create";
    }

    record Join( //
            @NonNull User user, //
            @NonNull String tableId, //
            @NonNull String name) implements TableCommand {
        public static final String LISTENER_NAME = "TableCommand.Join";
    }

    record AddAspect(//
            @NonNull User user, //
            @NonNull String tableId, //
            @NonNull String name, //
            String playerId) implements TableCommand {

        public static final String LISTENER_NAME = "TableCommand.AddAspect";
        public Optional<String> optionalPlayerId() {
            return Optional.ofNullable(playerId);
        }
    }

    record RemoveAspect(//
            @NonNull User user, //
            @NonNull String tableId, //
            @NonNull String id) implements TableCommand {
        public static final String LISTENER_NAME = "TableCommand.RemoveAspect";
    }

    record SpendFatePoint(@NonNull User user, @NonNull String tableId) implements TableCommand {
        public static final String LISTENER_NAME = "TableCommand.SpendFatePoint";
    }

    record UpdateFatePoints(//
            @NonNull User user, //
            @NonNull String tableId, //
            @NonNull String playerId, //
            @NonNull Integer fatePoints) implements TableCommand {
        public static final String LISTENER_NAME = "TableCommand.UpdateFatePoints";

    }

    record Leave(@NonNull User user) implements TableCommand {
        public static final String LISTENER_NAME = "TableCommand.Leave";
    }
}