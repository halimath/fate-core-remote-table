package com.github.halimath.fatecoreremotetable.boundary.websocket;

import java.io.IOException;
import java.util.Optional;

import javax.enterprise.context.ApplicationScoped;
import javax.websocket.Session;

import com.github.halimath.fatecoreremotetable.boundary.Response;
import com.github.halimath.fatecoreremotetable.boundary.ResponseSerializer;
import com.github.halimath.fatecoreremotetable.control.TableEvent;
import com.github.halimath.fatecoreremotetable.control.TableException;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import io.netty.handler.codec.http.HttpResponseStatus;
import io.quarkus.vertx.ConsumeEvent;
import io.smallrye.mutiny.Uni;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@ApplicationScoped
@RequiredArgsConstructor
@Slf4j
class TableEventConsumer {
    private final SessionRegistry sessionRegistry;
    private final ResponseSerializer responseSerializer;

    @ConsumeEvent(TableEvent.Updated.LISTENER_NAME)
    void handle(@NonNull final TableEvent.Updated event) {
        notifyUsers(event.table())
            .subscribe().with(
                ignored -> log.debug("Table update published"),
                t -> log.warn("Got error while publishing table update", t)
            );
    }

    @ConsumeEvent(TableEvent.Closed.LISTENER_NAME)
    void handle(@NonNull final TableEvent.Closed event) {
        Uni.combine().all().unis(event.users().stream().map(this::disconnectUser).toList()).discardItems()
                .subscribe().with(
                        ignored -> log.debug("Orphan users disconnected"),
                        t -> log.warn("Failed to close websockets", t));
    }

    @ConsumeEvent(TableEvent.Error.LISTENER_NAME)
    void handle(@NonNull final TableEvent.Error event) {
        final var session = sessionRegistry.get(event.user()).orElse(null);
        if (session == null) {
            return;
        }

        final Response res;

        if (event.t() instanceof TableException.TableNotFound e) {
            log.warn("Table not found: {}", e.getMessage());
            res = Response.error(session.getId(), HttpResponseStatus.NOT_FOUND.code(), e.getMessage());

        } else if (event.t() instanceof TableException.PlayerNotFound e) {
            log.warn("Player not found: {}", e.getMessage());
            res = Response.error(session.getId(), HttpResponseStatus.PRECONDITION_FAILED.code(),
                    e.getMessage());

        } else if (event.t() instanceof TableException.OperationForbidden e) {
            log.warn("Operation forbidden", e);
            res = Response.error(session.getId(), HttpResponseStatus.FORBIDDEN.code(), e.getMessage());

        } else {
            log.warn("Got unexpected error while processing command", event.t());
            res = Response.error(session.getId(), HttpResponseStatus.INTERNAL_SERVER_ERROR.code(),
                    event.t().getMessage());
        }

        sendResponse(session, res) //
                .subscribe().with(ignored -> log.debug("command dispatch complete"));        
    }

    private Uni<Void> disconnectUser(final User user) {
        return sessionRegistry.get(user)
                .<Uni<Void>>map(s -> {
                    try {
                        s.close();
                        return Uni.createFrom().item((Void) null);
                    } catch (IOException e) {
                        return Uni.createFrom().failure(e);
                    }
                })
                .orElseGet(() -> Uni.createFrom().item((Void) null));
    }

    private Uni<Void> notifyUsers(final Table table) {
        log.debug("Notifying users of {}", table.getId());

        return Uni.combine().all().unis(table.allUsers().map(u -> sessionRegistry.get(u))
                .filter(Optional::isPresent)
                .map(Optional::get)
                .map(s -> sendResponse(s, responseSerializer.serialize(Response.table(s.getId(), table)))).toList())
                .discardItems();
    }

    private Uni<Void> sendResponse(final Session session, final Response response) {
        return sendResponse(session, responseSerializer.serialize(response));
    }

    private Uni<Void> sendResponse(final Session session, final String message) {
        return Uni.createFrom().future(session.getAsyncRemote().sendText(message));
    }

}
