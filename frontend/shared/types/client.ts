export interface Client {
  id: string
  client_id: string
  name: string
  redirect_uris: string[]
  created_at: string
  updated_at: string
}

export interface CreateClientRequest {
  name: string
  redirect_uris: string[]
}

export interface UpdateClientRequest {
  name: string
  redirect_uris: string[]
}

export interface CreateClientResponse {
  client: Client
  client_secret: string
}
