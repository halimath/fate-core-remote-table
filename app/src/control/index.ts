import * as wecco from "@weccoframework/core"
import { GamemasterApi, PlayerCharacterApi, createApiClient } from "../api"
import { GamemasterScene, HomeScene, Model, Notification, PlayerCharacterScene, Scene } from "../models"
import { m } from "../utils/i18n"

export class ReplaceScene {
    readonly command = "replace-scene"

    public readonly notifications: Array<Notification>

    constructor(public readonly scene: Scene, ...notifications: Array<Notification>) { 
        this.notifications = notifications
    }
}

export class PostNotification {
    readonly command = "post-notification"

    public readonly notifications: Array<Notification>

    constructor(...notifications: Array<Notification>) {
        this.notifications = notifications
    }
}

export class NewSession {
    readonly command = "new-session"

    constructor(public readonly title: string) { }
}

export class Rejoin {
    readonly command = "rejoin"

    constructor(public readonly sessionId: string) { }
}

export class JoinAsPlayer {
    readonly command = "join-as-player"

    constructor(public readonly id: string, public readonly name: string) { }
}

export class UpdatePlayerFatePoints {
    readonly command = "update-fate-points"

    constructor(public readonly playerId: string, public readonly fatePoints: number) { }
}

export class SpendFatePoint {
    readonly command = "spend-fate-point"
}

export class AddAspect {
    readonly command = "add-aspect"

    constructor(public readonly name: string, public readonly targetPlayerId?: string) { }
}

export class RemoveAspect {
    readonly command = "remove-aspect"

    constructor(public readonly id: string) { }
}

export class SessionClosed {
    readonly command = "session-closed"
}

export type Message = ReplaceScene |
    PostNotification |
    NewSession |
    Rejoin |
    JoinAsPlayer |
    UpdatePlayerFatePoints |
    SpendFatePoint |
    AddAspect |
    RemoveAspect |
    SessionClosed

export class Controller {
    private api: GamemasterApi | PlayerCharacterApi | null = null

    private updateIntervalHandle: number | null = null
    private inUpdate = false

    private clearUpdate() {
        if (this.updateIntervalHandle !== null) {
            clearInterval(this.updateIntervalHandle)
        }
        this.inUpdate = false
    }
            
    private scheduleUpdates(emit: wecco.MessageEmitter<Message>) {
        const requestAndEmitUpdate = async () => {
            if (!this.api) {
                return
            }

            if (this.inUpdate) {
                return
            }

            try {
                this.inUpdate = true
                const session = await withMaxRetries(this.api.getSession.bind(this.api))
                if (session === null) {
                    // Session has been removed from the server.                
                    emit(new SessionClosed())
                    return
                }
                
                let scene: Scene
                
                if (this.api instanceof GamemasterApi) {
                    scene = new GamemasterScene(session)
                } else {
                    scene = new PlayerCharacterScene(session)
                }
                
                emit(new ReplaceScene(scene))            
            } catch (e) {
                console.log(`Got error while loading session: ${e}. Considering session closed.`)
                emit(new SessionClosed())
            } finally {
                this.inUpdate = false
            }
        }
    
        this.updateIntervalHandle = setInterval(requestAndEmitUpdate, 1000)
        requestAndEmitUpdate()
    }   

    async update({ model, message, emit }: wecco.UpdaterContext<Model, Message>): Promise<Model | typeof wecco.NoModelChange> {
        switch (message.command) {
            case "replace-scene":
                return new Model(model.versionInfo, message.scene, ...message.notifications)

            case "post-notification":
                return new Model(model.versionInfo, model.scene, ...message.notifications)

            case "session-closed":
                this.clearUpdate()
                this.api = null
                history.pushState(null, "", "/")
                return new Model(model.versionInfo, new HomeScene(), new Notification(m("sessionClosed.message")))

            case "new-session":
                this.clearUpdate()
                this.api = await GamemasterApi.createSession(message.title)
                history.pushState(null, "", `/session/${this.api.sessionId}`)
                this.scheduleUpdates(emit)
                break

            case "join-as-player":
                this.clearUpdate()
                this.api = await PlayerCharacterApi.joinGame(message.id, message.name)
                history.pushState(null, "", `/session/${this.api.sessionId}`)
                this.scheduleUpdates(emit)
                break

            case "rejoin":
                this.clearUpdate()
                try {
                    this.api = await rejoinSession(message.sessionId)
                    this.scheduleUpdates(emit)
                } catch {
                    document.location.pathname = ""
                    return new Model(model.versionInfo, new HomeScene())
                }
                break


            case "update-fate-points":
                if (this.api instanceof GamemasterApi) {
                    this.api.updateFatePoints(message.playerId, message.fatePoints)
                }
                break

            case "spend-fate-point":
                if (this.api instanceof PlayerCharacterApi) {
                    this.api.spendFatePoint()
                }
                break

            case "add-aspect":
                if (this.api instanceof GamemasterApi) {
                    this.api.addAspect(message.name, message.targetPlayerId)
                }
                break

            case "remove-aspect":
                if (this.api instanceof GamemasterApi) {
                    this.api.removeAspect(message.id)
                }
                break
        }

        return model.notifications.length === 0 ? wecco.NoModelChange : model.pruneNotifications()
    }
}

async function rejoinSession (sessionId: string): Promise<GamemasterApi | PlayerCharacterApi> {
    let apiClient = await createApiClient()

    let [authInfo, session] = await Promise.all([
        apiClient.authorization.getAuthenticationInfo(),
        apiClient.session.getSession({ id: sessionId }),
    ])

    if (authInfo.userId == session.ownerId) {
        return new GamemasterApi(apiClient, sessionId)
    } else {
        console.log(session)
        let characterId = session.characters.find(c => c.ownerId == authInfo.userId)?.id
        console.log(characterId)
        return new PlayerCharacterApi(apiClient, sessionId, characterId!)
    }
}

async function withMaxRetries<R>(fn: () => Promise<R>, maxTries=5, backOff=200): Promise<R> {
    let tries = 1
    let sleepTime = 0

    while (true) {
        try {
            return await fn()
        } catch (e) {
            if (tries == maxTries) {
                throw e
            }
            console.error(`Got error while executing call: ${e}. Retrying ${maxTries - tries} times...`)
            tries++
            sleepTime += backOff
            await sleep(sleepTime)
        }
    }
}

function sleep (ms: number): Promise<void> {
    return new Promise(resolve => {
        setTimeout(resolve, ms)
    })
}