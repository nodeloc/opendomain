<template>
  <div
    class="relative isolate flex flex-col overflow-hidden"
    :class="standalone ? 'min-h-screen' : 'min-h-[calc(100vh-4rem)]'"
  >
    <div class="y-dot-grid absolute inset-0 -z-20" aria-hidden="true"></div>
    <div class="y-glow absolute left-1/2 top-1/2 -z-10 h-[560px] w-[760px] max-w-[160vw] -translate-x-1/2 -translate-y-[55%]" aria-hidden="true"></div>

    <div v-if="standalone" class="flex items-center justify-between gap-3 px-4 pt-4 sm:px-6 sm:pt-5">
      <router-link to="/" class="inline-flex min-h-[2.25rem] items-center gap-1.5 rounded-sm px-2 text-sm font-medium text-y-muted transition-colors hover:text-y-text">
        <ArrowLeft :size="16" :stroke-width="1.5" />
        {{ $t('auth2.home') }}
      </router-link>
      <div class="flex items-center gap-1">
        <LanguageSwitcher />
        <ThemeToggle />
      </div>
    </div>

    <main class="flex flex-1 items-center justify-center px-4 py-12 sm:py-16">
      <div class="flex w-full max-w-[400px] flex-col items-center gap-7">
        <router-link to="/" class="inline-flex rounded-md" :aria-label="$t('auth2.home')">
          <span class="inline-flex h-14 w-14 items-center justify-center rounded-[14px] border border-y-line bg-y-surface">
            <YLogo variant="mark" class="h-9 w-9" />
          </span>
        </router-link>

        <div class="flex flex-col gap-1.5 text-center">
          <p v-if="eyebrow" class="y-eyebrow mb-1">{{ eyebrow }}</p>
          <h1 class="text-[28px] font-semibold leading-[34px] tracking-[-0.03em] text-y-text">{{ title }}</h1>
          <p v-if="$slots.description || description" class="text-sm leading-[22px] text-y-muted">
            <slot name="description">{{ description }}</slot>
          </p>
        </div>

        <div class="auth-surface w-full rounded-[18px] border border-y-line p-5 sm:p-6">
          <slot />
        </div>

        <div v-if="$slots.footer" class="flex flex-col items-center gap-2 text-center">
          <slot name="footer" />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ArrowLeft } from 'lucide-vue-next'
import YLogo from '../brand/YLogo.vue'
import LanguageSwitcher from '../LanguageSwitcher.vue'
import ThemeToggle from './ThemeToggle.vue'

defineProps({
  eyebrow: { type: String, default: '' },
  title: { type: String, required: true },
  description: { type: String, default: '' },
  // Pages rendered without the site navbar need their own way home and the theme/language controls.
  standalone: { type: Boolean, default: true },
})
</script>

<style scoped>
.auth-surface {
  background: rgb(var(--y-surface) / 0.85);
  -webkit-backdrop-filter: blur(10px);
  backdrop-filter: blur(10px);
}
</style>
