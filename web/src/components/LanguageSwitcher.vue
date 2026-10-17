<template>
  <div class="dropdown dropdown-end">
    <button tabindex="0" type="button" class="btn btn-ghost btn-sm gap-1.5 px-2" :aria-label="$t('nav.language')">
      <Languages :size="18" :stroke-width="1.5" />
      <span class="font-mono text-xs uppercase text-y-muted">{{ currentLocale === 'zh-CN' ? '中' : 'EN' }}</span>
    </button>
    <ul tabindex="0" class="dropdown-content menu z-[60] mt-2 w-40 rounded-md p-1.5">
      <li>
        <a :class="{ active: currentLocale === 'zh-CN' }" @click="changeLocale('zh-CN')">
          <span class="w-6 font-mono text-xs text-y-faint">中</span>简体中文
        </a>
      </li>
      <li>
        <a :class="{ active: currentLocale === 'en-US' }" @click="changeLocale('en-US')">
          <span class="w-6 font-mono text-xs text-y-faint">EN</span>English
        </a>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Languages } from 'lucide-vue-next'

const { locale } = useI18n()

const currentLocale = computed(() => locale.value)

const changeLocale = (newLocale) => {
  locale.value = newLocale
  localStorage.setItem('locale', newLocale)
  document.documentElement.setAttribute('lang', newLocale)
  if (document.activeElement instanceof HTMLElement) document.activeElement.blur()
}
</script>
