package com.github.halimath.fatecoreremotetable.boundary.websocket;

import javax.enterprise.context.ApplicationScoped;
import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.ServerEndpoint;

import com.github.halimath.fatecoreremotetable.boundary.CommandDispatcher;
import com.github.halimath.fatecoreremotetable.boundary.Request;
import com.github.halimath.fatecoreremotetable.boundary.RequestDeserializer;
import com.github.halimath.fatecoreremotetable.boundary.RequestDeserializer.RequestDeserializationFailedException;
import com.github.halimath.fatecoreremotetable.control.TableCommand;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.eventbus.EventBus;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ServerEndpoint("/table")
@ApplicationScoped
@Slf4j
@RequiredArgsConstructor
class TableWebsocketEndpoint {
    private final EventBus eventBus;
    private final CommandDispatcher commandDispatcher;
    private final RequestDeserializer requestDeserializer;
    private final SessionRegistry sessionRegistry;

    @OnOpen
    void onOpen(final Session session) {
        log.info("Session started: sessionId={}", session.getId());
        sessionRegistry.put(new User(session.getId()), session);
    }

    @OnClose
    void onClose(final Session session) {
        log.info("Session closed: sessionId={}", session.getId());

        final var user = new User(session.getId());

        sessionRegistry.remove(user);

        eventBus.publish(TableCommand.Leave.LISTENER_NAME, new TableCommand.Leave(user));
    }

    @OnError
    void onError(final Session session, final Throwable throwable) {
        log.warn("Session error: sessionId={}", session.getId(), throwable);
    }

    @OnMessage
    public void onMessage(final Session session, final String message) {
        if (isHeartbeatMessage(message)) {
            log.debug("Got heartbeat message from {}", session.getId());
            sendHeartbeatResponse(session);
            return;
        }

        log.info("Received message: sessionId={}", session.getId());
        final Request request;

        try {
            request = requestDeserializer.deserialize(message);
        } catch (RequestDeserializationFailedException e) {
            log.warn("Request parsing error: sessionId={} message={}", session.getId(), message, e);

            return;
        }

        final var user = new User(session.getId());

        commandDispatcher.dispatchCommand(user, request);
    }

    private boolean isHeartbeatMessage(final String message) {
        return "ping".equals(message);
    }

    private Uni<Void> sendHeartbeatResponse(final Session session) {
        return Uni.createFrom().future(session.getAsyncRemote().sendText("pong"));
    }    

}
