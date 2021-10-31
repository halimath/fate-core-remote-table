import * as wecco from "@weccoframework/core"
import { Message, UpdateTable } from "../control"
import { Aspect, Player, Table } from "../models"

export interface AspectMessage {
    id: string
    name: string
}
export interface PlayerMessage {
    id: string
    name: string
    fatePoints: number
    aspects: Array<AspectMessage>
}

export interface TableMessage {
    id: string
    title: string
    gamemaster: string
    players: Array<PlayerMessage>
    aspects: Array<AspectMessage>
}

export class API {
    static connect(context: wecco.AppContext<Message>, userId: string): Promise<API> {
        return new Promise((resolve, reject) => {            
            const websocket = new WebSocket(`ws${document.location.protocol === "https:" ? "s" : ""}://${document.location.host}/user/${userId}`)
            websocket.onopen = () => {
                resolve(new API(context, websocket))
            }
        })
    }

    constructor (private readonly context: wecco.AppContext<Message>, private readonly websocket: WebSocket) {
        this.websocket.onmessage = this.handleMessage.bind(this)
    }

    public createTable(title: string) {
        this.websocket.send(JSON.stringify({
            type: "new-table",
            title: title,
        }))
    }

    public joinTable(id: string, name: string) {
        this.websocket.send(JSON.stringify({
            type: "join-table",
            tableId: id,
            name: name,
        }))
    }

    public updateFatePoints(tableId: string, playerId: string, fatePoints: number) {
        this.websocket.send(JSON.stringify({
            type: "update-fate-points",
            tableId: tableId,
            playerId: playerId,
            fatePoints: fatePoints,
        }))
    }

    public spendFatePoint(tableId: string) {
        this.websocket.send(JSON.stringify({
            type: "spend-fate-point",
            tableId: tableId,
        }))
    }

    public addAspect(tableId: string, name: string, playerId?: string) {
        this.websocket.send(JSON.stringify({
            type: "add-aspect",
            tableId: tableId,
            name: name,
            playerId: playerId,
        }))
    }

    public removeAspect(tableId: string, id: string) {
        this.websocket.send(JSON.stringify({
            type: "remove-aspect",
            tableId: tableId,
            id: id,
        }))
    }

    private handleMessage (evt: MessageEvent<string>) {
        try {
            const tableMessage: TableMessage = JSON.parse(evt.data)
            this.context.emit(new UpdateTable(convertTable(tableMessage)))
        } catch (e) {
            console.error(e)
        }
    }
}

function convertTable(msg: TableMessage): Table {
    return new Table(
        msg.id, 
        msg.title, 
        msg.gamemaster,
        msg.players.map(p => new Player(p.id, p.name, p.fatePoints, p.aspects.map(a => new Aspect(a.id, a.name)))), 
        msg.aspects.map(a => new Aspect(a.id, a.name)),
    )
}