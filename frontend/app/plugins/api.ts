import { parse, serialize } from 'cookie-es'

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

      // Update requestEvent.headers with new cookie data, replacing duplicates
      if (setCookieHeaders.length > 0) {
        const existingCookies = parse(requestEvent.headers.get('cookie') || '')

        // Parse new cookies from set-cookie headers and merge, overwriting duplicates
        for (const setCookieHeader of setCookieHeaders) {
          const [nameValue] = setCookieHeader.split(';')
          if (nameValue === undefined) continue
          const [name, value] = nameValue.split('=')
          if (name && value !== undefined) {
            existingCookies[name.trim()] = value.trim()
          }
        }

        // Serialize back to cookie string
        const updatedCookies = Object.entries(existingCookies)
          .map(([name, value]) => serialize(name, value))
          .join('; ')

        requestEvent.headers.set('cookie', updatedCookies)
      }
    },
  })

  return {
    provide: {
      identityApi,
    },
  }
})
