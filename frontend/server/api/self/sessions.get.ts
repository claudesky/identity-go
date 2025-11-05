export default defineEventHandler(async (event): Promise<any> => {
  return createIdentityApiFetch(event)('/self/sessions')
})
