import * as wecco from "@weccoframework/core"
import { Message, PostNotification, ReplaceScene } from "../control"
import { Aspect, Gamemaster, Notification, Player, PlayerCharacter, Table } from "../models"
import { v4} from "uuid"

export interface AspectDto {
    id: string
    name: string
}
export interface PlayerDto {
    id: string
    name: string
    fatePoints: number
    aspects: Array<AspectDto>
}

export interface TableDto {
    id: string
    title: string
    gamemaster: string
    players: Array<PlayerDto>
    aspects: Array<AspectDto>
}


export interface ErrorDto {
    requestId?: string
    code: number
    reason: string
}
export interface Update {
    type: "table" | "error"
    self: string
    error?: ErrorDto
    table?: TableDto
}

export class API {
    static connect(context: wecco.AppContext<Message>): Promise<API> {
        return new Promise(resolve => {            
            const websocket = new WebSocket(`ws${document.location.protocol === "https:" ? "s" : ""}://${document.location.host}/table`)
            websocket.onopen = () => {
                resolve(new API(context, websocket))
            }
        })
    }

    constructor (private readonly context: wecco.AppContext<Message>, private readonly websocket: WebSocket) {
        this.websocket.onmessage = this.handleMessage.bind(this)
    }

    public createTable(title: string) {
        this.sendCommand({
            type: "create",
            title: title,
        })
    }

    public joinTable(tableId: string, name: string) {
        this.sendCommand({
            type: "join",
            tableId: tableId,
            name: name,
        })
    }

    public updateFatePoints(playerId: string, fatePoints: number) {
        this.sendCommand({
            type: "update-fate-points",
            playerId: playerId,
            fatePoints: fatePoints,
        })
    }

    public spendFatePoint() {
        this.sendCommand({
            type: "spend-fate-point",
        })
    }

    public addAspect(name: string, playerId?: string) {
        this.sendCommand({
            type: "add-aspect",
            name: name,
            playerId: playerId,
        })
    }

    public removeAspect(id: string) {
        this.sendCommand({
            type: "remove-aspect",
            id: id,
        })
    }

    private sendCommand (command: any) {
        const request = {
            id: v4(),
            command: command,
        }

        this.websocket?.send(JSON.stringify(request))
    }

    private handleMessage (evt: MessageEvent<string>) {
        try {
            const update: Update = JSON.parse(evt.data)
            
            if (update.table) {
                const table = convertTable(update.table)
                if (update.table.gamemaster === update.self) {
                    this.context.emit(new ReplaceScene(new Gamemaster(update.self, table)))
                    return
                }
            
                this.context.emit(new ReplaceScene(new PlayerCharacter(update.self, table)))
                return
            }

            this.context.emit(new PostNotification(new Notification(update.error?.reason ?? "Error", "error")))
        } catch (e) {
            console.error(e)
        }
    }
}

function convertTable(msg: TableDto): Table {
    return new Table(
        msg.id, 
        msg.title, 
        msg.gamemaster,
        msg.players.map(p => new Player(p.id, p.name, p.fatePoints, p.aspects.map(a => new Aspect(a.id, a.name)))), 
        msg.aspects.map(a => new Aspect(a.id, a.name)),
    )
}