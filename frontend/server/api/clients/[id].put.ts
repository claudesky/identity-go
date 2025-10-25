import type {
  Client,
  UpdateClientRequest,
} from '../../../composables/useClients'
import type { DataResponse } from '../../../lib/response'

export default defineEventHandler(async (event) => {
  const accessToken = getCookie(event, 'identity_access_token')

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'No access token found',
    })
  }

  const id = getRouterParam(event, 'id')

  if (!id) {
    throw createError({
      statusCode: 400,
      message: 'Client ID is required',
    })
  }

  const body = await readBody<UpdateClientRequest>(event)

  try {
    const response = await event.context.fetchWithAuth<DataResponse<Client>>(
      `/clients/${id}`,
      { method: 'PUT', body }
    )

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Failed to update client',
    })
  }
})
