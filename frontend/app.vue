<script setup lang="ts">
const config = useRuntimeConfig()

useHead({
  titleTemplate: (titleChunk) => {
    return titleChunk
      ? `${titleChunk} - ${config.public.appName}`
      : config.public.appName
  },
})
</script>
<template>
  <NuxtLoadingIndicator />
  <NuxtPage />
</template>
