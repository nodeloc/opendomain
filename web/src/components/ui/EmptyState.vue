<template>
  <div class="flex flex-col items-center justify-center px-6 py-14 text-center">
    <svg viewBox="0 0 200 96" class="mb-6 h-20 w-auto" aria-hidden="true">
      <path d="M40 48 H78 C96 48 96 18 118 18 H150 M78 48 H150 M78 48 C96 48 96 78 118 78 H150" fill="none" stroke="rgb(var(--y-line-strong))" stroke-width="2" stroke-linecap="round" />
      <circle cx="40" cy="48" r="5" fill="rgb(var(--y-text))" />
      <circle cx="152" cy="18" r="5" fill="rgb(var(--y-node))" />
      <circle cx="152" cy="48" r="5" fill="rgb(var(--y-node))" opacity="0.55" />
      <circle cx="152" cy="78" r="5" fill="rgb(var(--y-node))" opacity="0.3" />
    </svg>
    <h3 class="y-h3">{{ title }}</h3>
    <p v-if="description" class="mt-2 max-w-md text-sm leading-6 text-y-muted">{{ description }}</p>
    <div v-if="$slots.default" class="mt-6 flex flex-wrap items-center justify-center gap-2">
      <slot />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>
