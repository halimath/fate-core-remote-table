package com.github.halimath.fatecoreremotetable.boundary;

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
import com.github.halimath.fatecoreremotetable.boundary.dto.Request;
import com.github.halimath.fatecoreremotetable.boundary.dto.Response;
import com.github.halimath.fatecoreremotetable.control.TableController;
import com.github.halimath.fatecoreremotetable.control.TableController.OperationForbiddenException;
import com.github.halimath.fatecoreremotetable.control.TableController.PlayerNotFoundException;
import com.github.halimath.fatecoreremotetable.control.TableController.TableControllerException;
import com.github.halimath.fatecoreremotetable.control.TableController.TableNotFoundException;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.netty.handler.codec.http.HttpResponseStatus;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ServerEndpoint("/table")
@ApplicationScoped
@Slf4j
@RequiredArgsConstructor
public class TableWebsocketEndpoint {
    private final TableController tableController;
    private final CommandDispatcher commandDispatcher;
    private final RequestDeserializer requestDeserializer;
    private final ResponseSerializer responseSerializer;
    private final Map<String, Session> sessionMap = new ConcurrentHashMap<>();

    @OnOpen
    public void onOpen(final Session session) {
        log.info("Session started: sessionId={}", session.getId());
        sessionMap.put(session.getId(), session);
    }

    @OnClose
    public void onClose(final Session session) {
        log.info("Session closed: sessionId={}", session.getId());

        sessionMap.remove(session.getId());
        tableController.leave(new User(session.getId())).ifPresent(this::notifyUsers);
    }

    @OnError
    public void onError(final Session session, final Throwable throwable) {
        log.warn("Session error: sessionId={}", session.getId(), throwable);
    }

    @OnMessage
    public void onMessage(final Session session, final String message) {
        log.info("Received message: sessionId={}", session.getId());

        final Request request;

        try {
            request = requestDeserializer.deserialize(message);
        } catch (RequestDeserializationFailedException e) {
            log.warn("Request parsing error: sessionId={} message={}", session.getId(), message, e);
            sendResponse(session,
                    Response.error(session.getId(), HttpResponseStatus.BAD_REQUEST.code(), "invalid message: " + e.getMessage()));
            return;
        }

        try {
            final var user = new User(session.getId());
            sessionMap.put(user.getId(), session);

            notifyUsers(commandDispatcher.dispatchCommand(user, request.command()));

        } catch (TableNotFoundException e) {
            log.warn("Table not found: {}", e.getMessage());
            sendResponse(session, Response.error(session.getId(), request.id(), HttpResponseStatus.NOT_FOUND.code(), e.getMessage()));

        } catch (PlayerNotFoundException e) {
            log.warn("Player not found: {}", e.getMessage());
            sendResponse(session, Response.error(session.getId(), request.id(), HttpResponseStatus.PRECONDITION_FAILED.code(), e.getMessage()));

        } catch (OperationForbiddenException e) {
            log.warn("Got forbidden while processing {}", message, e);
            sendResponse(session, Response.error(session.getId(), request.id(), HttpResponseStatus.FORBIDDEN.code(), e.getMessage()));

        } catch (TableControllerException e) {
            log.warn("Got processing error while processing {}", message, e);
            sendResponse(session,
                    Response.error(session.getId(), request.id(), HttpResponseStatus.INTERNAL_SERVER_ERROR.code(), e.getMessage()));
        }
    }

    private void notifyUsers(final Table table) {
        table.allUsers()
            .map(u -> sessionMap.get(u.getId()))
            .filter(s -> s != null)
            .forEach(s -> {
                final var json = responseSerializer.serialize(Response.table(s.getId(), table));
                sendResponse(s, json);
            });
    }

    private void sendResponse(final Session session, final Response response) {
        sendResponse(session, responseSerializer.serialize(response));
    }

    private void sendResponse(final Session session, final String message) {
        session.getAsyncRemote().sendText(message);
    }
}
