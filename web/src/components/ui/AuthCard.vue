<template>
  <section class="relative isolate flex min-h-[calc(100vh-4rem)] items-center justify-center overflow-hidden px-4 py-16">
    <div class="y-dot-grid absolute inset-0 -z-20" aria-hidden="true"></div>
    <div class="y-glow absolute left-1/2 top-1/3 -z-10 h-[420px] w-[640px] -translate-x-1/2 -translate-y-1/2" aria-hidden="true"></div>
    <div class="w-full max-w-[420px]">
      <div class="mb-8 text-center">
        <YLogo variant="mark" class="mx-auto h-10 w-10" />
        <p v-if="eyebrow" class="y-eyebrow mt-6">{{ eyebrow }}</p>
        <h1 class="mt-3 text-[1.75rem] font-semibold leading-tight tracking-[-0.03em] text-y-text">{{ title }}</h1>
        <p v-if="description" class="mt-2 text-[15px] leading-6 text-y-muted">{{ description }}</p>
      </div>
      <div class="y-panel p-6 sm:p-7">
        <slot />
      </div>
      <div v-if="$slots.footer" class="mt-6 text-center text-sm text-y-muted">
        <slot name="footer" />
      </div>
    </div>
  </section>
</template>

<script setup>
import YLogo from '../brand/YLogo.vue'

defineProps({
  eyebrow: { type: String, default: '' },
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>
