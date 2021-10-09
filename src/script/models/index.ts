
/**
 * Sentinel value describing that no result is present.
 */
export const NoResult = {}

/**
 * The result of rolling a single fate die.
 */
export type Roll = 1 | 0 | -1

/**
 * Rating implements a character's skill or similar rating which gets added to the rolls.
 */
export type Rating = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8

/**
 * Total implements the overall result.
 */
export type Total = -4 | -3 | -2 | -1 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12

export class Result {
    constructor (public readonly rolls: [Roll, Roll, Roll, Roll], public readonly rating: Rating) {}

    get total(): Total {
        return (this.rating + this.rolls.reduce((p, c) => p + c, 0)) as Total
    }
}

export type Model = Result | typeof NoResult

export function roll(rating: Rating): Result {
    return new Result([
        randomRoll(),
        randomRoll(),
        randomRoll(),
        randomRoll(),
    ], rating)
}

function randomRoll(): Roll {
    return (Math.floor(Math.random() * 3) - 1) as Roll
}