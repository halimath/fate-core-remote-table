import * as wecco from "@weccoframework/core"
import { API } from "../api"
import { Gamemaster, Home, Model, Notification, PlayerCharacter, Rating, Scene } from "../models"
import { m } from "../utils/i18n"

export class ReplaceScene {
    readonly command = "replace-scene"

    constructor (public readonly scene: Scene) {}
}

export class PostNotification {
    readonly command = "post-notification"

    public readonly notifications: Array<Notification>
    
    constructor (...notifications: Array<Notification>) {
        this.notifications = notifications
    }
}

export class NewTable {
    readonly command = "new-table"

    constructor(public readonly title: string) { }
}

export class JoinTable {
    readonly command = "join-table"

    constructor(public readonly id: string, public readonly name: string) {}
}

export class UpdatePlayerFatePoints {
    readonly command = "update-fate-points"

    constructor (public readonly playerId: string, public readonly fatePoints: number) {}
}

export class SpendFatePoint {
    readonly command = "spend-fate-point"
}

export class AddAspect {
    readonly command = "add-aspect"

    constructor (public readonly name: string, public readonly targetPlayerId?: string) {}
}

export class RemoveAspect {
    readonly command = "remove-aspect"

    constructor (public readonly id: string) {}
}

export class RollDice {
    readonly command = "roll-dice"

    constructor(public readonly rating: Rating) { }
}

export class TableClosed {
    readonly command = "table-closed"
}

export type Message = ReplaceScene | 
    PostNotification |
    NewTable | 
    JoinTable | 
    UpdatePlayerFatePoints | 
    SpendFatePoint | 
    AddAspect |
    RemoveAspect |
    RollDice |
    TableClosed

export class Controller {
    private api: API
    
    async update(model: Model, message: Message, context: wecco.AppContext<Message>): Promise<Model | typeof wecco.NoModelChange> {
        if (typeof this.api === "undefined") {
            this.api = await API.connect(context)
        }

        switch (message.command) {
            case "replace-scene":
                return new Model(model.versionInfo, message.scene)

            case "post-notification":
                return new Model(model.versionInfo, model.scene, ...message.notifications)
            
            case "roll-dice":
                return new Model(model.versionInfo, model.scene.roll(message.rating))

            case "table-closed":
                return new Model(model.versionInfo, new Home(), new Notification(m("tableClosed.message")))

            case "new-table":
                this.api.createTable(message.title)
                break

            case "join-table":
                this.api.joinTable(message.id, message.name)
                break

            case "update-fate-points":
                if (model.scene instanceof Gamemaster) {
                    this.api.updateFatePoints(message.playerId, message.fatePoints)
                }
                break

            case "spend-fate-point":
                if (model.scene instanceof PlayerCharacter) {
                    this.api.spendFatePoint()
                }
                break

            case "add-aspect":
                if (model.scene instanceof Gamemaster) {
                    this.api.addAspect(message.name, message.targetPlayerId)
                }
                break

            case "remove-aspect":
                if (model.scene instanceof Gamemaster) {
                    this.api.removeAspect(message.id)
                }
                break
        }
        
        return model.notifications.length === 0 ? wecco.NoModelChange : model.pruneNotifications()
    }
}
