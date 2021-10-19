package com.github.halimath.fatetable.boundary.commands;

import com.github.halimath.fatetable.boundary.Command;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class UpdateFatePoints implements Command {
    private String type;
    private String tableId;
    private String playerId;
    private Integer fatePoints;
}
