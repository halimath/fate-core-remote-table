package com.github.halimath.fatecoreremotetable.boundary.dto;

import java.util.Optional;

import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeInfo.As;
import com.fasterxml.jackson.annotation.JsonTypeInfo.Id;

@JsonTypeInfo(use = Id.NAME, include = As.PROPERTY, property = "type")
@JsonSubTypes({ 
    @JsonSubTypes.Type(value = Command.Create.class, name = "create"),
    @JsonSubTypes.Type(value = Command.Join.class, name = "join"),
    @JsonSubTypes.Type(value = Command.AddAspect.class, name = "add-aspect"),
    @JsonSubTypes.Type(value = Command.RemoveAspect.class, name = "remove-aspect"),
    @JsonSubTypes.Type(value = Command.UpdateFatePoints.class, name = "update-fate-points"),
    @JsonSubTypes.Type(value = Command.SpendFatePoint.class, name = "spend-fate-point")
})
public interface Command {
  String type();

  public record Create(String type, String title) implements Command {}

  public record Join (String type, String tableId, String name) implements Command {}
  
  public record AddAspect(String type, String name, String playerId) implements Command {
    public Optional<String> optionalPlayerId() {
      return Optional.ofNullable(playerId);
    }
  }
  
  public record RemoveAspect (String type, String id) implements Command {}

  public record SpendFatePoint (String type) implements Command {}

  public record UpdateFatePoints (String type, String playerId, Integer fatePoints) implements Command {}
}
