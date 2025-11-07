<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

useHead({
  title: 'Authorize Application',
})

const route = useRoute()

// Extract and validate OAuth parameters from query string
const clientId = computed(() => route.query.client_id as string)
const redirectUri = computed(() => route.query.redirect_uri as string)

// Validate required parameters
const hasRequiredParams = computed(() => !!clientId.value && !!redirectUri.value)
</script>

<template>
  <div>
    <AppNav />
    <main
      class="flex justify-center items-center min-h-[calc(100vh-200px)] p-8"
    >
      <AuthorizeClient
        v-if="hasRequiredParams"
        :client-id="clientId"
        :redirect-uri="redirectUri"
        :response-type="(route.query.response_type as string) || 'code'"
        :scope="(route.query.scope as string)"
        :state="(route.query.state as string)"
        :code-challenge="(route.query.code_challenge as string)"
        :code-challenge-method="(route.query.code_challenge_method as string)"
      />
      <div v-else class="w-full max-w-md bg-white p-8 rounded-lg shadow">
        <h1 class="text-2xl font-bold text-red-600 mb-4">
          Authorization Error
        </h1>
        <div class="bg-red-50 text-red-800 p-3 rounded">
          Missing required parameters: client_id and redirect_uri are required
        </div>
      </div>
    </main>
  </div>
</template>
