import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"

export function appShell(content: wecco.ElementUpdate, title: string, additionalAppBarContent?: wecco.ElementUpdate): wecco.ElementUpdate {
    return wecco.html`
        <div class="flex flex-col min-h-screen">
            <header class="sticky top-0 z-30 w-full max-w-8xl mx-auto mb-2 bg-white flex-none flex bg-blue-900">
                <div class="flex-auto h-16 flex items-center justify-between px-4 sm:px-6 lg:mx-20 lg:px-0 xl:mx-8 text-white font-bold text-lg">
                    <span>${title}</span>
                    ${additionalAppBarContent ? wecco.html`<span>${additionalAppBarContent}</span>` : ""}
                </div>
            </header>
            <main class="flex-grow">
                ${content}
            </main>
            <footer class="bg-blue-200 h-20 text-gray-600 text-xs flex items-center justify-around px-2">
                <div>
                    Fate Core Table v0.1.0.
                    &copy; 2021 Alexander Metzner.
                    <a href="https://github.com/halimath/fate-table">github.com/halimath/fate-table</a><br><br>
                    The Fate Core font is © Evil Hat Productions, LLC and is used with permission. The Four Actions icons were 
                    designed by Jeremy Keller.
                </div>
            </footer>
        </div>
        `
}

export function container(content: wecco.ElementUpdate): wecco.ElementUpdate {
    return wecco.html`<div class="container lg:mx-auto lg:px-4">${content}</div>`
}

export type ButtonStyle = "primary" | "secondary"
export type ButtonCallback = () => void

export function button(content: wecco.ElementUpdate, onClick: ButtonCallback, style: ButtonStyle = "primary"): wecco.ElementUpdate {
    const col = color(style)
    return wecco.html`<button @click=${onClick} class="bg-${col}-500 hover:bg-${col}-700 text-white font-bold py-2 px-4 rounded shadow-lg mr-1">${content}</button>`
}

function color(style: ButtonStyle): string {
    switch (style) {
        case "primary":
            return "blue"
        case "secondary":
            return "yellow"
    }
}
