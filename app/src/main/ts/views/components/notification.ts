import * as wecco from "@weccoframework/core"

const OutletElementId = "notification-outlet"
const ShowNotificationTime = 5000

export interface NotificationOpts {
    content: wecco.ElementUpdate
    title?: string
    // type?: "info" | "warning" | "error"
}

function isOpts(input: wecco.ElementUpdate | NotificationOpts): input is NotificationOpts {
    return typeof((input as any).content) !== "undefined"
}

export function showNotification(input: wecco.ElementUpdate | NotificationOpts) {
    const opts = isOpts(input) ? input : {
        content: input
    }

    let outlet = document.getElementById(OutletElementId)

    if (outlet === null) {
        outlet = document.body.appendChild(document.createElement("div"))
        outlet.id = OutletElementId
    }

    wecco.removeAllChildren(outlet)

    const init = (evt: Event) => {
        const e = evt.target as HTMLElement
        e.classList.remove("opacity-0")
        setTimeout(() => {
            e.classList.add("opacity-0")
        }, ShowNotificationTime)
    }

    wecco.updateElement(outlet, wecco.html`
<div @update=${init} class="z-40 absolute top-1 right-1 flex w-full max-w-sm mx-auto overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800 my-5 transition-opacity duration-200 opacity-0">
  <div class="flex items-center justify-center w-12 bg-blue-500">
      <svg class="w-6 h-6 text-white fill-current" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <path d="M20 3.33331C10.8 3.33331 3.33337 10.8 3.33337 20C3.33337 29.2 10.8 36.6666 20 36.6666C29.2 36.6666 36.6667 29.2 36.6667 20C36.6667 10.8 29.2 3.33331 20 3.33331ZM21.6667 28.3333H18.3334V25H21.6667V28.3333ZM21.6667 21.6666H18.3334V11.6666H21.6667V21.6666Z"/>
      </svg>
  </div>
  
  <div class="px-4 py-2 -mx-3">
      <div class="mx-3">
          ${opts.title ? wecco.html`<span class="font-semibold text-blue-500 dark:text-blue-400">${opts.title}</span>` : ""}
          <p class="text-sm text-gray-600 dark:text-gray-200">${opts.content}</p>
      </div>
  </div>
</div>    
    `)
}