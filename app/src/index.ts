import * as wecco from "@weccoframework/core"
import "material-icons/iconfont/material-icons.css"
import "roboto-fontface/css/roboto/roboto-fontface.css"
import { Controller, Message, Rejoin, ReplaceScene } from "./control"
import "./index.css"
import { Aspect, GamemasterScene, HomeScene, Model, Player, PlayerCharacterScene, Session, VersionInfo } from "./models"
import { load } from "./utils/i18n"
import { root } from "./views"


// This eventlistener boostraps the wecco application
document.addEventListener("DOMContentLoaded", async () => {
    // Wait for i18n files to be loaded
    await load()
    const versionInfo = await loadVersionInfo()

    const controller = new Controller()

    const app = wecco.createApp<Model, Message>(() => new Model(versionInfo, new HomeScene()), controller.update.bind(controller), root).mount("#app")

    if (document.location.pathname.startsWith("/join/")) {
        const sessionId = document.location.pathname.replace("/join/", "")
        document.location.hash = ""
        app.emit(new ReplaceScene(new HomeScene(sessionId)))

    } else if (document.location.pathname.startsWith("/session/")) {        
        const sessionId = document.location.pathname.replace("/session/", "")
        app.emit(new Rejoin(sessionId))

        // The following branches serve as an easy way to "view" a scene in dev mode. 
        // TODO: Replace this with https://storybook.js.org/ or something similar.
    } else if (document.location.hash.startsWith("#dev/gamemaster")) {
        app.emit(new ReplaceScene(new GamemasterScene(new Session(
            "1",
            "Test Table",
            "1",
            [
                new Player("2", "Cynere", false, 2, [
                    new Aspect("5", "Cover"),
                ]),
                new Player("3", "Landon", false, 4, []),
                new Player("4", "Zird", false, 0, []),
            ],
            [
                new Aspect(
                    "1",
                    "Fog",
                ),
                new Aspect(
                    "2",
                    "Slippy Grounds",
                ),
            ]
        ))))
    } else if (document.location.hash.startsWith("#dev/player")) {
        app.emit(new ReplaceScene(new PlayerCharacterScene(new Session(
            "1",
            "Test Table",
            "1",
            [
                new Player("2", "Cynere", true, 2, [
                    new Aspect("5", "Cover"),
                ]),
                new Player("3", "Landon", false, 4, []),
                new Player("4", "Zird", false, 1, []),
            ],
            [
                new Aspect(
                    "1",
                    "Fog",
                ),
                new Aspect(
                    "2",
                    "Slippy Grounds",
                ),
            ]
        ))))
    }
})

async function loadVersionInfo(): Promise<VersionInfo> {
    const res = await fetch("/api/version-info")
    return await res.json()
}
