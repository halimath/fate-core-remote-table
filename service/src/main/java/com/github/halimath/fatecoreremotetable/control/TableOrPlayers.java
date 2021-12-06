package com.github.halimath.fatecoreremotetable.control;

import java.util.Set;

import com.github.halimath.fatecoreremotetable.entity.Player;
import com.github.halimath.fatecoreremotetable.entity.Table;

public record TableOrPlayers(Table table, Set<Player> players) {
}
