<script setup lang="ts">
import { useClients, type UpdateClientRequest } from '~/composables/useClients'

definePageMeta({
  middleware: 'auth'
})

useHead({
  title: 'Edit OAuth Client'
})

const route = useRoute()
const router = useRouter()
const { getClient, updateClient } = useClients()

const clientId = route.params.id as string

const { data: client, error: loadError, pending: loading } = await useAsyncData(
  `client-${clientId}`,
  async () => {
    const result = await getClient(clientId)
    if (result.data.value === undefined) {
      throw result.error.value || new Error('Failed to load client')
    }
    return result.data.value.data
  }
)

const formData = reactive<UpdateClientRequest>({
  name: '',
  redirect_uris: ['']
})

// Populate form when client data loads
watch(client, (newClient) => {
  if (newClient) {
    formData.name = newClient.name
    formData.redirect_uris = [...newClient.redirect_uris]
  }
}, { immediate: true })

const submitting = ref(false)
const error = ref('')

const addRedirectUri = () => {
  formData.redirect_uris.push('')
}

const removeRedirectUri = (index: number) => {
  formData.redirect_uris.splice(index, 1)
}

const handleSubmit = async () => {
  error.value = ''

  // Validate
  if (!formData.name.trim()) {
    error.value = 'Client name is required'
    return
  }

  const validUris = formData.redirect_uris.filter(uri => uri.trim())
  if (validUris.length === 0) {
    error.value = 'At least one redirect URI is required'
    return
  }

  submitting.value = true
  try {
    const result = await updateClient(clientId, {
      name: formData.name.trim(),
      redirect_uris: validUris
    })

    if (result.error.value) {
      throw result.error.value
    }

    router.push('/clients')
  } catch (err: any) {
    console.error('Failed to update client:', err)
    error.value = err.message || 'Failed to update client'
  } finally {
    submitting.value = false
  }
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  alert('Copied to clipboard!')
}
</script>

<template>
  <div>
    <AppNav />
    <main class="container">
      <div class="header">
        <h1>Edit OAuth Client</h1>
        <button class="btn-secondary" @click="router.push('/clients')">Back to Clients</button>
      </div>

      <div v-if="loading" class="loading">Loading client...</div>
      <div v-else-if="loadError" class="error">{{ loadError.message || loadError }}</div>

      <form v-else @submit.prevent="handleSubmit" class="form-container">
        <div class="info-section">
          <h3>Client Information</h3>
          <div class="info-item">
            <label>Client ID</label>
            <div class="copy-field">
              <code>{{ client?.client_id }}</code>
              <button type="button" class="btn-copy" @click="copyToClipboard(client?.client_id || '')">
                Copy
              </button>
            </div>
          </div>
          <div class="info-item">
            <label>Created At</label>
            <span>{{ client?.created_at ? new Date(client.created_at).toLocaleString() : '' }}</span>
          </div>
          <div class="info-item">
            <label>Last Updated</label>
            <span>{{ client?.updated_at ? new Date(client.updated_at).toLocaleString() : '' }}</span>
          </div>
        </div>

        <div class="form-group">
          <label for="name">Client Name</label>
          <input
            id="name"
            v-model="formData.name"
            type="text"
            placeholder="My Application"
            required
          />
        </div>

        <div class="form-group">
          <label>Redirect URIs</label>
          <div v-for="(uri, index) in formData.redirect_uris" :key="index" class="uri-input">
            <input
              v-model="formData.redirect_uris[index]"
              type="url"
              placeholder="https://example.com/callback"
              required
            />
            <button
              v-if="formData.redirect_uris.length > 1"
              type="button"
              class="btn-remove"
              @click="removeRedirectUri(index)"
            >
              Remove
            </button>
          </div>
          <button type="button" class="btn-add" @click="addRedirectUri">
            + Add Redirect URI
          </button>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div class="form-actions">
          <button type="button" class="btn-secondary" @click="router.push('/clients')" :disabled="submitting">
            Cancel
          </button>
          <button type="submit" class="btn-primary" :disabled="submitting">
            {{ submitting ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </main>
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

h1 {
  color: #333;
  margin: 0;
}

.loading,
.error {
  padding: 2rem;
  text-align: center;
  background: #f5f5f5;
  border-radius: 8px;
}

.error {
  background: #ffebee;
  color: #d32f2f;
}

.form-container {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.info-section {
  background: #f5f5f5;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
}

.info-section h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
}

.info-item {
  margin-bottom: 1rem;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-item label {
  display: block;
  font-weight: 500;
  color: #555;
  margin-bottom: 0.25rem;
  font-size: 0.875rem;
}

.info-item span {
  color: #666;
}

.copy-field {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.copy-field code {
  flex: 1;
  padding: 0.5rem;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.875rem;
  word-break: break-all;
}

.btn-copy {
  padding: 0.5rem 1rem;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 0.875rem;
}

.btn-copy:hover {
  background: #5a6268;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #555;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: #007bff;
}

.uri-input {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.uri-input input {
  flex: 1;
}

.btn-remove {
  padding: 0.75rem 1rem;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
}

.btn-remove:hover {
  background: #c82333;
}

.btn-add {
  padding: 0.5rem 1rem;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-add:hover {
  background: #218838;
}

.error-message {
  padding: 0.75rem;
  background: #ffebee;
  color: #d32f2f;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #0056b3;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #5a6268;
}

.btn-primary:disabled,
.btn-secondary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
