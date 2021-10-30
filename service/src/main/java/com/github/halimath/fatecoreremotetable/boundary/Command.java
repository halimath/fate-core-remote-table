package com.github.halimath.fatecoreremotetable.boundary;

import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeInfo.As;
import com.fasterxml.jackson.annotation.JsonTypeInfo.Id;
import com.github.halimath.fatecoreremotetable.boundary.commands.JoinTable;
import com.github.halimath.fatecoreremotetable.boundary.commands.NewTable;
import com.github.halimath.fatecoreremotetable.boundary.commands.SpendFatePoint;
import com.github.halimath.fatecoreremotetable.boundary.commands.UpdateFatePoints;

@JsonTypeInfo(use = Id.NAME, include = As.PROPERTY, property = "type")
@JsonSubTypes({ 
    @JsonSubTypes.Type(value = NewTable.class, name = "new-table"),
    @JsonSubTypes.Type(value = JoinTable.class, name = "join-table"), 
    @JsonSubTypes.Type(value = UpdateFatePoints.class, name = "update-fate-points"), 
    @JsonSubTypes.Type(value = SpendFatePoint.class, name = "spend-fate-point") 
})
public interface Command {
  String getType();
}
