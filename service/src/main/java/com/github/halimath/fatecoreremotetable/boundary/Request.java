package com.github.halimath.fatecoreremotetable.boundary;

import java.util.Optional;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeInfo.As;
import com.fasterxml.jackson.annotation.JsonTypeInfo.Id;

import lombok.NonNull;

public record Request(String id, String tableId, Command command) {
    public static record User (String id) {}

    @JsonTypeInfo(use = Id.NAME, include = As.PROPERTY, property = "type")
    @JsonSubTypes({ @JsonSubTypes.Type(value = Command.Create.class, name = "create"),
            @JsonSubTypes.Type(value = Command.Join.class, name = "join"),
            @JsonSubTypes.Type(value = Command.AddAspect.class, name = "add-aspect"),
            @JsonSubTypes.Type(value = Command.RemoveAspect.class, name = "remove-aspect"),
            @JsonSubTypes.Type(value = Command.UpdateFatePoints.class, name = "update-fate-points"),
            @JsonSubTypes.Type(value = Command.SpendFatePoint.class, name = "spend-fate-point") })
    interface Command {
        
        record Create(@JsonProperty(required = true) @NonNull String title) implements Command {
        }

        record Join(@JsonProperty(required = true) @NonNull String name) implements Command {
        }

        record AddAspect(@JsonProperty(required = true) @NonNull String name, String playerId) implements Command {
            Optional<String> optionalPlayerId() {
                return Optional.ofNullable(playerId);
            }
        }

        record RemoveAspect(@JsonProperty(required = true) @NonNull String id) implements Command {
        }

        record SpendFatePoint() implements Command {
        }

        record UpdateFatePoints(@JsonProperty(required = true) @NonNull String playerId,
                @JsonProperty(required = true) @NonNull Integer fatePoints) implements Command {
        }
    }
}
