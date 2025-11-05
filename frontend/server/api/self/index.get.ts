export default defineEventHandler(async (event) => {
  return createIdentityApiFetch(event)('/self')
})
