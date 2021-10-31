package com.github.halimath.fatecoreremotetable.boundary.commands;

import com.github.halimath.fatecoreremotetable.boundary.Command;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class RemoveAspect implements Command {
    private String type;
    private String tableId;
    private String id;    
}
