import * as wecco from "@weccoframework/core"
import { Controller, JoinCharacter, RejoinSession, ReplaceScene } from "./control"
import { Aspect, Gamemaster, Home, Model, Player, PlayerCharacter, Session, VersionInfo } from "./models"
import { load, m } from "./utils/i18n"
import { root } from "./views"

import "roboto-fontface/css/roboto/roboto-fontface.css"
import "material-icons/iconfont/material-icons.css"
import "./index.css"

// This eventlistener boostraps the wecco application
document.addEventListener("DOMContentLoaded", async () => {
    // Wait for i18n files to be loaded
    await load()
    const versionInfo = await loadVersionInfo()

    const controller = new Controller()

    const context = wecco.app(() => new Model(versionInfo, new Home()), controller.update.bind(controller), root, "#app")

    if (document.location.pathname.startsWith("/join/")) {
        const sessionId = document.location.pathname.replace("/join/", "")
        const name = prompt(m("home.joinSession.promptName"))
        if (name !== null) {
            document.location.hash = ""
            context.emit(new JoinCharacter(sessionId, name.trim()))
        }

    } else if (document.location.pathname.startsWith("/session/")) {
        const sessionId = document.location.pathname.replace("/session/", "")
        context.emit(new RejoinSession(sessionId))

        // The following branches serve as an easy way to "view" a scene in dev mode. 
        // TODO: Replace this with https://storybook.js.org/ or something similar.
    } else if (document.location.hash.startsWith("#dev/gamemaster")) {
        context.emit(new ReplaceScene(new Gamemaster(new Session(
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
        context.emit(new ReplaceScene(new PlayerCharacter(new Session(
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
