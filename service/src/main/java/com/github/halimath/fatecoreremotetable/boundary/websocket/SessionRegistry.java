package com.github.halimath.fatecoreremotetable.boundary.websocket;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import javax.enterprise.context.ApplicationScoped;
import javax.websocket.Session;

import com.github.halimath.fatecoreremotetable.entity.User;

import lombok.NonNull;

@ApplicationScoped
class SessionRegistry {
    private final Map<String, Session> sessionMap = new ConcurrentHashMap<>();

    void put(@NonNull final User user, @NonNull Session session) {
        sessionMap.put(user.getId(), session);
    }
    
    Optional<Session> get(@NonNull final User user) {
        return Optional.ofNullable(sessionMap.get(user.getId()));
    }

    void remove(@NonNull final User user) {
        sessionMap.remove(user);
    }
}
