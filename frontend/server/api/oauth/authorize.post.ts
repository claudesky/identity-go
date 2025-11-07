export default defineEventHandler(async (event) => {
  const body = await readBody(event)

  return createIdentityApiFetch(event)('/auth/authorize', {
    method: 'POST',
    body,
  })
})
