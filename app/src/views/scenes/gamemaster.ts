import * as wecco from "@weccoframework/core"
import { AddAspect, Message, RemoveAspect, UpdatePlayerFatePoints } from "../../control"
import { Aspect, GamemasterScene, Player, VersionInfo } from "../../models"
import { m } from "../../utils/i18n"
import { modal, modalCloseAction } from "../widgets/modal"
import { showNotification } from "../widgets/notification"
import { appShell, button, card, container } from "../widgets/ui"

export function gamemaster(versionInfo: VersionInfo, model: GamemasterScene, emit: wecco.MessageEmitter<Message>): wecco.ElementUpdate {
    const title = `${m("gamemaster.title.gm")} @ ${model.session.title}`
    document.title = title
    return appShell({
        body: container(content(model, emit)),
        title,
        versionInfo,
        additionalAppBarContent: [
            button({
                label: wecco.html`<i class="material-icons">share</i>`,
                onClick: share.bind(undefined, model),
                testId: "join-session-link"
            }),
        ]
    })
}

function content(model: GamemasterScene, emit: wecco.MessageEmitter<Message>): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1 divide-y">
        <fcrt-skillcheck></fcrt-skillcheck>

        <div class="grid grid-cols-1 lg:grid-cols-2 place-content-start">
            <div class="flex flex-col" data-testid="aspects">
                ${model.session.aspects.map(aspect.bind(undefined, emit))}
                <div class="flex justify-center ml-2 mr-2 mt-2">
                    ${button({
        label: wecco.html`<i class="material-icons">add</i> ${m("gamemaster.addAspect")}`,
        onClick: addAspect.bind(null, emit, undefined),
        testId: "add-aspect",
    })}
                </div>
            </div>
            <div class="flex flex-col" data-testid="players">            
                ${model.session.players.map(player.bind(undefined, emit))}
            </div>        
        </div>
    </div>`
}

function aspect(emit: wecco.MessageEmitter<Message>, aspect: Aspect): wecco.ElementUpdate {
    return card(wecco.html`
        <div class="flex justify-between">
            <span class="text-lg text-blue-800 dark:text-blue-400 flex-grow-1">${aspect.name}</span>
            <a href="#" @click=${() => emit(new RemoveAspect(aspect.id))}><i class="material-icons text-gray-600">close</i></a>
        </div>
    `)
}

function player(emit: wecco.MessageEmitter<Message>, player: Player): wecco.ElementUpdate {
    return card(wecco.html`
        <h3 class="text-lg font-bold text-yellow-700" data-testid="name">${player.name}</h3>
        <div class="grid grid-cols-2">
            <div data-testid="aspects">
                ${player.aspects.map(aspect.bind(undefined, emit))}
                <div class="flex justify-center">
                    ${button({
        label: m("gamemaster.addAspect"),
        onClick: addAspect.bind(null, emit, player.id),
        size: "s",
        testId: "player-add-aspect",
    })}
                </div>
            </div>
            ${fatePoints(player.fatePoints, fp => emit(new UpdatePlayerFatePoints(player.id, fp)))}
        </div>
    `)
}

function fatePoints(fatePoints: number, onChange: (value: number) => void): wecco.ElementUpdate {
    return wecco.html`<div class="flex items-center justify-around">
        ${button({
        label: "-",
        onClick: onChange.bind(undefined, -1),
        color: "yellow",
        size: "s",
        disabled: fatePoints === 0,
        testId: "dec-fate-points",
    })}
        <span class="text-3xl font-bold text-yellow-600" data-testid="fate-points">${fatePoints}</span>
        ${button({
        label: "+",
        onClick: onChange.bind(undefined, 1),
        color: "yellow",
        size: "s",
        testId: "inc-fate-points",
    })}
    </div>`
}

function share(model: GamemasterScene) {
    const url = `${document.location.protocol}//${document.location.host}/join/${model.session.id}`
    navigator.clipboard.writeText(url)
    showNotification(m("gamemaster.shareLink.notification"))
}

function addAspect(emit: wecco.MessageEmitter<Message>, characterId?: string) {
    let nameInput: HTMLInputElement

    const bindNameInput = (e: Event) => {
        nameInput = e.target as HTMLInputElement
        nameInput.focus()
    }

    modal({
        title: m("gamemaster.addAspect"),
        body: wecco.html`
            <label for="aspect-name">${m("gamemaster.addAspect.prompt")}</label>
            <input type="text" id="aspect-name" data-testid="aspect-name" @update=${bindNameInput}>
        `,
        actions: [
            {
                label: m("ok"),
                kind: "ok",
                action: m => {
                    const name = nameInput.value.trim()
                    if (name === "") {
                        return
                    }
               
                    emit(new AddAspect(name, characterId))
                    m.hide()
                },
            },
            modalCloseAction(),
        ],
    }).show()    
}
