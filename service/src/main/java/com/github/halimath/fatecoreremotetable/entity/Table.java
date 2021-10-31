package com.github.halimath.fatecoreremotetable.entity;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Stream;

import lombok.Data;
import lombok.NonNull;

@Data
public class Table {
    private final String id;
    private final String title;
    private final User gameMaster;
    private final Set<Player> players = new HashSet<>();
    private final List<Aspect> aspects = new ArrayList<>();

    public void join(@NonNull final Player player) {
        players.add(player);
    }

    public void removePlayer(@NonNull final User user) {
        this.players.removeIf(p -> p.getUser().getId().equals(user.getId()));
    }

    public Optional<Player> findPlayer(@NonNull final String id) {
        return players.stream().filter(p -> p.getUser().getId().equals(id)).findFirst();
    }

    public Stream<User> allUsers() {
        return Stream.concat(Stream.of(gameMaster), players.stream().map(Player::getUser));
    }

    public void addAspect(@NonNull final Aspect aspect) {
        this.aspects.add(aspect);
    }

    public void removeAspect(@NonNull final String id) {
        for (var a: aspects) {
            if (a.getId().equals(id)) {
                aspects.remove(a);
                return;
            }
        }

        for (var p: players) {
            if (p.removeAspect(id)) {
                return;
            }
        }
    }
}
