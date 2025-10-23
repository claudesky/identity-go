import type { CreateClientRequest, CreateClientResponse } from '../../../composables/useClients'
import type { DataResponse } from '../../../lib/response'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()

  const accessToken = getCookie(event, 'identity_access_token')

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'No access token found',
    })
  }

  const body = await readBody<CreateClientRequest>(event)

  try {
    const response = await $fetch<DataResponse<CreateClientResponse>>(
      `${config.public.apiBaseURL}/clients`,
      {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        body,
      }
    )

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Failed to create client',
    })
  }
})
