package com.github.halimath.fatetable.boundary;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;
import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.Session;
import javax.websocket.server.PathParam;
import javax.websocket.server.ServerEndpoint;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.halimath.fatetable.control.TableController;
import com.github.halimath.fatetable.control.TableController.TableControllerException;
import com.github.halimath.fatetable.entity.Table;
import com.github.halimath.fatetable.entity.User;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ServerEndpoint("/user/{id}")
@ApplicationScoped
@Slf4j
@RequiredArgsConstructor
public class UserWebSocket {
    private final Map<String, UserAndSession> sessionMap = new ConcurrentHashMap<>();
    private final ObjectMapper objectMapper = new ObjectMapper();
    @Inject
    private final TableController tableController;
    @Inject
    private final CommandDispatcher commandDispatcher;

    @OnOpen
    public void onOpen(final Session session, @PathParam("id") final String id) {
        log.info("onOpen: id={}", id);
        final var user = new User(id);
        sessionMap.put(id, new UserAndSession(user, session));
    }

    @OnClose
    public void onClose(final Session session, @PathParam("id") final String id) {
        log.info("onClose: id={}", id);
        final var user = sessionMap.get(id).user;        
        sessionMap.remove(id);
        tableController.disconnect(user).ifPresent(this::notifyUsers);        
    }

    @OnError
    public void onError(final Session session, @PathParam("id") final String id, final Throwable throwable) {
        log.warn("onError: id={}", id, throwable);
    }

    @OnMessage
    public void onMessage(final String message, @PathParam("id") final String id) {
        try {
            log.debug("onMessage: id={} message={}", id, message);
            final var user = sessionMap.get(id).user;
            final var command = objectMapper.readValue(message, Command.class);
            notifyUsers(commandDispatcher.dispatchCommand(user, command));
        } catch (IOException e) {
            log.error("onMessage: failed to parse command: id={} message={}", id, message, e);
        } catch (TableControllerException e) {
            log.warn("Table not found while processing {}", message, e);
        }
    }

    private void notifyUsers(final Table table) {
        try {
            final var msg = Message.fromEntity(table);
            final var json = objectMapper.writeValueAsString(msg);

            table.allUsers().forEach(u -> notifyUser(u, json));
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    private void notifyUser(final User user, final String message) {
        sessionMap.get(user.getId()).session.getAsyncRemote().sendText(message);
    }

    private static record UserAndSession(User user, Session session) {
    }
}
