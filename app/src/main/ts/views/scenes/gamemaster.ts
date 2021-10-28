import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"
import { Message, UpdatePlayerFatePoints } from "../../control"
import { Gamemaster, Player } from "../../models"
import { appShell, button, container } from "../components"
import { showNotification } from "../components/notification"
import { result } from "../components/result"

export function gamemaster(model: Gamemaster, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return appShell(
        container(content(model, context)),
        model.table.title,
        [
            button(wecco.html`<i class="material-icons">share</i>`, share.bind(undefined, model)),
            button(wecco.html`<i class="material-icons">power_off</i>`, () => void 0),
        ]
    )
}

function content(model: Gamemaster, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return wecco.html`
    <div class="grid grid-cols-1 lg:grid-cols-2">
        <div class="grid grid-cols-1 divide-y-2">            
            ${model.table.players.map(player.bind(undefined, context))}
        </div>
        ${result(context, model.result)}
    </div>`
}

function player (context: wecco.AppContext<Message>, player: Player): wecco.ElementUpdate {
    return wecco.html`<div class="shadow-lg p-2 m-2">
        <h3 class="text-lg font-bold text-blue-800">${player.name}</h3>
        ${fatePoints(player.fatePoints, fp => context.emit(new UpdatePlayerFatePoints(player.id, fp))) }
    </div>`
}

function fatePoints(fatePoints: number, onChange: (value: number) => void): wecco.ElementUpdate {
    return wecco.html`<div class="flex items-center justify-around">
    <button class="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded shadow-lg"
        ?disabled=${fatePoints === 0} @click=${onChange.bind(undefined, fatePoints - 1)}>-</button>
    <span class="text-3xl font-bold text-yellow-600">${fatePoints}</span>
    <button class="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded shadow-lg" @click=${onChange.bind(undefined, fatePoints + 1)}>+</button>
</div>`
}

function share(model: Gamemaster) {
    const url = `${document.location.protocol}//${document.location.host}#join/${model.table.id}`
    navigator.clipboard.writeText(url)
    showNotification(m("gamemaster.shareLink.notification"))
}
