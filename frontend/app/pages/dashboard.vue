<script setup lang="ts">
const config = useRuntimeConfig()

definePageMeta({
  middleware: 'auth'
})

useHead({
  title: 'Dashboard'
})
</script>

<template>
  <div>
    <AppNav />
    <main>
      <h1>Welcome to {{ config.public.appName }}</h1>
      <sub>Dashboard Page</sub>
    </main>
  </div>
</template>

<style scoped>
main {
  display: block;
  padding: 2rem;
  text-align: center;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

sub {
  font-size: 1rem;
  color: #666;
}
</style>
