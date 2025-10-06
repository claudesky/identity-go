export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()

  let accessToken;

  const authHeader = getHeader(event, 'authorization')
  if (authHeader && authHeader.startsWith('Bearer ')) {
    accessToken = authHeader.substring(7)
  }

  if (!accessToken) {
    throw createError({
      statusCode: 401,
      message: 'No access token found'
    })
  }

  try {
    const response = await $fetch(`${config.public.apiBaseURL}/auth/validate`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      }
    })

    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Validation failed'
    })
  }
})
