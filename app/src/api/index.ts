import * as wecco from "@weccoframework/core"
import { ApiClient, CreateCharacter, Session as SessionDto } from "../../generated"
import { Message, ReplaceScene } from "../control"
import { Aspect, Gamemaster, Player, PlayerCharacter, Session } from "../models"

abstract class ApiBase {
    protected constructor(
        protected readonly context: wecco.AppContext<Message>,
        protected readonly apiClient: ApiClient,
        readonly sessionId: string,
        private readonly modelCtor: new (table: Session) => Gamemaster | PlayerCharacter,
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

const AuthTokenSessionStorageKey = "auth-token"

async function createApiClient(): Promise<ApiClient> {
    let authToken = sessionStorage.getItem(AuthTokenSessionStorageKey)
    let apiClient: ApiClient

    if (authToken !== null) {
        apiClient = new ApiClient({
            BASE: "/api",
            TOKEN: authToken,
            CREDENTIALS: "include",
            WITH_CREDENTIALS: true,
        })
    } else {
        apiClient = new ApiClient({
            BASE: "/api",
        })
    }

    const token = await apiClient.authorization.createAuthToken()
    sessionStorage.setItem(AuthTokenSessionStorageKey, token)

    return new ApiClient({
        BASE: "/api",
        TOKEN: token,
        CREDENTIALS: "include",
        WITH_CREDENTIALS: true,
    })
}

export class GamemasterApi extends ApiBase {
    static async createSession(context: wecco.AppContext<Message>, title: string): Promise<GamemasterApi> {
        const apiClient = await createApiClient()

        const sessionId = await apiClient.session.createSession({
            requestBody: {
                title: title,
            }
        })

        return new GamemasterApi(context, apiClient, sessionId)
    }

    static async joinSession(context: wecco.AppContext<Message>, sessionId: string): Promise<GamemasterApi> {
        const apiClient = await createApiClient()

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

function convertTable(msg: SessionDto, characterId?: string): Session {
    return new Session(
        msg.id,
        msg.title,
        msg.ownerId,
        msg.characters.map(p => new Player(p.id, p.name, p.id === characterId, p.fatePoints, p.aspects.map(a => new Aspect(a.id, a.name)))),
        msg.aspects.map(a => new Aspect(a.id, a.name)),
    )
}