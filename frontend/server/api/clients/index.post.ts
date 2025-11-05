export default defineEventHandler(async (event) => {
  const body = await readBody<CreateClientRequest>(event)

  return createIdentityApiFetch(event)<DataResponse<CreateClientResponse>>(
    '/clients',
    {
      method: 'POST',
      body,
    }
  )
})
