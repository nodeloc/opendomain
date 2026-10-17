<template>
  <section class="relative isolate flex min-h-[calc(100vh-4rem)] flex-col items-center justify-center overflow-hidden px-4 pb-56 pt-16 text-center">
    <img
      src="/brand/planet-horizon-2880x1800.jpg"
      alt=""
      aria-hidden="true"
      class="absolute inset-0 -z-10 h-full w-full object-cover object-bottom"
      :class="isDark ? 'opacity-100' : 'opacity-0'"
    />
    <div v-if="!isDark" class="y-dot-grid absolute inset-0 -z-10" aria-hidden="true"></div>
    <p class="font-mono text-sm tracking-[0.1em] text-y-mint-text">ERROR 404</p>
    <h1 class="y-headline mt-5 text-[clamp(56px,12vw,160px)] [&:lang(zh)]:text-[clamp(44px,8.5vw,112px)]">
      {{ $t('notFound.titleLead') }} <em>{{ $t('notFound.titleAccent') }}</em>
    </h1>
    <p class="mt-6 max-w-md text-lg leading-7 text-y-muted">
      {{ $t('notFound.before') }}
      <span class="font-mono text-y-text">{{ path }}</span>
      {{ $t('notFound.after') }}
    </p>
    <div class="mt-8 flex flex-wrap justify-center gap-3">
      <router-link to="/dashboard" class="btn btn-primary">{{ $t('shell.openDashboard') }}</router-link>
      <router-link to="/" class="btn btn-outline">{{ $t('shell.claimName') }}</router-link>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTheme } from '../composables/useTheme'

const route = useRoute()
const { isDark } = useTheme()
const path = computed(() => route.fullPath)
</script>
