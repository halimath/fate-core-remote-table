import * as wecco from "@weccoframework/core"
import { Message, SpendFatePoint } from "../../control"
import { Aspect, Player, PlayerCharacter, Session, VersionInfo } from "../../models"
import { m } from "../../utils/i18n"
import { appShell, button, card, container } from "../widgets/ui"

export function player(versionInfo: VersionInfo, model: PlayerCharacter, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    const title = `${model.session.self?.name} @ ${model.session.title}`
    document.title = title
    return appShell({
        body: container(content(model, context)),
        title,
        versionInfo,
    })
}

function content(player: PlayerCharacter, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1 divide-y">
        <fcrt-skillcheck></fcrt-skillcheck>
        ${fatePoints(player.fatePoints, context)}
        ${aspects(player.session)}
    </div>`
}

function aspects(table: Session): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1">
        ${table.aspects.map(a => aspect(a))}
        ${table.players.map(p => p.aspects.map(a => aspect(a, p)))}
    </div>`
}

function aspect(aspect: Aspect, player?: Player): wecco.ElementUpdate {
    return card(wecco.html`
        <span class="text-lg font-bold text-blue-800 flex-grow-1">* ${aspect.name}</span>
        ${player ? wecco.html`<span class="text-sm bg-blue-200 rounded p-1 ml-2">${player.name}</span>` : ""}
    `)
}

const fatePointActions = [
    "invokeAspect",
    "powerStunt",
    "refuseCompel",
    "storyDetail",
]

function fatePoints(fatePoints: number, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return wecco.html`<div class="flex items-center justify-around">
        <div class="flex items-center justify-center flex-col">
            <span class="text-3xl font-bold text-yellow-600">${fatePoints}</span>
            ${button({
        label: m(`player.spendFatePoint`),
        onClick: () => context.emit(new SpendFatePoint()),
        color: "yellow",
        disabled: fatePoints === 0,
    })}
        </div>
        <ul class="text-gray-400">
            ${fatePointActions.map(a => wecco.html`<li>${m(`player.spendFatePoint.${a}`)}</li>`)}
        </ul>
    </div>`
}