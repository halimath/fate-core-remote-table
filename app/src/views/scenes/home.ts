import * as wecco from "@weccoframework/core"
import { JoinCharacter, Message, NewSession } from "../../control"
import { Home, VersionInfo } from "../../models"
import { m } from "../../utils/i18n"
import { appShell, button, container } from "../widgets/ui"

export function home(versionInfo: VersionInfo, model: Home, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return appShell({
        body: container(
            wecco.html`<div class="grid grid-cols-1 divide-y">
                <fcrt-skillcheck></fcrt-skillcheck>
                <div class="grid grid-col-1 items-center justify-around pt-2 gap-y-2">
                    ${button({
                label: m("home.joinSession"),
                onClick: () => {
                    const idOrUrl = prompt(m("home.joinSession.promptId"))
                    if (idOrUrl === null || idOrUrl.trim().length === 0) {
                        return
                    }

                    let id: string
                    if (idOrUrl.indexOf("/") >= 0) {
                        id = idOrUrl.substring(idOrUrl.lastIndexOf("/") + 1).trim()
                    } else {
                        id = idOrUrl.trim()
                    }

                    const name = prompt(m("home.joinSession.promptName"))
                    if (name === null || name.trim().length === 0) {
                        return
                    }

                    context.emit(new JoinCharacter(id, name.trim()))
                }
            })}

                    ${button({
                label: m("home.createNewSession"),
                color: "yellow",
                onClick: () => {
                    const title = prompt(m("home.createNewSession.prompt"))

                    if (title === null) {
                        return
                    }

                    context.emit(new NewSession(title))
                }
            })}                    
            </div>
        </div>`
        ),
        title: m("title"),
        versionInfo,
    })
}