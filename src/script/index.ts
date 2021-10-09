import * as wecco from "@weccoframework/core"

import { load } from "./utils/i18n"
import { root } from "./views"

import "../styles/index.css"
import { update } from "./control"
import { NoResult } from "./models"

document.addEventListener("DOMContentLoaded", async () => {
    await load()

    const context = wecco.app(() => NoResult, update, root, "#app")
})
