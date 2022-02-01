import * as wecco from "@weccoframework/core"
import { Notification } from "../../models"

const ShowNotificationTime = 5000

export interface NotificationOpts {
    content: wecco.ElementUpdate
    style?: "info" | "error"
}

function isOpts(input: wecco.ElementUpdate | NotificationOpts): input is NotificationOpts {
    return typeof((input as any).content) !== "undefined"
}

export function showNotifications(notifications: Array<Notification>) {
    notifications.forEach(n => showNotification({
        content: n.content,
        style: n.style
    }))
}

export function showNotification(input: wecco.ElementUpdate | NotificationOpts) {
    const opts = isOpts(input) ? input : {
        content: input,
        style: "info",
    }

    const outlet = document.body.appendChild(document.createElement("div"))

    const init = (evt: Event) => {
        const e = evt.target as HTMLElement

        e.addEventListener("transitionend", () => {
            setTimeout(() => {
                e.addEventListener("transitionend", () => {
                    document.body.removeChild(outlet)
                }, { once: true })
                setTimeout(() => e.style.opacity = "0")
            }, ShowNotificationTime)
        }, { once: true })
        setTimeout(() => e.style.opacity = "1")
    }

    let style = "z-40 absolute top-1 right-1 flex w-full max-w-sm mx-auto overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800 my-5 transition-opacity duration-200 opacity-0 border-l-8"
    if (typeof opts.style === "undefined" || opts.style === "info") {
        style += " border-blue-600 text-blue-700"
    } else {
        style += " border-red-600 text-red-700"
    }

    wecco.updateElement(outlet, wecco.html`
<div @update=${init} class=${style}> 
  <div class="px-4 py-2 -mx-3">
      <div class="mx-3">
          <p class="text-sm">${opts.content}</p>
      </div>
  </div>
</div>    
    `)
}