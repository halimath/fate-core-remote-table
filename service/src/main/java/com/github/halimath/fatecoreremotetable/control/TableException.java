package com.github.halimath.fatecoreremotetable.control;

import lombok.NonNull;

public abstract class TableException extends RuntimeException {
    protected TableException(@NonNull final String message) {
        super(message);
    }

    public static class TableNotFound extends TableException {
        TableNotFound() {
            super("Table not found");
        }
    }

    public static class Conflict extends TableException {
        Conflict() {
            super("Conflict");
        }
    }

    public static class PlayerNotFound extends TableException {
        PlayerNotFound() {
            super("Player not found");
        }
    }

    public static class OperationForbidden extends TableException {
        OperationForbidden() {
            super("Operation forbidden");
        }
    }
    
}
