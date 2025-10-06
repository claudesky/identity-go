<script setup lang="ts">
import { useAuthAPI } from '~/lib/api/auth-api'

definePageMeta({
  middleware: 'guest'
})

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const firstName = ref('')
const lastName = ref('')
const error = ref('')
const loading = ref(false)

const { signUp } = useAuthAPI()
const { setUser } = useAuth()
const router = useRouter()

const validateForm = (): boolean => {
  if (!firstName.value || !lastName.value || !email.value || !password.value || !confirmPassword.value) {
    error.value = 'Please fill in all fields'
    return false
  }

  if (password.value.length < 8) {
    error.value = 'Password must be at least 8 characters long'
    return false
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return false
  }

  return true
}

const handleSubmit = async () => {
  error.value = ''

  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    const response = await signUp({
      firstName: firstName.value,
      lastName: lastName.value,
      email: email.value,
      password: password.value
    })

    if (response.user) {
      setUser(response.user)
      await router.push('/')
    } else {
      error.value = response.message || 'Sign up failed'
    }
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Sign up failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <AppNav />
    <div class="page-container">
      <div class="container">
        <h1>Sign Up</h1>
        <div class="form">
          <div v-if="error" class="error">{{ error }}</div>

          <form @submit.prevent="handleSubmit">
            <div class="form-row">
              <div class="form-group">
                <label for="firstName">First Name</label>
                <input
                  type="text"
                  id="firstName"
                  v-model="firstName"
                  required
                  autocomplete="given-name"
                />
              </div>

              <div class="form-group">
                <label for="lastName">Last Name</label>
                <input
                  type="text"
                  id="lastName"
                  v-model="lastName"
                  required
                  autocomplete="family-name"
                />
              </div>
            </div>

            <div class="form-group">
              <label for="email">Email</label>
              <input
                type="email"
                id="email"
                v-model="email"
                required
                autocomplete="email"
              />
            </div>

            <div class="form-group">
              <label for="password">Password</label>
              <input
                type="password"
                id="password"
                v-model="password"
                required
                autocomplete="new-password"
              />
              <div class="password-hint">Must be at least 8 characters</div>
            </div>

            <div class="form-group">
              <label for="confirmPassword">Confirm Password</label>
              <input
                type="password"
                id="confirmPassword"
                v-model="confirmPassword"
                required
                autocomplete="new-password"
              />
            </div>

            <button type="submit" :disabled="loading">
              {{ loading ? 'Creating account...' : 'Sign Up' }}
            </button>
          </form>

          <div class="link">
            Already have an account? <NuxtLink to="/sign-in">Sign in</NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
  padding: 2rem;
}

.container {
  max-width: 400px;
  width: 100%;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

.form {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-weight: 500;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

input:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
}

button {
  width: 100%;
  padding: 0.75rem;
  background: #0066cc;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

button:hover:not(:disabled) {
  background: #0052a3;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.error {
  color: #d32f2f;
  margin-bottom: 1rem;
  padding: 0.75rem;
  background: #ffebee;
  border-radius: 4px;
  font-size: 0.9rem;
}

.link {
  text-align: center;
  margin-top: 1.5rem;
}

.link a {
  color: #0066cc;
  text-decoration: none;
}

.link a:hover {
  text-decoration: underline;
}

.password-hint {
  font-size: 0.85rem;
  color: #666;
  margin-top: 0.25rem;
}
</style>
