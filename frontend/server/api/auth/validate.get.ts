export default defineEventHandler(async (event) => {
  let accessToken

  const authHeader = getHeader(event, 'authorization')
  if (authHeader && authHeader.startsWith('Bearer ')) {
    accessToken = authHeader.substring(7)
  }

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'No access token found',
    })
  }

  try {
    const response = await event.context.fetchWithAuth('/auth/validate', {
      method: 'GET',
    })

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Validation failed',
    })
  }
})
