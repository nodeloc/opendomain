<template>
  <footer class="border-t border-y-line">
    <div class="y-container flex flex-wrap items-center justify-between gap-x-6 gap-y-4 py-8 text-sm text-y-faint">
      <span>© {{ year }} {{ siteName }}</span>
      <nav class="flex flex-wrap gap-x-6 gap-y-2" :aria-label="$t('shell.footerNav')">
        <template v-for="item in links" :key="item.label">
          <a v-if="item.href" :href="item.href" target="_blank" rel="noopener" class="text-y-muted transition-colors hover:text-y-text">{{ item.label }}</a>
          <router-link v-else :to="item.to" class="text-y-muted transition-colors hover:text-y-text">{{ item.label }}</router-link>
        </template>
      </nav>
    </div>
  </footer>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSiteConfigStore } from '../stores/siteConfig'
import { DOCS_PATH, GITHUB_URL, TELEGRAM_GROUP_URL } from '../config/links'

const { t } = useI18n()
const siteConfigStore = useSiteConfigStore()
const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const year = new Date().getFullYear()

const links = computed(() => [
  { to: DOCS_PATH, label: t('shell.docs') },
  { to: '/status', label: t('shell.status') },
  { to: '/sponsor', label: t('sponsor.nav') },
  { href: GITHUB_URL, label: t('shell.github') },
  { href: TELEGRAM_GROUP_URL, label: t('shell.telegram') },
  { to: '/acceptable-use', label: t('shell.acceptableUse') },
  { to: '/report', label: t('shell.reportAbuse') },
  { to: '/pages/terms-of-service', label: t('shell.terms') },
  { to: '/pages/privacy-policy', label: t('shell.privacy') },
])
</script>
