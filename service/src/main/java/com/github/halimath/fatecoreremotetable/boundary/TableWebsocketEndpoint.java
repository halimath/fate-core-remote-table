package com.github.halimath.fatecoreremotetable.boundary;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

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
import com.github.halimath.fatecoreremotetable.control.AsyncTableController;
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
    private final AsyncTableController tableController;
    private final CommandDispatcher commandDispatcher;
    private final RequestDeserializer requestDeserializer;
    private final ResponseSerializer responseSerializer;
    private final Map<String, Session> sessionMap = new ConcurrentHashMap<>();
    private final ExecutorService executor = Executors.newFixedThreadPool(10); // TODO: Make this configurable

    @OnOpen
    public void onOpen(final Session session) {
        log.info("Session started: sessionId={}", session.getId());
        sessionMap.put(session.getId(), session);
    }

    @OnClose
    public void onClose(final Session session) {
        log.info("Session closed: sessionId={}", session.getId());

        sessionMap.remove(session.getId());
        tableController.leave(new User(session.getId())).thenAccept(table -> table.ifPresent(this::notifyUsers)).join();
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
            sendResponse(session, Response.error(session.getId(), HttpResponseStatus.BAD_REQUEST.code(),
                    "invalid message: " + e.getMessage())).join();
            return;
        }

        try {
            final var user = new User(session.getId());
            sessionMap.put(user.getId(), session);

            commandDispatcher.dispatchCommand(user, request.command()).thenAccept(this::notifyUsers).join();

        } catch (TableNotFoundException e) {
            log.warn("Table not found: {}", e.getMessage());
            sendResponse(session,
                    Response.error(session.getId(), request.id(), HttpResponseStatus.NOT_FOUND.code(), e.getMessage()))
                            .join();

        } catch (PlayerNotFoundException e) {
            log.warn("Player not found: {}", e.getMessage());
            sendResponse(session, Response.error(session.getId(), request.id(),
                    HttpResponseStatus.PRECONDITION_FAILED.code(), e.getMessage())).join();

        } catch (OperationForbiddenException e) {
            log.warn("Got forbidden while processing {}", message, e);
            sendResponse(session,
                    Response.error(session.getId(), request.id(), HttpResponseStatus.FORBIDDEN.code(), e.getMessage()))
                            .join();

        } catch (TableControllerException e) {
            log.warn("Got processing error while processing {}", message, e);
            sendResponse(session, Response.error(session.getId(), request.id(),
                    HttpResponseStatus.INTERNAL_SERVER_ERROR.code(), e.getMessage())).join();
        }
    }

    private CompletableFuture<Void> notifyUsers(final Table table) {
        return CompletableFuture.allOf(table.allUsers().map(u -> sessionMap.get(u.getId())).filter(s -> s != null)
                .map(s -> sendResponse(s, responseSerializer.serialize(Response.table(s.getId(), table))))
                .toArray(CompletableFuture[]::new));
    }

    private CompletableFuture<Void> sendResponse(final Session session, final Response response) {
        return sendResponse(session, responseSerializer.serialize(response));
    }

    private CompletableFuture<Void> sendResponse(final Session session, final String message) {
        return CompletableFuture.supplyAsync(() -> {
            try {
                session.getBasicRemote().sendText(message);
                return null;
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }, executor);
    }
}
