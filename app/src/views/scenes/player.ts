import * as wecco from "@weccoframework/core"
import { Message, SpendFatePoint } from "../../control"
import { Aspect, Player, PlayerCharacterScene, VersionInfo } from "../../models"
import { m } from "../../utils/i18n"
import { appShell, button, card, container } from "../widgets/ui"

export function player(versionInfo: VersionInfo, model: PlayerCharacterScene, emit: wecco.MessageEmitter<Message>): wecco.ElementUpdate {
    const title = `${model.session.self?.name} @ ${model.session.title}`
    document.title = title
    return appShell({
        body: container(content(model, emit)),
        title,
        versionInfo,
    })
}

function content(scene: PlayerCharacterScene, emit: wecco.MessageEmitter<Message>): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1">
        <fcrt-skillcheck></fcrt-skillcheck>
        ${card([
            wecco.html`<span class="text-yellow-700 text-xl font-bold">${scene.player.name}</span>`,
            fatePoints(scene.fatePoints, emit),
            aspects(scene.player.aspects)
        ])}

        ${aspects(scene.session.aspects)}
        ${players(scene)}
    </div>`
}

function players(scene: PlayerCharacterScene): wecco.ElementUpdate {
    return scene.session.players.filter(p => p !== scene.player).map(p => otherPlayer(p))
}

function otherPlayer(player: Player): wecco.ElementUpdate {
    return card(wecco.html`<div class="grid grid-cols-1">
        <span class="text-lg font-bold text-yellow-700">${player.name}</span>
        ${aspects(player.aspects)}
    </div>`)
}

// 
//         ${table.players.map(p => p.aspects.map(a => aspect(a, p)))}


function aspects(aspects: Array<Aspect>): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1" data-testid="aspects">
        ${aspects.map(a => aspect(a))}
    </div>`
}

function aspect(aspect: Aspect, player?: Player): wecco.ElementUpdate {
    return card(wecco.html`
        <span class="text-lg text-blue-800 flex-grow-1">${aspect.name}</span>
        ${player ? wecco.html`<span class="text-sm bg-blue-200 rounded p-1 ml-2">${player.name}</span>` : ""}
    `)
}

const fatePointActions = [
    "invokeAspect",
    "powerStunt",
    "refuseCompel",
    "storyDetail",
]

function fatePoints(fatePoints: number, emit: wecco.MessageEmitter<Message>): wecco.ElementUpdate {
    return wecco.html`<div class="flex items-center justify-around">
        <div class="flex items-center justify-center flex-col">
            <span class="text-3xl font-bold text-yellow-600" data-testid="fate-points">${fatePoints}</span>
            ${button({
        label: m(`player.spendFatePoint`),
        onClick: () => emit(new SpendFatePoint()),
        color: "yellow",
        disabled: fatePoints === 0,
        testId: "spend-fate-point",
    })}
        </div>
        <ul class="text-gray-400">
            ${fatePointActions.map(a => wecco.html`<li>${m(`player.spendFatePoint.${a}`)}</li>`)}
        </ul>
    </div>`
}