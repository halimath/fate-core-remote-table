import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"
import { button } from "../widgets/ui"

/**
 * The result of rolling a single fate die.
 */
export type Roll = 1 | 0 | -1

/**
 * Rating implements a character's skill or similar rating which gets added to the rolls.
 */
export type Rating = 0 | 1 | 2 | 3 | 4 | 5

/**
 * Total implements the overall result. It contains the numbers of "the ladder" (from -2 up to +8)
 * and contains values to describe "below" or "above" the ladder.
 */
export type Total = "below" | -2 | -1 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | "above"

export class Result {
    static roll(rating: Rating): Result {
        return new Result([
            randomRoll(),
            randomRoll(),
            randomRoll(),
            randomRoll(),
        ], rating)
    }

    constructor(
        public readonly rolls: [Roll, Roll, Roll, Roll],
        public readonly rating: Rating,
    ) { }

    get total(): Total {
        const val = (this.rating + this.rolls.reduce((p, c) => p + (c as number), 0 as number))
        if (val <= -3) {
            return "below"
        }

        if (val >= 9) {
            return "above"
        }

        return val as Total
    }
}

function randomRoll(): Roll {
    return (Math.floor(Math.random() * 3) - 1) as Roll
}

export interface SkillCheckData {
    result?: Result
}

export const SkillCheck = wecco.define("fcrt-skillcheck", ({data, requestUpdate}: wecco.RenderContext<SkillCheckData>): wecco.ElementUpdate => {
    return wecco.html`<div class="grid grid-cols-1 mt-2 mb-2 pt-2 pb-2 lg:grid-cols-3">
        <div class="fate-icon text-4xl text-blue-700 flex items-center justify-around">OCAD</div>
        <div class="flex items-center justify-around">
            ${[
            button({ label: "+0", testId: "skill-check-0-btn", onClick: () => { data.result = Result.roll(0); requestUpdate() } }),
            button({ label: "+1", testId: "skill-check-1-btn", onClick: () => { data.result = Result.roll(1); requestUpdate() } }),
            button({ label: "+2", testId: "skill-check-2-btn", onClick: () => { data.result = Result.roll(2); requestUpdate() } }),
            button({ label: "+3", testId: "skill-check-3-btn", onClick: () => { data.result = Result.roll(3); requestUpdate() } }),
            button({ label: "+4", testId: "skill-check-4-btn", onClick: () => { data.result = Result.roll(4); requestUpdate() } }),
            button({ label: "+5", testId: "skill-check-5-btn", onClick: () => { data.result = Result.roll(5); requestUpdate() } }),
        ]}
        </div>

        ${data.result ? resultView(data.result) : ""}
    </div>`

})

function resultView(result: Result): wecco.ElementUpdate {
    const total = result.total === "below" ? -2 : (result.total === "above" ? 8 : result.total)

    return wecco.html`
        <div class="flex flex-row items-center justify-center text-blue-700" data-testid="skill-check-result">
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
