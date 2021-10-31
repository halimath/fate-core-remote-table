package com.github.halimath.fatecoreremotetable.entity;

import java.util.ArrayList;
import java.util.List;

import lombok.Data;
import lombok.NonNull;

@Data
public class Player {
    private final User user;
    private final String name;
    private Integer fatePoints = 0;
    private final List<Aspect> aspects = new ArrayList<>();

    public void addAspect(@NonNull final Aspect aspect) {
        this.aspects.add(aspect);
    }
    
    public boolean removeAspect(@NonNull final String id) {
        for (var a: aspects) {
            if (a.getId().equals(id)) {
                aspects.remove(a);
                return true;
            }
        }        

        return false;
    }
}
