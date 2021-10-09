import * as wecco from "@weccoframework/core"
import { Message, Roll } from "../control"
import { Model, Result } from "../models"
import { m } from "../utils/i18n"
import { appShell, button, container } from "./components"

const commonStyle = "text-lg mb-2"
const selectedStyle = `${commonStyle} text-blue-700 font-bold`
const unselectedStyle = `${commonStyle} text-gray-400`

export function root (model: Model, context: wecco.AppContext<Message>): wecco.ElementUpdate {
    return appShell(
            container(
                wecco.html`<div class="grid grid-cols-1 lg:grid-cols-3 lg:gap-2 lg:divide-x divide-gray-400">
                    <div class="flex items-center justify-around lg:col-span-2">
                        ${[
                            button("+0", () => context.emit(new Roll(0))),
                            button("+1", () => context.emit(new Roll(1))),
                            button("+2", () => context.emit(new Roll(2))),
                            button("+3", () => context.emit(new Roll(3))),
                            button("+4", () => context.emit(new Roll(4))),
                            button("+5", () => context.emit(new Roll(5))),
                        ]}
                    </div>
                    <div class="flex items-center justify-center lg:justify-start lg:pl-8">
                        ${model instanceof Result ? result(model) : noResult()}
                    </div>
                </div>`
            )
        )
}

function result(result: Result): wecco.ElementUpdate {
    return wecco.html`
        <div>
            <ul>
                <li class=${result.total >= 8 ? selectedStyle : unselectedStyle }>${m("result.legendary")}</li>
                <li class=${result.total == 7 ? selectedStyle : unselectedStyle }>${m("result.epic")}</li>
                <li class=${result.total == 6 ? selectedStyle : unselectedStyle }>${m("result.fantastic")}</li>
                <li class=${result.total == 5 ? selectedStyle : unselectedStyle }>${m("result.superb")}</li>
                <li class=${result.total == 4 ? selectedStyle : unselectedStyle }>${m("result.great")}</li>
                <li class=${result.total == 3 ? selectedStyle : unselectedStyle }>${m("result.good")}</li>
                <li class=${result.total == 2 ? selectedStyle : unselectedStyle }>${m("result.fair")}</li>
                <li class=${result.total == 1 ? selectedStyle : unselectedStyle }>${m("result.average")}</li>
                <li class=${result.total == 0 ? selectedStyle : unselectedStyle }>${m("result.mediocre")}</li>
                <li class=${result.total == -1 ? selectedStyle : unselectedStyle }>${m("result.poor")}</li>
                <li class=${result.total <= -2 ? selectedStyle : unselectedStyle }>${m("result.terrible")}</li>
            </ul>
            <div class="mt-4 flex items-center justify-center text-blue-700">
                <span class="text-lg font-bold mr-2">${result.rating} +</span>
                ${result.rolls.map(r => wecco.html`<span class="text-lg font-bold mr-1 rounded border border-gray-400 w-8 h-8 flex items-center justify-center">${(r == -1) ? "-" : ((r == 1) ? "+" : " ")}</span>`)}
                <span class="text-lg font-bold">= ${result.total}</span>
            </div>
        </div>            
    `
}

function noResult(): wecco.ElementUpdate {
    return wecco.html`<ul>
        <li class=${unselectedStyle}>${m("result.legendary")}</li>
        <li class=${unselectedStyle}>${m("result.epic")}</li>
        <li class=${unselectedStyle}>${m("result.fantastic")}</li>
        <li class=${unselectedStyle}>${m("result.superb")}</li>
        <li class=${unselectedStyle}>${m("result.great")}</li>
        <li class=${unselectedStyle}>${m("result.good")}</li>
        <li class=${unselectedStyle}>${m("result.fair")}</li>
        <li class=${unselectedStyle}>${m("result.average")}</li>
        <li class=${unselectedStyle}>${m("result.mediocre")}</li>
        <li class=${unselectedStyle}>${m("result.poor")}</li>
        <li class=${unselectedStyle}>${m("result.terrible")}</li>
    </ul>`
}

