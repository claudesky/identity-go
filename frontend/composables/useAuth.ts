export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
}

export const useAuth = () => {
  const user = useState<User | null>('user', () => null)
  const accessToken = useCookie('identity_access_token', {
    maxAge: 60 * 15, // 15 minutes
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    httpOnly: true
  })
  const refreshToken = useCookie('identity_refresh_token', {
    maxAge: 60 * 60 * 24 * 7, // 7 days
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    httpOnly: true
  })

  const setUser = (userData: User | null) => {
    user.value = userData
  }

  const validateToken = async (): Promise<boolean> => {
    if (!accessToken.value) {
      user.value = null
      return false
    }

    try {
      const response = await $fetch<{ user?: User }>('/api/auth/validate', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${accessToken.value}`
        }
      })
      // TODO: add a self endpoint somewhere
      if (response) {
        user.value = {
          id: '1234',
          email: 'test',
          firstName: 'test',
          lastName: 'test',
        }
        return true
      } else {
        user.value = null
        return false
      }
    } catch (error) {
      console.error('Token validation failed:', error)
      user.value = null
      return false
    }
  }

  const logout = async () => {
    try {
      await $fetch('/api/auth/logout', {
        method: 'POST'
      })
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      user.value = null
      accessToken.value = null
      refreshToken.value = null
    }
  }

  return {
    user: readonly(user),
    setUser,
    validateToken,
    logout,
    isAuthenticated: computed(() => !!user.value)
  }
}
