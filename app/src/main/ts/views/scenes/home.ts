import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"
import { JoinTable, Message, NewTable } from "../../control"
import { Model } from "../../models"
import { appShell, button, container } from "../components/ui"
import { result } from "../components/result"

export function home(model: Model, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return appShell(
        container(
            wecco.html`<div class="grid grid-cols-1 divide-y">
                ${result(context, model.result)}
                <div class="grid grid-col-1 flex items-center justify-around pt-2 gap-y-2">
                    ${button({
                        label: m("home.joinTable"),
                        onClick: () => {
                            const idOrUrl = prompt(m("home.joinTable.promptId"))                            
                            if (idOrUrl === null || idOrUrl.trim().length === 0) {
                                return
                            }

                            let id: string
                            if (idOrUrl.indexOf("/") >= 0) {
                                id = idOrUrl.substring(idOrUrl.lastIndexOf("/") + 1).trim()
                            } else {
                                id = idOrUrl.trim()
                            }

                            const name = prompt(m("home.joinTable.promptName"))
                            if (name === null || name.trim().length === 0) {
                                return
                            }

                            context.emit(new JoinTable(id, name.trim()))
                        }
                    })}

                    ${button({
                        label: m("home.createNewTable"),
                        color: "yellow",
                        onClick: () => {
                            const title = prompt(m("home.createNewTable.prompt"))
                            
                            if (title === null) {
                                return
                            }
                            
                            context.emit(new NewTable(title))
                        }
                    })}                    
            </div>
        </div>`
        ),
        m("title")
    )
}