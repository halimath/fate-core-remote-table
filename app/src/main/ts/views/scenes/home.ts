import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"
import { Message, NewTable } from "../../control"
import { Model } from "../../models"
import { appShell, button, container } from "../components"

export function home(model: Model, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return appShell(
        container(
            wecco.html`<div class="grid grid-col-1 flex items-center justify-around">
                ${button(m("createNewTable"), () => {
                const title = prompt(m("createNewTable.prompt"))

                if (title === null) {
                    return
                }
                
                context.emit(new NewTable(title))
            })}
            </div>`
        ),
        m("title")
    )
}