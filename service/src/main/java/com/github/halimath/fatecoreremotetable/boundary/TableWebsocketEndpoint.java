package com.github.halimath.fatecoreremotetable.boundary;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import javax.enterprise.context.ApplicationScoped;
import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.ServerEndpoint;

import com.github.halimath.fatecoreremotetable.boundary.RequestDeserializer.RequestDeserializationFailedException;
import com.github.halimath.fatecoreremotetable.control.TableController;
import com.github.halimath.fatecoreremotetable.control.TableException;
import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.netty.handler.codec.http.HttpResponseStatus;
import io.smallrye.mutiny.Uni;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ServerEndpoint("/table")
@ApplicationScoped
@Slf4j
@RequiredArgsConstructor
class TableWebsocketEndpoint {
    private final TableController tableController;
    private final CommandDispatcher commandDispatcher;
    private final RequestDeserializer requestDeserializer;
    private final ResponseSerializer responseSerializer;
    private final Map<String, Session> sessionMap = new ConcurrentHashMap<>();

    @OnOpen
    void onOpen(final Session session) {
        log.info("Session started: sessionId={}", session.getId());
        sessionMap.put(session.getId(), session);
    }

    @OnClose
    void onClose(final Session session) {
        log.info("Session closed: sessionId={}", session.getId());

        sessionMap.remove(session.getId());

        tableController.applyCommand(new TableController.Command.Leave(new User(session.getId()))) //
                .flatMap(tableOrPlayers -> {
                    if (tableOrPlayers.table() != null) {
                        return notifyUsers(tableOrPlayers.table());
                    }
                    return Uni.combine().all().unis(tableOrPlayers.players().stream().map(Player::getUser).map(this::disconnectUser).toList())
                            .discardItems();

                }).subscribe().with(ignored -> log.debug("Leave complete"));
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

            sendResponse(session,
                    Response.error(session.getId(), HttpResponseStatus.BAD_REQUEST.code(),
                            "invalid message: " + e.getMessage())) //
                                    .subscribe().with(ignored -> log.debug("command dispatch complete"));
            return;
        }

        final var user = new User(session.getId());
        sessionMap.put(user.getId(), session);

        commandDispatcher.dispatchCommand(user, request) //
                .onItem().transform(this::notifyUsers) //
                .onFailure().invoke(t -> handleException(session, request, t)) //
                .subscribe().with(ignored -> log.debug("command dispatch complete"));
    }

    private void handleException(final Session session, final Request request, final Throwable t) {
        final Response res;

        if (t instanceof TableException.TableNotFound e) {
            log.warn("Table not found: {}", e.getMessage());
            res = Response.error(session.getId(), request.id(), HttpResponseStatus.NOT_FOUND.code(), e.getMessage());

        } else if (t instanceof TableException.PlayerNotFound e) {
            log.warn("Player not found: {}", e.getMessage());
            res = Response.error(session.getId(), request.id(), HttpResponseStatus.PRECONDITION_FAILED.code(),
                    e.getMessage());

        } else if (t instanceof TableException.OperationForbidden e) {
            log.warn("Operation forbidden", e);
            res = Response.error(session.getId(), request.id(), HttpResponseStatus.FORBIDDEN.code(), e.getMessage());

        } else {
            log.warn("Got unexpected error while processing command", t);
            res = Response.error(session.getId(), request.id(), HttpResponseStatus.INTERNAL_SERVER_ERROR.code(),
                    t.getMessage());
        }

        sendResponse(session, res) //
                .subscribe().with(ignored -> log.debug("command dispatch complete"));
    }

    private Uni<Void> notifyUsers(final Table table) {
        log.debug("Notifying users of {}", table.getId());

        return Uni.combine().all().unis(table.allUsers().map(u -> sessionMap.get(u.getId())).filter(s -> s != null)
                .map(s -> sendResponse(s, responseSerializer.serialize(Response.table(s.getId(), table)))).toList())
                .discardItems();
    }

    private Uni<Void> sendResponse(final Session session, final Response response) {
        return sendResponse(session, responseSerializer.serialize(response));
    }

    private Uni<Void> sendResponse(final Session session, final String message) {
        return Uni.createFrom().future(session.getAsyncRemote().sendText(message));
    }

    private boolean isHeartbeatMessage(final String message) {
        return "ping".equals(message);
    }

    private Uni<Void> sendHeartbeatResponse(final Session session) {
        return Uni.createFrom().future(session.getAsyncRemote().sendText("pong"));
    }

    private Uni<Void> disconnectUser(final User user) {
        try {
            sessionMap.get(user.getId()).close();
            return Uni.createFrom().item((Void)null);
        } catch (IOException e) {
            return Uni.createFrom().failure(e);
        }
    }
}
