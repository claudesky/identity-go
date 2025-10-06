<script setup lang="ts">

definePageMeta({
  middleware: 'auth'
})

</script>

<template>
  <div>
    <AppNav />
    <main>
      <h1>Welcome to Identity - Go</h1>
      <sub>Support Page</sub>
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
