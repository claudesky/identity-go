import assert from "node:assert"

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, 'id')

  assert(id)

  return createIdentityApiFetch(event)(`/clients/${id}`)
})
