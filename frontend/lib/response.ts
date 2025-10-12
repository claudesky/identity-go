export interface AuthData {
  sub: string
  name: string
  email: string
}

export interface TokenData {
  access_token: string
  refresh_token: string
}

export interface DataResponse<T> {
  message: string
  status: number
  data: T
}
