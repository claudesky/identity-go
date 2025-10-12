import type { SignUpRequest, SignInRequest } from '../lib/request'
import { useSelf } from './useSelf'

export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
}

export const useAuth = () => {
  const user = useState<User | null>('user', () => null)
  const { accessToken, clearTokens } = useTokens()
  const { getSelf } = useSelf()

  const setUser = (userData: User | null) => {
    user.value = userData
  }

  const loadSelf = async () => {
    const response = await getSelf()

    if (response.data.value === null) return false

    user.value = response.data.value.data
    return true
  }

  const validateToken = async (): Promise<boolean> => {
    if (!accessToken.value) {
      user.value = null
      return false
    }

    return loadSelf()
  }

  const register = async (data: SignUpRequest) => {}

  const login = async (data: SignInRequest) => {
    await $fetch<void>('/api/auth/login', {
      method: 'POST',
      body: data,
    })

    return loadSelf()
  }

  const logout = async () => {
    try {
      await $fetch('/api/auth/logout', {
        method: 'POST',
      })
    } finally {
      user.value = null
      clearTokens()
    }
  }

  return {
    user: readonly(user),
    setUser,
    validateToken,
    register,
    login,
    logout,
    isAuthenticated: computed(() => !!user.value),
  }
}
