import type { DataResponse } from '../lib/response'
import type { User } from './useAuth'

export interface Session {
  id: string
  sub: string
  last_issued: string
  created_at: string
  last_issued_at: string
  expires_at: string
  revoked: boolean
}

export const useSelf = () => {
  const getSelf = () => {
    return useFetch<DataResponse<User>>('/api/self', {
      method: 'GET',
    })
  }

  const getSessions = async () => {
    const response = await $fetch<DataResponse<Session[]>>('/api/self/sessions', {
      method: 'GET',
    })

    return response
  }

  return {
    getSessions,
    getSelf,
  }
}
