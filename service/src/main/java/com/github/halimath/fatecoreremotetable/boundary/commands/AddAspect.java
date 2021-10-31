package com.github.halimath.fatecoreremotetable.boundary.commands;

import java.util.Optional;

import com.github.halimath.fatecoreremotetable.boundary.Command;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class AddAspect implements Command {
    private String type;
    private String tableId;
    private String name;
    private String playerId;

    public Optional<String> getOptionalPlayerId() {
        return Optional.ofNullable(playerId);
    }
}
