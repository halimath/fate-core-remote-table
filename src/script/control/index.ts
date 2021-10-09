import { Model, NoResult, Rating, roll } from "../models"

export class Roll {
    readonly command = "roll"

    constructor (public readonly rating: Rating) {}
}

export type Message = Roll

export function update (model: Model, message: Message): Model {
    switch (message.command) {
        case "roll": 
            return roll(message.rating)
    }

    return NoResult
}
