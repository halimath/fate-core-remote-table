import * as wecco from "@weccoframework/core"
import { m } from "../../utils/i18n"

export interface Modal {
    show(): Promise<void>
    hide(): Promise<void>
}

export type ModalActionKind = "close" | "ok" | "danger"

export interface ModalAction {
    label: wecco.ElementUpdate
    kind?: ModalActionKind
    action: (modal: Modal) => void
}

export function modalCloseAction(label?: string): ModalAction {
    return {
        label: label ?? m("close"),
        kind: "close",
        action: modal => modal.hide(),
    }
}

export interface ModalOptions {
    title?: wecco.ElementUpdate
    body: wecco.ElementUpdate
    actions?: Array<ModalAction>
}

export function modal(opts: ModalOptions): Modal {
    const backdrop = document.createElement("div")
    // backdrop.classList.add("modal-backdrop")
    backdrop.classList.add("fixed", "inset-0", "bg-gray-500", "z-40", "backdrop-filter", "backdrop-blur-sm", "transition", "duration-100", "bg-opacity-0")

    const modalElement = document.createElement("div")
    modalElement.classList.add("fixed", "z-50", "inset-0", "overflow-y-auto", "transition-opacity", "duration-200", "opacity-0")

    let keyboardListener: (evt: KeyboardEvent) => void

    const modal = {
        show() {
            return new Promise<void>(resolve => {
                document.body.appendChild(backdrop)
                document.body.appendChild(modalElement)
                wecco.updateElement(modalElement, wecco.html`
                <div class="flex items-end justify-center lg:min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0 mt-24">
                    <div class="relative inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
                        <div class="bg-gray-100 px-4 py-3 sm:px-6 sm:flex"><h3>${opts.title ?? ""}</h3></div>
                        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                            <div class="sm:flex sm:items-start">
                                <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">                            
                                    ${opts.body}
                                </div>
                            </div>
                        </div>
                        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse border-top border-gray-100">
                            ${(opts.actions ?? []).map(a => renderModalAction(a, modal))}
                        </div>
                    </div>
                </div>    
                `)
                document.addEventListener("keydown", keyboardListener)
                backdrop.classList.add("bg-opacity-75")
                modalElement.classList.add("opacity-100")
                modalElement.addEventListener("transitionend", resolve.bind(null, void null), {
                    once: true,
                })
            })
        },

        hide() {
            return new Promise<void>(resolve => {
                document.removeEventListener("keydown", keyboardListener)
                backdrop.classList.remove("bg-opacity-75")
                modalElement.classList.remove("opacity-100")
                modalElement.addEventListener("transitionend", () => {
                    document.body.removeChild(modalElement)
                    document.body.removeChild(backdrop)
                    resolve(void null)    
                }, {
                    once: true,
                })
            })
        },
    }

    keyboardListener = (evt: KeyboardEvent) => {
        if (evt.key === "Escape") {
            document.removeEventListener("keydown", keyboardListener)
            modal.hide()
            return
        }

        if (evt.key === "Enter") {
            if (opts.actions) {
                opts.actions.find(a => (a.kind ?? "close") !== "close")?.action(modal)
            }
        }
    }

    return modal
}

function renderModalAction(a: ModalAction, modal: Modal): wecco.ElementUpdate {
    let classes = "mt-3 w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"

    switch (a.kind ?? "close") {
        case "close":
            classes += "border-gray-300 bg-white text-gray-700 hover:bg-gray-50"
            break
        case "danger":
            classes += "border-transparent bg-red-600 text-white hover:bg-red-700"
            break
        case "ok":
            classes += "border-transparent bg-blue-500 text-white hover:bg-blue-700"
            break
    }

    return wecco.html`<button type="button" @click=${() => a.action(modal)} class=${classes}>${a.label}</button>`
}