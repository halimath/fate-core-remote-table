
export type NotificationStyle = "info" | "error"
export class Notification {
    constructor(public readonly content: string, public readonly style: NotificationStyle = "info") { }
}

// --

export class Aspect {
    constructor(
        public readonly id: string,
        public readonly name: string,
    ) { }
}

// --

export class Player {
    constructor(
        public readonly id: string,
        public readonly name: string,
        public readonly isSelf: boolean,
        public readonly fatePoints: number,
        public readonly aspects: Array<Aspect>,
    ) { }
}

export class Session {
    constructor(
        public readonly id: string,
        public readonly title: string,
        public readonly gamemasterId: string,
        public readonly players: Array<Player>,
        public readonly aspects: Array<Aspect>,
    ) { }

    get self(): Player | undefined {
        return this.players.find(p => p.isSelf)
    }
}

export class PlayerCharacterScene {
    constructor(
        public readonly session: Session,
    ) { }

    get fatePoints(): number {
        return this.session.self?.fatePoints ?? 0
    }
}

export class GamemasterScene {
    constructor(
        public readonly session: Session,
    ) { }
}

export class HomeScene {
    constructor(public readonly joinSessionId?: string) { }
}

export type Scene = HomeScene | PlayerCharacterScene | GamemasterScene

export interface VersionInfo {
    version: string
    commit: string
}

export class Model {
    public readonly notifications: Array<Notification>

    constructor(
        public readonly versionInfo: VersionInfo,
        public readonly scene: Scene,
        ...notifications: Array<Notification>
    ) {
        this.notifications = notifications
    }

    pruneNotifications(): Model {
        return new Model(this.versionInfo, this.scene)
    }
}

