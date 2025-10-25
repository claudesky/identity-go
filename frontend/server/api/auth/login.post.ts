import type { DataResponse, TokenData } from "../../../lib/response"
import { setTokensFromData } from "../../utils/refreshToken"

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const body = await readBody(event)

  try {
    const response = await $fetch<DataResponse<TokenData>>(
      `${config.public.apiBaseURL}/auth/login`,
      {
        method: 'POST',
        body: body,
      }
    )

    setTokensFromData(event, response.data)

    // Return success without exposing tokens
    return
  } catch (error: any) {
    throw createError({
      statusCode: error.response?.status || 500,
      message: error.data?.message || 'Login failed',
    })
  }
})
