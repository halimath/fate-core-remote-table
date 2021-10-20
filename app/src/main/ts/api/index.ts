import * as wecco from "@weccoframework/core"
import { Message, UpdateTable } from "../control"

export interface PlayerMessage {
    id: string
    name: string
    fatePoints: number
}

export interface TableMessage {
    id: string
    title: string
    gamemaster: string
    players: Array<PlayerMessage>
}

export class API {
    static connect(context: wecco.AppContext<Message>, userId: string, baseUrl: string = "ws://localhost:8080"): Promise<API> {
        return new Promise((resolve, reject) => {
            const websocket = new WebSocket(`${baseUrl}/user/${userId}`)
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

    private handleMessage (evt: MessageEvent<string>) {
        try {
            const table: TableMessage = JSON.parse(evt.data)
            this.context.emit(new UpdateTable(table))
        } catch (e) {
            console.error(e)
        }
    }
}