export const useSelf = () => {
  const { $identityApi } = useNuxtApp()
  const getSelf = async () => {
    return useAsyncData<DataResponse<User>>(() => $identityApi('/api/self'))
  }

  const getSessions = () => {
    return useAsyncData<DataResponse<User>>(() =>
      $identityApi('/api/self/sessions')
    )
  }

  return {
    getSessions,
    getSelf,
  }
}
