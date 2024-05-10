import { ApiClient, ApiError, Session as SessionDto } from "../../generated"
import { Aspect, Player, Session } from "../models"

const AuthTokenSessionStorageKey = "auth-token"

export async function createApiClient(): Promise<ApiClient> {
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

export class GamemasterApi {
    static async createSession(title: string): Promise<GamemasterApi> {
        const apiClient = await createApiClient()

        const sessionId = await apiClient.session.createSession({
            requestBody: {
                title: title,
            }
        })

        return new GamemasterApi(apiClient, sessionId)
    }

    constructor(
        private readonly apiClient: ApiClient,
        readonly sessionId: string,
    ) { }

    async getSession(): Promise<Session | null> {
        try {
            let dto = await this.apiClient.session.getSession({ id: this.sessionId })
            return convertTable(dto)
        } catch (e) {
            if (e instanceof ApiError) {
                if (e.status === 404) {
                    return null
                }
            }

            throw e
        }
    }

    async updateFatePoints(characterId: string, delta: number) {
        await this.apiClient.session.updateFatePoints({
            id: this.sessionId,
            characterId: characterId,
            requestBody: {
                fatePointsDelta: delta,
            }
        })
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
    }

    async removeAspect(id: string) {
        await this.apiClient.session.deleteAspect({
            id: this.sessionId,
            aspectId: id,
        })
    }
}

export class PlayerCharacterApi {
    static async joinGame(id: string, name: string): Promise<PlayerCharacterApi> {
        const apiClient = await createApiClient()

        const characterId = await apiClient.session.joinSession({
            id: id,
            requestBody: {
                name: name,
            }
        })

        return new PlayerCharacterApi(apiClient, id, characterId)
    }

    constructor(
        private readonly apiClient: ApiClient,
        readonly sessionId: string,
        readonly characterId: string,
    ) { }

    async getSession(): Promise<Session | null> {
        try {
            let dto = await this.apiClient.session.getSession({ id: this.sessionId })
            return convertTable(dto, this.characterId)
        } catch (e) {
            if (e instanceof ApiError) {
                if (e.status === 404) {
                    return null
                }
            }

            throw e
        }
    }

    async spendFatePoint() {
        await this.apiClient.session.updateFatePoints({
            id: this.sessionId,
            characterId: this.characterId!,
            requestBody: {
                fatePointsDelta: -1,
            }
        })
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