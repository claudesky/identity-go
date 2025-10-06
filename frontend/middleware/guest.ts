export default defineNuxtRouteMiddleware(async (to, from) => {
  const { validateToken, isAuthenticated } = useAuth()

  // Only validate on server-side or initial client navigation
  // Prevents double validation on hydration
  if (import.meta.server || !from) {
    await validateToken()
  }

  // If authenticated, redirect to dashboard
  if (isAuthenticated.value) {
    return navigateTo('/dashboard')
  }
})
