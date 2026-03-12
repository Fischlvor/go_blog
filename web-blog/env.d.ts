/// <reference types="vite/client" />

export interface ImportMetaEnv {
    VITE_SERVER_URL: string
    VITE_BASE_API: string
    VITE_ADMIN_API: string
}

declare module 'vue-router' {
    interface RouteMeta {
        requiresAuth?: boolean
        requiresAdmin?: boolean
        title?: string
    }
}