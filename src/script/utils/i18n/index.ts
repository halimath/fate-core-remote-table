import * as i18n from "@weccoframework/i18n"

let messageResolver: i18n.MessageResolver

export const load = () => 
    i18n.MessageResolver.create(
        new i18n.CascadingBundleLoader(
            new i18n.JsonBundleLoader(i18n.fetchJson("/messages/default.json")),
            new i18n.JsonBundleLoader(i18n.fetchJsonByLocale("/messages")))
        )
        .then(mr => messageResolver = i18n.registerStandardFormatters(mr))

export function m (key: string, ...args: Array<unknown>): string {
    return messageResolver.m(key, ...args)
}        

export function mpl (key: string, amount: number, ...args: Array<unknown>): string {
    return messageResolver.mpl(key, amount, ...args)
}        
