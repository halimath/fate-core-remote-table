package com.github.halimath.fatecoreremotetable.boundary;

import java.util.List;

import com.github.halimath.fatecoreremotetable.entity.Table;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NonNull;

@AllArgsConstructor
@Builder
@Getter
public class Message {
    public static Message fromEntity(@NonNull final Table table) {
        return new Message(
            table.getId(), 
            table.getTitle(), 
            table.getGameMaster().getId(), 
            table.getPlayers().stream()
                .map(p -> new Message.Player(
                    p.getUser().getId(), 
                    p.getName(), 
                    p.getFatePoints(), 
                    p.getAspects().stream().map(a -> new Aspect(a.getId(), a.getName())).toList())).toList(),
            table.getAspects().stream().map(a -> new Aspect(a.getId(), a.getName())).toList()
            );
    }

    private final String id;
    private final String title;
    private final String gamemaster;
    private final List<Player> players;
    private final List<Aspect> aspects;

    public static record Player(String id, String name, Integer fatePoints, List<Aspect> aspects) {
    }

    public static record Aspect(String id, String name) {}
}
