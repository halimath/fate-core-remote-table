import * as wecco from "@weccoframework/core"
import { AddAspect, Message, RemoveAspect, UpdatePlayerFatePoints } from "../../control"
import { Aspect, Gamemaster, Player, VersionInfo } from "../../models"
import { m } from "../../utils/i18n"
import { showNotification } from "../widgets/notification"
import { appShell, button, card, container } from "../widgets/ui"

export function gamemaster(versionInfo: VersionInfo, model: Gamemaster, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    const title = `${m("gamemaster.title.gm")} @ ${model.session.title}`
    document.title = title
    return appShell({
        body: container(content(model, context)),
        title,
        versionInfo,
        additionalAppBarContent: [
            button({
                label: wecco.html`<i class="material-icons">share</i>`,
                onClick: share.bind(undefined, model),
            }),
        ]
    })
}

function content(model: Gamemaster, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1 divide-y">
        <fcrt-skillcheck></fcrt-skillcheck>

        <div class="grid grid-cols-1 lg:grid-cols-2 place-content-start">
            <div class="flex flex-col">
                ${model.session.aspects.map(aspect.bind(undefined, context))}
                <div class="flex justify-center ml-2 mr-2 mt-2">
                    ${button({
        label: wecco.html`<i class="material-icons">add</i> ${m("gamemaster.addAspect")}`,
        onClick: () => {
            const name = prompt(m("gamemaster.addAspect.prompt"))
            if (name !== null) {
                context.emit(new AddAspect(name))
            }
        },
    })}
                </div>
            </div>
            <div class="flex flex-col">            
                ${model.session.players.map(player.bind(undefined, context))}
            </div>        
        </div>
    </div>`
}

function aspect(context: wecco.AppContext<Message>, aspect: Aspect): wecco.ElementUpdate {
    return card(wecco.html`
        <div class="flex justify-between">
            <span class="text-lg font-bold text-blue-800 flex-grow-1">* ${aspect.name}</span>
            <a href="#" @click=${() => context.emit(new RemoveAspect(aspect.id))}><i class="material-icons">close</i></a>
        </div>
    `)
}

function player(context: wecco.AppContext<Message>, player: Player): wecco.ElementUpdate {
    return card(wecco.html`
        <h3 class="text-lg font-bold text-yellow-700">${player.name}</h3>
        <div class="grid grid-cols-2">
            <div>
                ${player.aspects.map(aspect.bind(undefined, context))}
                <div class="flex justify-center">
                    ${button({
        label: m("gamemaster.addAspect"),
        onClick: () => {
            const name = prompt(m("gamemaster.addAspect.prompt"))
            if (name !== null) {
                context.emit(new AddAspect(name, player.id))
            }
        },
        size: "s",
    })}
                </div>
            </div>
            ${fatePoints(player.fatePoints, fp => context.emit(new UpdatePlayerFatePoints(player.id, fp)))}
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
    })}
        <span class="text-3xl font-bold text-yellow-600">${fatePoints}</span>
        ${button({
        label: "+",
        onClick: onChange.bind(undefined, 1),
        color: "yellow",
        size: "s",
    })}
    </div>`
}

function share(model: Gamemaster) {
    const url = `${document.location.protocol}//${document.location.host}/join/${model.session.id}`
    navigator.clipboard.writeText(url)
    showNotification(m("gamemaster.shareLink.notification"))
}
