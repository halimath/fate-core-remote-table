import * as wecco from "@weccoframework/core"
import { JoinCharacter, Message, NewSession } from "../../control"
import { Home, VersionInfo } from "../../models"
import { m } from "../../utils/i18n"
import { modal, modalCloseAction } from "../widgets/modal"
import { appShell, button, container } from "../widgets/ui"

export function home(versionInfo: VersionInfo, model: Home, emit: wecco.MessageEmitter<Message>): wecco.ElementUpdate {
    const onInit = (evt: Event) => {
        if (!model.joinSessionId) {
            return
        }

        joinSession(emit, model.joinSessionId)
    }

    return appShell({
        body: container(
            wecco.html`<div class="grid grid-cols-1 divide-y" @update=${onInit}>
                <fcrt-skillcheck></fcrt-skillcheck>
                <div class="grid grid-col-1 items-center justify-around pt-2 gap-y-2">
                    ${button({
                label: m("home.joinSession"),
                onClick: joinSession.bind(null, emit, undefined),
            })}

            ${button({
                label: m("home.createNewSession"),
                color: "yellow",
                onClick: startNewSession.bind(null, emit),
            })}                    
            </div>
        </div>`
        ),
        title: m("title"),
        versionInfo,
    })
}

function startNewSession(emit: wecco.MessageEmitter<Message>) {
    let titleInput: HTMLInputElement

    const bindTitleInput = (e: Event) => {
        titleInput = e.target as HTMLInputElement
        titleInput.focus()
    }

    modal({
        title: m("home.createNewSession"),
        body: wecco.html`
        <p>${m("home.createNewSession.prompt")}</p>
        <input type="text" @update=${bindTitleInput}>`,
        actions: [
            {
                label: m("ok"),
                kind: "ok",
                action: m => {
                    const title = titleInput.value.trim()
                    if (title === "") {
                        titleInput.setCustomValidity("Missing value")
                        return
                    }
                    
                    m.hide().then(() => emit(new NewSession(title)))
                },
            },
            modalCloseAction(),
        ],
    }).show()
}

function joinSession(emit: wecco.MessageEmitter<Message>, urlOrId?: string) {
    let idInput: HTMLInputElement
    let nameInput: HTMLInputElement

    const bindIdInput = (e: Event) => {
        idInput = e.target as HTMLInputElement
        if (!urlOrId) {
            idInput.focus()
        }
    }

    const bindNameInput = (e: Event) => {
        nameInput = e.target as HTMLInputElement
        if (urlOrId) {
            nameInput.focus()
        }
    }

    modal({
        title: m("home.joinSession"),
        body: wecco.html`
            <p>${m("home.joinSession.promptId")}</p>
            <input type="text" @update=${bindIdInput} value=${urlOrId ?? ""}>
            <p>${m("home.joinSession.promptName")}</p>
            <input type="text" @update=${bindNameInput}>
        `,
        actions: [
            {
                label: m("ok"),
                kind: "ok",
                action: m => {
                    const idOrUrl = idInput.value.trim()
                    idInput.setCustomValidity(idOrUrl === "" ? "Missing value" : "")

                    const name = nameInput.value.trim()
                    nameInput.setCustomValidity(name === "" ? "Missing value" : "")
                    
                    if (idOrUrl === "" || name === "") {
                        return
                    }

                    let id: string
                    if (idOrUrl.indexOf("/") >= 0) {
                        id = idOrUrl.substring(idOrUrl.lastIndexOf("/") + 1).trim()
                    } else {
                        id = idOrUrl.trim()
                    }
                
                    m.hide().then(() => emit(new JoinCharacter(id, name.trim())))                    
                },
            },
            modalCloseAction(),
        ],
    }).show()
}