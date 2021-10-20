import * as wecco from "@weccoframework/core"
import { v4 } from "uuid"
import { Controller, JoinTable } from "./control"
import { Home } from "./models"
import { load, m } from "./utils/i18n"
import { root } from "./views"

import "../css/index.css"

document.addEventListener("DOMContentLoaded", async () => {
    await load()

    const userId = v4()

    const controller = new Controller()

    const context = wecco.app(() => new Home(userId), controller.update.bind(controller), root, "#app")

    if (document.location.hash.startsWith("#join/")) {
        const tableId = document.location.hash.replace("#join/", "")
        const name = prompt(m("joinTable.prompt"))
        if (name !== null) {
            document.location.hash = ""
            context.emit(new JoinTable(tableId, name.trim()))
        }
    }
})
