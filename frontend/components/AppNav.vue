<script setup lang="ts">
const { user, isAuthenticated, logout } = useAuth()
const router = useRouter()

const handleLogout = async () => {
  await logout()
  router.push('/')
}
</script>

<template>
  <nav>
    <ul>
      <template v-if="isAuthenticated">
        <li><NuxtLink to="/dashboard">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/support">Support</NuxtLink></li>
        <li><NuxtLink to="/profile/sessions">Sessions</NuxtLink></li>
        <li><NuxtLink to="/clients">Clients</NuxtLink></li>
      </template>
      <li v-else><NuxtLink to="/">Home</NuxtLink></li>
      <li class="spacer"></li>
      <li v-if="isAuthenticated && user">
        <span class="user-info">{{ user.firstName }} {{ user.lastName }}</span>
      </li>
      <li v-if="isAuthenticated">
        <button @click="handleLogout" class="logout-btn">Sign Out</button>
      </li>
      <template v-else>
        <li><NuxtLink to="/sign-up">Sign Up</NuxtLink></li>
        <li><NuxtLink to="/sign-in">Sign In</NuxtLink></li>
      </template>
    </ul>
  </nav>
</template>

<style scoped>
nav {
  background-color: #f5f5f5;
  padding: 1rem;
}

ul {
  list-style: none;
  display: flex;
  gap: 1.5rem;
  margin: 0;
  padding: 0;
}

.spacer {
  flex: 1;
}

a {
  text-decoration: none;
  color: #333;
  font-weight: 500;
}

a:hover {
  color: #0066cc;
}

.user-info {
  color: #333;
  font-weight: 500;
}

.logout-btn {
  background: none;
  border: none;
  color: #333;
  font-weight: 500;
  cursor: pointer;
  padding: 0;
  font-size: inherit;
}

.logout-btn:hover {
  color: #0066cc;
}
</style>
