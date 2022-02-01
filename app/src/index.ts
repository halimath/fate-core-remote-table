import * as wecco from "@weccoframework/core"
import "./index.css"
import { Controller, JoinTable, ReplaceScene } from "./control"
import { Aspect, Gamemaster, Home, Model, Player, PlayerCharacter, Table, VersionInfo } from "./models"
import { load, m } from "./utils/i18n"
import { root } from "./views"

// This eventlistener boostraps the wecco application
document.addEventListener("DOMContentLoaded", async () => {
    // Wait for i18n files to be loaded
    await load()
    const versionInfo = await loadVersionInfo()

    const controller = new Controller()

    const context = wecco.app(() => new Model(versionInfo, new Home()), controller.update.bind(controller), root, "#app")

    if (document.location.pathname.startsWith("/join/")) {
        const tableId = document.location.pathname.replace("/join/", "")
        const name = prompt(m("home.joinTable.promptName"))
        if (name !== null) {
            document.location.hash = ""
            context.emit(new JoinTable(tableId, name.trim()))
        }

        // The following branches serve as an easy way to "view" a scene in dev mode. 
        // TODO: Replace this with https://storybook.js.org/ or something similar.
    } else if (document.location.hash.startsWith("#dev/gamemaster")) {
        context.emit(new ReplaceScene(new Gamemaster(new Table(
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
        context.emit(new ReplaceScene(new PlayerCharacter(new Table(
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
