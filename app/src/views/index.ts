import * as wecco from "@weccoframework/core"
import { Message } from "../control"
import { GamemasterScene, HomeScene, Model, PlayerCharacterScene } from "../models"
import { showNotifications } from "./widgets/notification"
import { gamemaster } from "./scenes/gamemaster"
import { home } from "./scenes/home"
import { player } from "./scenes/player"

// We must import this component here to define the custom element
import "./components/skillcheck"

/**
 * `root` is the root view function executed by the wecco framework to apply model changes.
 * The main purpose is to dispatch based on the model's wrapped scene and show any notifications.
 */
export function root({model, emit}: wecco.ViewContext<Model, Message>): wecco.ElementUpdate {
    showNotifications(model.notifications)

    if (model.scene instanceof HomeScene) {
        return home(model.versionInfo, model.scene, emit)
    }

    if (model.scene instanceof PlayerCharacterScene) {
        return player(model.versionInfo, model.scene, emit)
    }

    if (model.scene instanceof GamemasterScene) {
        return gamemaster(model.versionInfo, model.scene, emit)
    }

    return "Unknown model"
}
