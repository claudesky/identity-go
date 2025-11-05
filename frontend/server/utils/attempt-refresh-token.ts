import type { TokenData } from "../../shared/types/auth"
import type { DataResponse } from "../../shared/types/response"
import { useRuntimeConfig } from "nuxt/app"

export async function attemptRefreshToken(refreshToken: string) {
  const config = useRuntimeConfig()
  const response = await $fetch<DataResponse<TokenData>>(
    `${config.public.apiBaseURL}/auth/refresh`,
    {
      method: 'POST',
      body: { refresh_token: refreshToken },
    }
  )
  return response.data
}
