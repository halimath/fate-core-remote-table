import * as wecco from "@weccoframework/core"
import { Message, RollDice } from "../../control"
import { Rating, Result } from "../../models"
import { m } from "../../utils/i18n"
import { button } from "./ui"

export function result(context: wecco.AppContext<Message>, result?: Result): wecco.ElementUpdate {
    return wecco.html`<div class="grid grid-cols-1 mt-2 mb-2 pt-2 pb-2 lg:grid-cols-3">
        <div class="fate-icon text-4xl text-blue-700 flex items-center justify-around">OCAD</div>
        <div class="flex items-center justify-around">
            ${[
            button({ label: "+0", onClick: () => context.emit(new RollDice(0)), }),
            button({ label: "+1", onClick: () => context.emit(new RollDice(1)), }),
            button({ label: "+2", onClick: () => context.emit(new RollDice(2)), }),
            button({ label: "+3", onClick: () => context.emit(new RollDice(3)), }),
            button({ label: "+4", onClick: () => context.emit(new RollDice(4)), }),
            button({ label: "+5", onClick: () => context.emit(new RollDice(5)), }),
        ]}
        </div>

        ${result ? resultView(result) : ""}
    </div>`
}

function resultView(result: Result): wecco.ElementUpdate {
    const total = result.total === "below" ? -2 : (result.total === "above" ? 8 : result.total)

    return wecco.html`
        <div class="flex flex-row items-center justify-center text-blue-700">
            <span class="fate-icon text-xl lg:text-2xl">${result.rolls.map(r => (r == -1) ? "-" : ((r == 1) ? "+" : "0"))}</span>
            <span class="text-lg mr-2">${rating(result.rating)}</span>
            <span class="text-lg mr-2">=</span>
            <span class="text-lg lg:text-2xl font-bold mr-2">${m(`result.${total}`)}</span>
        </div>
    `
}

function rating(rating: Rating): string {
    if (rating >= 0) {
        return `+ ${rating}`
    }

    return `- ${rating}`
}

// function r(roll: Roll): wecco.ElementUpdate {
//     return wecco.html`<span class="text-lg font-bold mr-1 rounded border border-gray-400 w-8 h-8 flex items-center justify-center shadow">${(roll == -1) ? "-" : ((roll == 1) ? "+" : " ")}</span>`
// }

// function ladder(styler: (r: Total) => string): wecco.ElementUpdate {
//     return wecco.html`<ul class="text-lg mb-2 text-gray-400">
//         <li class=${styler(8)}>${m("result.legendary")}</li>
//         <li class=${styler(7)}>${m("result.epic")}</li>
//         <li class=${styler(6)}>${m("result.fantastic")}</li>
//         <li class=${styler(5)}>${m("result.superb")}</li>
//         <li class=${styler(4)}>${m("result.great")}</li>
//         <li class=${styler(3)}>${m("result.good")}</li>
//         <li class=${styler(2)}>${m("result.fair")}</li>
//         <li class=${styler(1)}>${m("result.average")}</li>
//         <li class=${styler(0)}>${m("result.mediocre")}</li>
//         <li class=${styler(-1)}>${m("result.poor")}</li>
//     </ul>`
// }