import { defineConfig } from "vite"

export default defineConfig({
    server: {
        proxy: {
            "/api": "http://localhost:8080",
            "/ws": {
                target: "http://localhost:8080",
                ws: true
            },
            "/edit/*": {
                forward: "http://localhost:3000/",
            },
            "/view/*": {
                forward: "http://localhost:3000/",
            },
        }
    }
})
