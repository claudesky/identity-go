export interface Session {
  id: string
  sub: string
  last_issued: string
  created_at: string
  last_issued_at: string
  expires_at: string
  revoked: boolean
}
