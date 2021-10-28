import * as wecco from "@weccoframework/core"
import { API, TableMessage } from "../api"
import { Gamemaster, Model, Player, PlayerCharacter, Rating, Table } from "../models"

export class UpdateTable {
    readonly command = "update-table"

    constructor (public readonly table: TableMessage) {}
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

export class RollDice {
    readonly command = "roll-dice"

    constructor(public readonly rating: Rating) { }
}

export type Message = UpdateTable | NewTable | JoinTable | UpdatePlayerFatePoints | SpendFatePoint | RollDice

export class Controller {
    private api: API

    async update(model: Model, message: Message, context: wecco.AppContext<Message>): Promise<Model | typeof wecco.NoModelChange> {
        if (typeof this.api === "undefined") {
            this.api = await API.connect(context, model.userId)
        }

        switch (message.command) {
            case "update-table": 
                return this.convertTable(model, message.table)

            case "new-table":
                this.api.createTable(message.title)
                return wecco.NoModelChange

            case "join-table":
                this.api.joinTable(message.id, message.name)
                return wecco.NoModelChange

            case "update-fate-points":
                if (model instanceof Gamemaster) {
                    this.api.updateFatePoints(model.table.id, message.playerId, message.fatePoints)
                }
                return wecco.NoModelChange

            case "spend-fate-point":
                if (model instanceof PlayerCharacter) {
                    this.api.spendFatePoint(model.table.id)
                }
                return wecco.NoModelChange

            case "roll-dice":
                if ((model instanceof PlayerCharacter) || (model instanceof Gamemaster)) {
                    return model.roll(message.rating)
                }
        }

        return model
    }

    private convertTable(model: Model, msg: TableMessage): Model {
        const table = new Table(msg.id, msg.title, msg.players.map(p => new Player(p.id, p.name, p.fatePoints)))

        if (msg.gamemaster === model.userId) {
            return new Gamemaster(model.userId, table)
        }

        return new PlayerCharacter(model.userId, table)
    }
}
