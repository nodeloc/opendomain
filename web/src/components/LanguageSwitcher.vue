<template>
  <div class="dropdown dropdown-end">
    <label tabindex="0" class="btn btn-ghost btn-circle">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
      </svg>
    </label>
    <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-40 mt-3">
      <li>
        <a @click="changeLocale('zh-CN')" :class="{ 'active': currentLocale === 'zh-CN' }">
          <span class="text-xl">🇨🇳</span>
          <span>简体中文</span>
        </a>
      </li>
      <li>
        <a @click="changeLocale('en-US')" :class="{ 'active': currentLocale === 'en-US' }">
          <span class="text-xl">🇺🇸</span>
          <span>English</span>
        </a>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

const currentLocale = computed(() => locale.value)

const changeLocale = (newLocale) => {
  locale.value = newLocale
  localStorage.setItem('locale', newLocale)
}
</script>

<style scoped>
.active {
  @apply bg-primary text-primary-content;
}
</style>
