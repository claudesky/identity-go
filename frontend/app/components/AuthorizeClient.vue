<script setup lang="ts">
const props = defineProps<{
  clientId: string
  redirectUri: string
  responseType: string
  scope?: string
  state?: string
  codeChallenge?: string
  codeChallengeMethod?: string
}>()

const { getClient } = useClients()
const { user } = useAuth()

// Component state
const error = ref('')
const approving = ref(false)

const {
  data: clientData,
  error: getClientError,
  pending,
} = await getClient(props.clientId)

const client = computed(() => clientData.value?.data)

const handleApprove = async () => {
  approving.value = true
  error.value = ''

  try {
    const body: any = {
      response_type: props.responseType,
      client_id: props.clientId,
      redirect_uri: props.redirectUri,
    }

    if (props.scope) body.scope = props.scope
    if (props.state) body.state = props.state
    if (props.codeChallenge) body.code_challenge = props.codeChallenge
    if (props.codeChallengeMethod)
      body.code_challenge_method = props.codeChallengeMethod

    await $fetch('/api/oauth/authorize', {
      method: 'POST',
      body,
    })
  } catch (err: any) {
    error.value = err.data?.message || 'Authorization failed'
    approving.value = false
  }
}

const handleDeny = () => {
  const url = new URL(props.redirectUri)
  url.searchParams.set('error', 'access_denied')
  url.searchParams.set('error_description', 'User denied authorization')
  if (props.state) url.searchParams.set('state', props.state)
  window.location.href = url.toString()
}

const scopeList = computed(() =>
  props.scope ? props.scope.split(' ').filter(Boolean) : []
)
</script>

<template>
  <div class="w-full max-w-md">
    <div v-if="pending" class="bg-white p-8 rounded-lg shadow text-center">
      <p>Loading application details...</p>
    </div>

    <div v-else-if="getClientError" class="bg-white p-8 rounded-lg shadow">
      <h1 class="text-2xl font-bold text-red-600 mb-4">Authorization Error</h1>
      <div class="bg-red-50 text-red-800 p-3 rounded">{{ getClientError }}</div>
    </div>

    <div v-else-if="error && !client" class="bg-white p-8 rounded-lg shadow">
      <h1 class="text-2xl font-bold text-red-600 mb-4">Authorization Error</h1>
      <div class="bg-red-50 text-red-800 p-3 rounded">{{ error }}</div>
    </div>

    <div v-else-if="client" class="bg-white p-8 rounded-lg shadow">
      <div class="text-center mb-6 pb-6 border-b">
        <h1 class="text-2xl font-bold text-gray-900 mb-2">
          Authorize Application
        </h1>
        <p class="text-gray-600">
          {{ client.name }} wants to access your account
        </p>
      </div>

      <div class="space-y-4 mb-6">
        <div>
          <p class="text-xs font-semibold text-gray-500 uppercase mb-1">
            Signed in as
          </p>
          <p class="text-gray-900">{{ user?.email }}</p>
        </div>

        <div>
          <p class="text-xs font-semibold text-gray-500 uppercase mb-1">
            Application
          </p>
          <p class="text-gray-900">{{ client.name }}</p>
        </div>

        <div v-if="scopeList.length > 0">
          <p class="text-xs font-semibold text-gray-500 uppercase mb-1">
            Requested permissions
          </p>
          <ul class="space-y-1">
            <li
              v-for="s in scopeList"
              :key="s"
              class="bg-gray-100 px-3 py-1 rounded text-sm font-mono"
            >
              {{ s }}
            </li>
          </ul>
        </div>

        <div>
          <p class="text-xs font-semibold text-gray-500 uppercase mb-1">
            Will redirect to
          </p>
          <p class="text-gray-600 text-sm font-mono break-all">
            {{ redirectUri }}
          </p>
        </div>
      </div>

      <div v-if="error" class="bg-red-50 text-red-800 p-3 rounded mb-4 text-sm">
        {{ error }}
      </div>

      <div class="flex gap-3 mb-6">
        <button
          @click="handleDeny"
          :disabled="approving"
          class="flex-1 px-4 py-2 border border-gray-300 rounded font-semibold hover:bg-gray-100 disabled:opacity-50"
        >
          Deny
        </button>
        <button
          @click="handleApprove"
          :disabled="approving"
          class="flex-1 px-4 py-2 bg-blue-600 text-white rounded font-semibold hover:bg-blue-700 disabled:opacity-50"
        >
          {{ approving ? 'Authorizing...' : 'Authorize' }}
        </button>
      </div>

      <p
        class="text-xs text-gray-500 text-center pt-4 border-t leading-relaxed"
      >
        By clicking "Authorize", you allow this application to access your
        information according to the requested permissions.
      </p>
    </div>
  </div>
</template>
