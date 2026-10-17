<template>
  <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
    <div class="min-w-0 space-y-2">
      <p v-if="eyebrow" class="y-eyebrow">{{ eyebrow }}</p>
      <h1 class="text-[1.75rem] font-semibold leading-tight tracking-[-0.03em] text-y-text sm:text-[2rem]">
        <slot name="title">{{ title }}</slot>
      </h1>
      <p v-if="description || $slots.description" class="max-w-2xl text-[15px] leading-6 text-y-muted">
        <slot name="description">{{ description }}</slot>
      </p>
    </div>
    <div v-if="$slots.actions" class="flex shrink-0 flex-wrap items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
defineProps({
  eyebrow: { type: String, default: '' },
  title: { type: String, default: '' },
  description: { type: String, default: '' },
})
</script>
