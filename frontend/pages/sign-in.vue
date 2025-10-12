<script setup lang="ts">

definePageMeta({
  middleware: 'guest',
})

useHead({
  title: 'Sign In'
})

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const emailInput = ref<HTMLInputElement | null>(null)

const { login } = useAuth()
const router = useRouter()

onMounted(() => {
  emailInput.value?.focus()
})

const handleSubmit = async () => {
  error.value = ''
  loading.value = true

  try {
    await login({
      email: email.value,
      password: password.value
    })

    await router.push('/dashboard')
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Sign in failed'
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
        <h1>Sign In</h1>
        <div class="form">
          <div v-if="error" class="error">{{ error }}</div>

          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="email">Email</label>
              <input
                ref="emailInput"
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
                autocomplete="current-password"
              />
            </div>

            <button type="submit" :disabled="loading">
              {{ loading ? 'Signing in...' : 'Sign In' }}
            </button>
          </form>

          <div class="link">
            Don't have an account? <NuxtLink to="/sign-up">Sign up</NuxtLink>
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
</style>
