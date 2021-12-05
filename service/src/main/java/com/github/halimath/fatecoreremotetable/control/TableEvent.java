package com.github.halimath.fatecoreremotetable.control;

import java.util.Set;

import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;

public interface TableEvent {

    record Updated(@NonNull Table table) {
        public static final String LISTENER_NAME = "TableEvent.Updated";
    }
    
    record Closed(@NonNull Set<User> users) {
        public static final String LISTENER_NAME = "TableEvent.Closed";
    }

    record Error(@NonNull User user, @NonNull Throwable t) {
        public static final String LISTENER_NAME = "TableEvent.Error";
    }    
}
