export const useSelf = () => {
  const getSelf = () => {
    return useFetch<DataResponse<User>>('/api/self', {
      method: 'GET',
    })
  }

  const getSessions = () => {
    return useFetch<DataResponse<Session[]>>('/api/self/sessions', {
      method: 'GET',
    })
  }

  return {
    getSessions,
    getSelf,
  }
}
