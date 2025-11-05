export default defineNuxtPlugin((nuxtApp) => {
  const requestEvent = useRequestEvent()

  const identityApi = $fetch.create({
    async onRequest({ options }) {
      if (!requestEvent) return

      const cookies = requestEvent.headers.get('cookie')
      if (cookies) {
        options.headers = options.headers || {}
        options.headers = {
          ...options.headers,
          cookie: cookies,
        } as any
      }
    },
    async onResponse({ response }) {
      if (!requestEvent) return

      const setCookieHeaders = response.headers.getSetCookie?.() || []
      for (const cookie of setCookieHeaders) {
        requestEvent.node.res.appendHeader('set-cookie', cookie)
      }
    },
  })

  return {
    provide: {
      identityApi,
    },
  }
})
