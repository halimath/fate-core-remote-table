package com.github.halimath.fatecoreremotetable.control;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

import com.github.halimath.fatecoreremotetable.entity.Aspect;
import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;
import com.github.halimath.fatecoreremotetable.entity.User;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class TableServiceTest {

    private TableRepository repository;
    private TableService service;

    @BeforeEach
    void init() {
        repository = new TableRepository();
        service = new TableService(repository);
    }

    @Test
    void createTable_shouldCreateTable() {
        final var user = new User("test");
        final var got = service.create(user, "1", "Test").await().indefinitely();

        assertEquals(new Table("1", "Test", user), got);
    }

    @Test
    void createTable_shouldNotCreateTableWhenTableAlreadyExists() {
        final var user = new User("test");
        service.create(user, "1", "Already exists").await().indefinitely();
        assertThrows(TableException.Conflict.class,
                () -> service.create(user, "1", "Test").await().indefinitely());
    }

    @Test
    void createTable_shouldNotCreateTableWhenUserIsAlreadyGamemaster() {
        final var user = new User("test");
        service.create(user, "1", "Already exists").await().indefinitely();
        assertThrows(TableException.OperationForbidden.class,
                () -> service.create(user, "2", "Test").await().indefinitely());
    }

    @Test
    void createTable_shouldNotCreateTableWhenUserIsAlreadyAPlayer() {
        final var user = new User("test");
        final var table = new Table("1", "", new User("gm"));
        table.join(new Player(user, "Player Unknown"));
        repository.save(table).await().indefinitely();

        assertThrows(TableException.OperationForbidden.class,
                () -> service.create(user, "2", "Test").await().indefinitely());
    }

    @Test
    void join_shouldJoinPlayer() {
        final var table = new Table("1", "", new User("gm"));

        repository.save(table).await().indefinitely();

        final var user = new User("test");
        var got = service.join(user, table.getId(), "Player Unknown").await().indefinitely();

        assertEquals(1, got.getPlayers().size());
    }

    @Test
    void join_shouldNotJoinPlayerWhenTableIsNotFound() {
        assertThrows(TableException.TableNotFound.class,
                () -> service.join(new User("1"), "not found", "Player Unknown").await().indefinitely());
    }

    @Test
    void join_shouldNotJoinPlayerWhenUserIsGamemaster() {
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);

        repository.save(table).await().indefinitely();

        assertThrows(TableException.OperationForbidden.class,
                () -> service.join(gm, table.getId(), "Player Unknown").await().indefinitely());
    }

    @Test
    void join_shouldNotJoinPlayerWhenUserJoinedAnotherGame() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("gm", "", gm);
        table.join(new Player(user, "Player Unknwon"));
        repository.save(table).await().indefinitely();

        repository.save(new Table("2", "", gm)).await().indefinitely();

        assertThrows(TableException.OperationForbidden.class,
                () -> service.join(user, "2", "Player Unknown").await().indefinitely());
    }

    @Test
    void updateFatePoints_shouldUpdateFatePoints() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        table.join(new Player(user, "Player Unknwon"));
        repository.save(table).await().indefinitely();

        var got = service.updateFatePoints(gm, "1", user.getId(), 2).await().indefinitely();

        assertEquals(2, got.findPlayer(user.getId()).get().getFatePoints());
    }

    @Test
    void updateFatePoints_shouldNotUpdateFatePointsWhenTableDoesNotExist() {
        final var user = new User("user");
        final var gm = new User("gm");

        assertThrows(TableException.TableNotFound.class,
                () -> service.updateFatePoints(gm, "not-found", user.getId(), 2).await().indefinitely());
    }

    @Test
    void updateFatePoints_shouldNotUpdateFatePointsWhenUserIsNotTheGameMaster() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        table.join(new Player(user, "Player Unknwon"));
        repository.save(table).await().indefinitely();

        assertThrows(TableException.OperationForbidden.class,
                () -> service.updateFatePoints(user, "1", user.getId(), 2).await().indefinitely());
    }

    @Test
    void spendFatePoint_shouldReduceFatePoints() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        final var player = new Player(user, "Player Unknwon");
        table.join(player);
        player.setFatePoints(2);
        repository.save(table).await().indefinitely();

        final var got = service.spendFatePoint(user, "1").await().indefinitely();

        assertEquals(1, got.findPlayer("user").get().getFatePoints());
    }

    @Test
    void spendFatePoint_shouldNotReduceFatePointsWhenTableIsNotFound() {
        final var user = new User("user");
        assertThrows(TableException.TableNotFound.class,
                () -> service.spendFatePoint(user, "not-found").await().indefinitely());
    }

    @Test
    void spendFatePoint_shouldNotReduceFatePointsWhenPlayerIsNotPartOfTheTable() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        repository.save(table).await().indefinitely();

        assertThrows(TableException.PlayerNotFound.class,
                () -> service.spendFatePoint(user, "1").await().indefinitely());
    }

    @Test
    void addAspect_shouldAddGlobalAspect() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        table.join(new Player(user, "Player Unknown"));
        repository.save(table).await().indefinitely();

        final var got = service.addAspect(gm, "1", "test").await().indefinitely();

        assertEquals(1, got.getAspects().size());
    }

    @Test
    void addAspect_shouldAddPlayerAspect() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        table.join(new Player(user, "Player Unknown"));
        repository.save(table).await().indefinitely();

        final var got = service.addAspect(gm, "1", "test", user).await().indefinitely();

        assertEquals(0, got.getAspects().size());
        assertEquals(1, got.findPlayer(user.getId()).get().getAspects().size());
    }

    @Test
    void addAspect_shouldNotAddAspectWhenTableIsNotFound() {
        final var gm = new User("gm");

        assertThrows(TableException.TableNotFound.class,
                () -> service.addAspect(gm, "2", "test").await().indefinitely());
    }

    @Test
    void removeAspect_shouldRemoveGlobalAspect() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        table.join(new Player(user, "Player Unknown"));
        table.addAspect(new Aspect("1", "Test"));
        repository.save(table).await().indefinitely();

        final var got = service.removeAspect(gm, "1", "1").await().indefinitely();

        assertEquals(0, got.getAspects().size());
    }

    @Test
    void removeAspect_shouldPlayerGlobalAspect() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        final var player = new Player(user, "Player Unknown");
        player.addAspect(new Aspect("1", "Test"));
        table.join(player);
        repository.save(table).await().indefinitely();

        final var got = service.removeAspect(gm, "1", "1").await().indefinitely();

        assertEquals(0, got.findPlayer(user.getId()).get().getAspects().size());
    }

    @Test
    void leave_shouldRemovePlayer() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        final var player = new Player(user, "Player Unknown");
        table.join(player);
        repository.save(table).await().indefinitely();

        final var got = service.leave(user).await().indefinitely();

        assertNull(got.players());
        assertEquals(0, got.table().getPlayers().size());
    }

    @Test
    void leave_shouldRemoveGamemaster() {
        final var user = new User("user");
        final var gm = new User("gm");
        final var table = new Table("1", "", gm);
        final var player = new Player(user, "Player Unknown");
        table.join(player);
        repository.save(table).await().indefinitely();

        final var got = service.leave(gm).await().indefinitely();

        assertNull(got.table());
        assertEquals(1, got.players().size());
        assertNull(repository.findById("1").await().indefinitely());
    }
}
