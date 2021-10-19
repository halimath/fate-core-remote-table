package com.github.halimath.fatetable.entity;

import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import lombok.Setter;

@RequiredArgsConstructor
@Getter
@Setter
@EqualsAndHashCode
public class Player {
    private final User user;
    private final String name;
    private Integer fatePoints = 0;
}
