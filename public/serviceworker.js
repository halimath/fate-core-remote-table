const CacheVersion = 8;
const CacheName = `fate-diceroller-cache-v${CacheVersion}`;

self.addEventListener("install", e => {
    e.waitUntil(caches.open(CacheName)
        .then(cache => {
            return cache.addAll([
                "/index.html",
                "/manifest.json",
                "/img/icon.png",
                "/img/icon-small.png",
                `/serviceworker.${CacheVersion}.js`,
                `/fate-diceroller.${CacheVersion}.js`,
                `/css/fate-diceroller.${CacheVersion}.css`,
                "/messages/default.json",
                "/messages/de.json",
            ]);
        })
    );

    self.addEventListener("fetch", event => {
        event.respondWith(caches.match(event.request).then(response => response || fetch(event.request)));
    });
});