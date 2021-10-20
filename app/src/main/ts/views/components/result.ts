import * as wecco from "@weccoframework/core"
import { Message, RollDice } from "../../control"
import { Result, Roll, Total } from "../../models"
import { m } from "../../utils/i18n"
import { button } from "."

const selectedStyle = "text-blue-700 font-bold"

export function result(context: wecco.AppContext<Message>, result?: Result): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1 mt-2 pt-2">
        <div class="flex items-center justify-around">
        ${[
            button("+0", () => context.emit(new RollDice(0))),
            button("+1", () => context.emit(new RollDice(1))),
            button("+2", () => context.emit(new RollDice(2))),
            button("+3", () => context.emit(new RollDice(3))),
            button("+4", () => context.emit(new RollDice(4))),
            button("+5", () => context.emit(new RollDice(5))),
        ]}
        </div>
        <div class="flex items-center justify-around mt-2 pt-2">
            ${result ? resultView(result) : noResult()}
        </div>
    </div>`
}

function resultView(result: Result): wecco.ElementUpdate {
    return wecco.html`
        <div>
            ${ladder(r => (r === result.total) || (result.total === "above" && r === 8) || (result.total === "below" && r === -2) ? selectedStyle : "")}
        </div>
        <div class="mt-4 flex items-center justify-center text-blue-700">
            <span class="text-lg font-bold mr-2">${result.rating} +</span>
            ${result.rolls.map(r => roll(r))}
            <span class="text-lg font-bold">= ${result.total}</span>
        </div>
    `
}

function roll(roll: Roll): wecco.ElementUpdate {
    return wecco.html`<span class="text-lg font-bold mr-1 rounded border border-gray-400 w-8 h-8 flex items-center justify-center shadow">${(roll == -1) ? "-" : ((roll == 1) ? "+" : " ")}</span>`
}

function noResult(): wecco.ElementUpdate {
    return ladder(() => "") 
}

function ladder(styler: (r: Total) => string): wecco.ElementUpdate {
    return wecco.html`<ul class="text-lg mb-2 text-gray-400">
        <li class=${styler(8)}>${m("result.legendary")}</li>
        <li class=${styler(7)}>${m("result.epic")}</li>
        <li class=${styler(6)}>${m("result.fantastic")}</li>
        <li class=${styler(5)}>${m("result.superb")}</li>
        <li class=${styler(4)}>${m("result.great")}</li>
        <li class=${styler(3)}>${m("result.good")}</li>
        <li class=${styler(2)}>${m("result.fair")}</li>
        <li class=${styler(1)}>${m("result.average")}</li>
        <li class=${styler(0)}>${m("result.mediocre")}</li>
        <li class=${styler(-1)}>${m("result.poor")}</li>
    </ul>`
}