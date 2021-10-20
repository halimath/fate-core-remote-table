import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"
import { Message, SpendFatePoint } from "../../control"
import { PlayerCharacter } from "../../models"
import { appShell, container } from "../components"
import { result } from "../components/result"

export function player(model: PlayerCharacter, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return appShell(
        container(
            wecco.html`<div class="grid grid-cols-1 divide-y-2 divide-blue-800">
    ${fatePoints(model.fatePoints, context)}
    ${result(context, model.result)}
</div>`
        ),
        model.table.title
    )
}

const fatePointActions = [
    "invokeAspect",
    "powerStunt",
    "refuseCompel",
    "storyDetail",
]

function fatePoints(fatePoints: number, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    const buttonStyles = `bg-yellow-${fatePoints === 0 ? "300" : "500 hover:bg-yellow-700"} text-white font-bold py-2 px-4 rounded shadow-lg`

    return wecco.html`<div class="flex items-center justify-around">
        <div class="flex items-center justify-center flex-col">
            <span class="text-3xl font-bold text-yellow-600">${fatePoints}</span>
            <button class=${buttonStyles} ?disabled=${fatePoints === 0} @click=${() => context.emit(new SpendFatePoint())}>${m(`player.spendFatePoint`)}</button>
        </div>
        <ul class="text-gray-400">
            ${fatePointActions.map(a => wecco.html`<li>${m(`player.spendFatePoint.${a}`)}</li>`)}
        </ul>
    </div>`
}