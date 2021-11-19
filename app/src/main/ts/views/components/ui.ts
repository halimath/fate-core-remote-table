import * as wecco from "@weccoframework/core"
import { VersionInfo } from "../../models"

export interface AppShellOptions {
    body: wecco.ElementUpdate
    title: string
    versionInfo: VersionInfo
    additionalAppBarContent?: wecco.ElementUpdate
}

export function appShell(opts: AppShellOptions): wecco.ElementUpdate {
    return wecco.html`
        <div class="flex flex-col min-h-screen">
            <header class="sticky top-0 z-30 w-full max-w-8xl mx-auto mb-2 bg-white flex-none flex bg-blue-900">
                <div
                    class="flex-auto h-16 flex items-center justify-between px-4 sm:px-6 lg:mx-20 lg:px-0 xl:mx-8 text-white font-bold text-lg">
                    <span>${opts.title}</span>
                    ${opts.additionalAppBarContent ? wecco.html`<span>${opts.additionalAppBarContent}</span>` : ""}
                </div>
            </header>
            <main class="flex-grow">
                ${opts.body}
            </main>
            <footer class="bg-blue-200 h-20 text-gray-600 text-xs flex items-center justify-around px-2">
                <div>
                    Fate Core Table v${opts.versionInfo.version} (${opts.versionInfo.commit}).
                    &copy; 2021 Alexander Metzner.
                    <a href="https://github.com/halimath/fate-core-remote-table">github.com/halimath/fate-core-remote-table</a><br><br>
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

export function card (content: wecco.ElementUpdate): wecco.ElementUpdate {
    return wecco.html`<div class="shadow-lg p-2 m-2 border rounded">${content}</div>`
}

export type ButtonCallback = () => void

export interface ButtonOpts {
    label?: wecco.ElementUpdate
    color?: string
    onClick?: ButtonCallback
    size?: "s" | "m" | "l"
    disabled?: boolean
}

export function button(opts: ButtonOpts): wecco.ElementUpdate {
    const options: ButtonOpts = {
        label: opts.label ?? "",
        color: opts.color ?? "blue",
        onClick: opts.onClick ?? (() => void(0)),
        size: opts.size ?? "m",
        disabled: !!opts.disabled
    }

    const padding = (options.size === "s") ? 1 : ((options.size === "m") ? 2 : 4)

    const color = options.disabled ? "gray" : options.color

    let style = `bg-${color}-500 text-white font-bold py-${padding} px-${padding * 2} rounded mr-1`
    if (!options.disabled) {
        style += " hover:bg-${color}-700 shadow-lg"
    }

    return wecco.html`<button @click=${options.onClick} ?disabled=${options.disabled}
    class=${style}>${options.label}</button>`
}
