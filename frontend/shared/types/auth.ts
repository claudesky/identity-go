export interface SignUpRequest {
  name: string
  email: string
  password: string
}

export interface SignInRequest {
  email: string
  password: string
}

export interface AuthData {
  sub: string
  name: string
  email: string
}

export interface TokenData {
  access_token: string
  refresh_token: string
}
