export interface SignUpRequest {
  firstName: string
  lastName: string
  email: string
  password: string
}

export interface SignInRequest {
  email: string
  password: string
}

export interface AuthResponse {
  token?: string
  user?: {
    id: string
    email: string
    firstName: string
    lastName: string
  }
  message?: string
}

export const useAuthAPI = () => {
  const signUp = async (data: SignUpRequest): Promise<AuthResponse> => {
    const response = await $fetch<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: data,
    })

    return response
  }

  const signIn = async (data: SignInRequest): Promise<AuthResponse> => {
    const response = await $fetch<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: data,
    })

    return response
  }

  const signOut = async (): Promise<void> => {
    await $fetch('/api/auth/logout', {
      method: 'POST',
    })
  }

  return {
    signUp,
    signIn,
    signOut
  }
}
