export default defineEventHandler(async (event) => {
  // TODO: Call backend logout API when it's implemented
  deleteCookie(event, 'identity_access_token')
  deleteCookie(event, 'identity_refresh_token')
  return { success: true }
})
