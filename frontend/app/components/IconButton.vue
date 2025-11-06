<template>
  <button
    :class="['icon-button', `variant-${variant}`, `color-${color}`]"
    :type="type"
    @click="$emit('click', $event)"
  >
    <slot v-if="$slots.default"></slot>
    <font-awesome-icon v-else-if="icon" :icon="icon" class="w-4 h-4" />
  </button>
</template>

<script setup lang="ts">
interface Props {
  icon?: string | string[]
  variant?: 'filled' | 'gradient' | 'outlined' | 'text'
  color?: 'slate' | 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info'
  type?: 'button' | 'submit' | 'reset'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'filled',
  color: 'slate',
  type: 'button'
})

defineEmits<{
  click: [event: MouseEvent]
}>()
</script>

<style scoped>
@reference "tailwindcss";

.icon-button {
  @apply rounded-md p-2.5 text-center text-sm transition-all;
  @apply disabled:pointer-events-none disabled:opacity-50 disabled:shadow-none;
}

/* Filled variant */
.variant-filled.color-slate {
  @apply bg-slate-800 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-slate-700;
  @apply focus:bg-slate-700 focus:shadow-none;
  @apply active:bg-slate-700 active:shadow-none;
}

.variant-filled.color-primary {
  @apply bg-blue-600 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-blue-700;
  @apply focus:bg-blue-700 focus:shadow-none;
  @apply active:bg-blue-700 active:shadow-none;
}

.variant-filled.color-secondary {
  @apply bg-gray-600 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-gray-700;
  @apply focus:bg-gray-700 focus:shadow-none;
  @apply active:bg-gray-700 active:shadow-none;
}

.variant-filled.color-success {
  @apply bg-green-600 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-green-700;
  @apply focus:bg-green-700 focus:shadow-none;
  @apply active:bg-green-700 active:shadow-none;
}

.variant-filled.color-danger {
  @apply bg-red-600 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-red-700;
  @apply focus:bg-red-700 focus:shadow-none;
  @apply active:bg-red-700 active:shadow-none;
}

.variant-filled.color-warning {
  @apply bg-yellow-600 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-yellow-700;
  @apply focus:bg-yellow-700 focus:shadow-none;
  @apply active:bg-yellow-700 active:shadow-none;
}

.variant-filled.color-info {
  @apply bg-cyan-600 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-cyan-700;
  @apply focus:bg-cyan-700 focus:shadow-none;
  @apply active:bg-cyan-700 active:shadow-none;
}

/* Gradient variant */
.variant-gradient.color-slate {
  @apply bg-gradient-to-tr from-slate-800 to-slate-700 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-slate-700;
  @apply focus:bg-slate-700 focus:shadow-none;
  @apply active:bg-slate-700 active:shadow-none;
}

.variant-gradient.color-primary {
  @apply bg-gradient-to-tr from-blue-600 to-blue-500 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-blue-700;
  @apply focus:bg-blue-700 focus:shadow-none;
  @apply active:bg-blue-700 active:shadow-none;
}

.variant-gradient.color-secondary {
  @apply bg-gradient-to-tr from-gray-600 to-gray-500 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-gray-700;
  @apply focus:bg-gray-700 focus:shadow-none;
  @apply active:bg-gray-700 active:shadow-none;
}

.variant-gradient.color-success {
  @apply bg-gradient-to-tr from-green-600 to-green-500 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-green-700;
  @apply focus:bg-green-700 focus:shadow-none;
  @apply active:bg-green-700 active:shadow-none;
}

.variant-gradient.color-danger {
  @apply bg-gradient-to-tr from-red-600 to-red-500 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-red-700;
  @apply focus:bg-red-700 focus:shadow-none;
  @apply active:bg-red-700 active:shadow-none;
}

.variant-gradient.color-warning {
  @apply bg-gradient-to-tr from-yellow-600 to-yellow-500 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-yellow-700;
  @apply focus:bg-yellow-700 focus:shadow-none;
  @apply active:bg-yellow-700 active:shadow-none;
}

.variant-gradient.color-info {
  @apply bg-gradient-to-tr from-cyan-600 to-cyan-500 border border-transparent text-white shadow-sm;
  @apply hover:shadow-lg hover:bg-cyan-700;
  @apply focus:bg-cyan-700 focus:shadow-none;
  @apply active:bg-cyan-700 active:shadow-none;
}

/* Outlined variant */
.variant-outlined.color-slate {
  @apply border border-slate-300 shadow-sm text-slate-600;
  @apply hover:shadow-lg hover:text-white hover:bg-slate-800 hover:border-slate-800;
  @apply focus:text-white focus:bg-slate-800 focus:border-slate-800;
  @apply active:text-white active:bg-slate-800 active:border-slate-800;
}

.variant-outlined.color-primary {
  @apply border border-blue-300 shadow-sm text-blue-600;
  @apply hover:shadow-lg hover:text-white hover:bg-blue-600 hover:border-blue-600;
  @apply focus:text-white focus:bg-blue-600 focus:border-blue-600;
  @apply active:text-white active:bg-blue-600 active:border-blue-600;
}

.variant-outlined.color-secondary {
  @apply border border-gray-300 shadow-sm text-gray-600;
  @apply hover:shadow-lg hover:text-white hover:bg-gray-600 hover:border-gray-600;
  @apply focus:text-white focus:bg-gray-600 focus:border-gray-600;
  @apply active:text-white active:bg-gray-600 active:border-gray-600;
}

.variant-outlined.color-success {
  @apply border border-green-300 shadow-sm text-green-600;
  @apply hover:shadow-lg hover:text-white hover:bg-green-600 hover:border-green-600;
  @apply focus:text-white focus:bg-green-600 focus:border-green-600;
  @apply active:text-white active:bg-green-600 active:border-green-600;
}

.variant-outlined.color-danger {
  @apply border border-red-300 shadow-sm text-red-600;
  @apply hover:shadow-lg hover:text-white hover:bg-red-600 hover:border-red-600;
  @apply focus:text-white focus:bg-red-600 focus:border-red-600;
  @apply active:text-white active:bg-red-600 active:border-red-600;
}

.variant-outlined.color-warning {
  @apply border border-yellow-300 shadow-sm text-yellow-600;
  @apply hover:shadow-lg hover:text-white hover:bg-yellow-600 hover:border-yellow-600;
  @apply focus:text-white focus:bg-yellow-600 focus:border-yellow-600;
  @apply active:text-white active:bg-yellow-600 active:border-yellow-600;
}

.variant-outlined.color-info {
  @apply border border-cyan-300 shadow-sm text-cyan-600;
  @apply hover:shadow-lg hover:text-white hover:bg-cyan-600 hover:border-cyan-600;
  @apply focus:text-white focus:bg-cyan-600 focus:border-cyan-600;
  @apply active:text-white active:bg-cyan-600 active:border-cyan-600;
}

/* Text variant */
.variant-text.color-slate {
  @apply border border-transparent text-slate-600;
  @apply hover:bg-slate-100;
  @apply focus:bg-slate-100;
  @apply active:bg-slate-100;
}

.variant-text.color-primary {
  @apply border border-transparent text-blue-600;
  @apply hover:bg-blue-50;
  @apply focus:bg-blue-50;
  @apply active:bg-blue-50;
}

.variant-text.color-secondary {
  @apply border border-transparent text-gray-600;
  @apply hover:bg-gray-100;
  @apply focus:bg-gray-100;
  @apply active:bg-gray-100;
}

.variant-text.color-success {
  @apply border border-transparent text-green-600;
  @apply hover:bg-green-50;
  @apply focus:bg-green-50;
  @apply active:bg-green-50;
}

.variant-text.color-danger {
  @apply border border-transparent text-red-600;
  @apply hover:bg-red-50;
  @apply focus:bg-red-50;
  @apply active:bg-red-50;
}

.variant-text.color-warning {
  @apply border border-transparent text-yellow-600;
  @apply hover:bg-yellow-50;
  @apply focus:bg-yellow-50;
  @apply active:bg-yellow-50;
}

.variant-text.color-info {
  @apply border border-transparent text-cyan-600;
  @apply hover:bg-cyan-50;
  @apply focus:bg-cyan-50;
  @apply active:bg-cyan-50;
}
</style>
