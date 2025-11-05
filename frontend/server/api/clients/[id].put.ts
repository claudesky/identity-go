import assert from 'node:assert'

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, 'id')

  assert(id)

  const body = await readBody<UpdateClientRequest>(event)

  return createIdentityApiFetch(event)(`/clients/${id}`, {
    method: 'PUT',
    body,
  })
})
