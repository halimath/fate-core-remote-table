import * as wecco from "@weccoframework/core"
import { Message } from "../control"
import { Gamemaster, Home, Model, Player, PlayerCharacter } from "../models"
import { gamemaster } from "./scenes/gamemaster"
import { home } from "./scenes/home"
import { player } from "./scenes/player"

export function root (model: Model, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    if (model instanceof Home) {
        return home(model, context)
    }

    if (model instanceof PlayerCharacter) {
        return player(model, context)
    }

    if (model instanceof Gamemaster) {
        return gamemaster(model, context)
    }

    return "Unknown model"
}
