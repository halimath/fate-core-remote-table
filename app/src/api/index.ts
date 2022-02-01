import * as wecco from "@weccoframework/core"
import { ApiClient, CreateCharacter, Session } from "../../generated"
import { Message, ReplaceScene } from "../control"
import { Aspect, Gamemaster, Player, PlayerCharacter, Table } from "../models"

abstract class ApiBase {
    protected constructor(
        protected readonly context: wecco.AppContext<Message>,
        protected readonly apiClient: ApiClient,
        protected readonly sessionId: string,
        private readonly modelCtor: new (table: Table) => Gamemaster | PlayerCharacter,
        protected readonly characterId?: string,
    ) {
        this.requestUpdate()
        setInterval(this.requestUpdate.bind(this), 5000)
    }

    protected async requestUpdate() {
        const session = await this.apiClient.session.getSession({
            id: this.sessionId,
        })

        this.context.emit(new ReplaceScene(new this.modelCtor(convertTable(session, this.characterId))))
    }
}

async function createApiClient(): Promise<ApiClient> {
    const c = new ApiClient({
        BASE: "/api",
    })

    const token = await c.authorization.createAuthToken()

    return new ApiClient({
        BASE: "/api",
        CREDENTIALS: "include",
        WITH_CREDENTIALS: true,
        TOKEN: token,
    })
}

export class GamemasterApi extends ApiBase {
    static async createGame(context: wecco.AppContext<Message>, title: string): Promise<GamemasterApi> {
        const apiClient = await createApiClient()

        const sessionId = await apiClient.session.createSession({
            requestBody: {
                title: title,
            }
        })

        return new GamemasterApi(context, apiClient, sessionId)
    }

    private constructor(
        context: wecco.AppContext<Message>,
        unauthorizedApiClient: ApiClient,
        sessionId: string,
    ) {
        super(context, unauthorizedApiClient, sessionId, Gamemaster)
    }

    async updateFatePoints(characterId: string, delta: number) {
        await this.apiClient.session.updateFatePoints({
            id: this.sessionId,
            characterId: characterId,
            requestBody: {
                fatePointsDelta: delta,
            }
        })
        this.requestUpdate()
    }

    async addAspect(name: string, playerId?: string) {
        if (playerId) {
            await this.apiClient.session.createCharacterAspect({
                id: this.sessionId,
                characterId: playerId,
                requestBody: {
                    name: name,
                }
            })

        } else {
            await this.apiClient.session.createAspect({
                id: this.sessionId,
                requestBody: {
                    name: name,
                }
            })
        }

        this.requestUpdate()
    }

    async removeAspect(id: string) {
        await this.apiClient.session.deleteAspect({
            id: this.sessionId,
            aspectId: id,
        })
        this.requestUpdate()
    }
}

export class PlayerCharacterApi extends ApiBase {
    static async joinGaim(context: wecco.AppContext<Message>, id: string, name: string): Promise<PlayerCharacterApi> {
        const apiClient = await createApiClient()

        const characterId = await apiClient.session.createCharacter({
            id: id,
            requestBody: {
                name: name,
                type: CreateCharacter.type.PC,
            }
        })

        return new PlayerCharacterApi(context, apiClient, id, characterId)
    }

    private constructor(
        context: wecco.AppContext<Message>,
        apiClient: ApiClient,
        sessionId: string,
        characterId: string,
    ) {
        super(context, apiClient, sessionId, PlayerCharacter, characterId)
    }

    async spendFatePoint() {
        await this.apiClient.session.updateFatePoints({
            id: this.sessionId,
            characterId: this.characterId!,
            requestBody: {
                fatePointsDelta: -1,
            }
        })
        this.requestUpdate()
    }
}

function convertTable(msg: Session, characterId?: string): Table {
    return new Table(
        msg.id,
        msg.title,
        msg.ownerId,
        msg.characters.map(p => new Player(p.id, p.name, p.id === characterId, p.fatePoints, p.aspects.map(a => new Aspect(a.id, a.name)))),
        msg.aspects.map(a => new Aspect(a.id, a.name)),
    )
}