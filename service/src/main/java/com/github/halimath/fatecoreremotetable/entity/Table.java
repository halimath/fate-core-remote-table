package com.github.halimath.fatecoreremotetable.entity;

import java.util.HashSet;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Stream;

import lombok.Builder;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
@Getter
@EqualsAndHashCode
@Builder
public class Table {
    private final String id;
    private final String title;
    private final User gameMaster;
    private final Set<Player> players = new HashSet<>();

    public void join (@NonNull final Player player) {
        players.add(player);
    }

    public void removePlayer(@NonNull final User user) {
        this.players.removeIf(p -> p.getUser().getId().equals(user.getId()));
    }

    public Optional<Player> findPlayer (@NonNull final String id) {
        return players.stream().filter(p -> p.getUser().getId().equals(id)).findFirst();
    }

    public Stream<User> allUsers() {
        return Stream.concat(Stream.of(gameMaster), players.stream().map(Player::getUser));
    }
}
