export default defineNuxtRouteMiddleware(async (to, from) => {
  const { validateToken, isAuthenticated } = useAuth()

  // Only validate on server-side or initial client navigation
  // Prevents double validation on hydration
  if (import.meta.server || !from) {
    await validateToken()
  }

  // If not authenticated, redirect to sign-in
  if (!isAuthenticated.value) {
    return navigateTo('/sign-in')
  }
})
