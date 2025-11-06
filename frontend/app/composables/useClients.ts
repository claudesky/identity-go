export const useClients = () => {
  const { $identityApi } = useNuxtApp()

  const getClients = () => {
    return useAsyncData<DataResponse<Client[]>>(() =>
      $identityApi('/api/clients')
    )
  }

  const getClient = (id: string) => {
    return useAsyncData<DataResponse<Client>>(() =>
      $identityApi(`/api/clients/${id}`)
    )
  }

  const createClient = (data: CreateClientRequest) => {
    return useFetch<DataResponse<CreateClientResponse>>('/api/clients', {
      method: 'POST',
      body: data,
    })
  }

  const updateClient = (id: string, data: UpdateClientRequest) => {
    return useFetch<DataResponse<Client>>(`/api/clients/${id}`, {
      method: 'PUT',
      body: data,
    })
  }

  const deleteClient = (id: string) => {
    return useFetch<DataResponse<void>>(`/api/clients/${id}`, {
      method: 'DELETE',
    })
  }

  return {
    getClients,
    getClient,
    createClient,
    updateClient,
    deleteClient,
  }
}
