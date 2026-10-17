<template>
  <button
    type="button"
    class="btn btn-ghost btn-sm btn-square"
    :title="isDark ? $t('nav.lightMode') : $t('nav.darkMode')"
    :aria-label="isDark ? $t('nav.lightMode') : $t('nav.darkMode')"
    @click="toggleTheme"
  >
    <Sun v-if="isDark" :size="18" :stroke-width="1.5" />
    <Moon v-else :size="18" :stroke-width="1.5" />
  </button>
</template>

<script setup>
import { Sun, Moon } from 'lucide-vue-next'
import { useTheme } from '../../composables/useTheme'

const { isDark, toggleTheme } = useTheme()
</script>
