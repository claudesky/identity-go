// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: {
    enabled: process.env.NODE_ENV !== 'production',

    timeline: {
      enabled: true
    }
  },
  debug: process.env.NODE_ENV !== 'production',

  devServer: {
    port: 3000
  },

  css: [
    '~/assets/css/main.css',
    '@fortawesome/fontawesome-svg-core/styles.css'
  ],

  app: {
    head: {
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
    }
  },

  runtimeConfig: {
    public: {
      apiBaseURL: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:9102',
      appName: process.env.NUXT_PUBLIC_APP_NAME || 'Identity - Go'
    }
  },

  modules: ['@nuxtjs/tailwindcss']
})