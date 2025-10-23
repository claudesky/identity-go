<script setup lang="ts">
import { useClients, type CreateClientRequest } from '~/composables/useClients'

const emit = defineEmits<{
  close: []
  created: [clientSecret: string]
}>()

const { createClient } = useClients()

const formData = reactive<CreateClientRequest>({
  name: '',
  redirect_uris: ['']
})

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
    const result = await createClient({
      name: formData.name.trim(),
      redirect_uris: validUris
    })

    if (result.error.value) {
      throw result.error.value
    }

    if (result.data.value?.data) {
      emit('created', result.data.value.data.client_secret)
    }
  } catch (err: any) {
    console.error('Failed to create client:', err)
    error.value = err.message || 'Failed to create client'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="modal-overlay" @click="emit('close')">
    <div class="modal" @click.stop>
      <h2>Create OAuth Client</h2>

      <form @submit.prevent="handleSubmit">
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
          <button type="button" class="btn-secondary" @click="emit('close')" :disabled="submitting">
            Cancel
          </button>
          <button type="submit" class="btn-primary" :disabled="submitting">
            {{ submitting ? 'Creating...' : 'Create Client' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #333;
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
