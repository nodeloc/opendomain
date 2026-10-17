<template>
  <div class="flex min-h-screen flex-col">
    <header class="sticky top-0 z-40 border-b border-y-line bg-y-bg/95 backdrop-blur-md">
      <nav
        class="mx-auto flex w-full max-w-[1360px] flex-wrap items-center gap-x-6 gap-y-3 px-4 py-3 sm:px-6 lg:h-16 lg:flex-nowrap lg:py-0"
        :aria-label="$t('shell.mainNav')"
      >
        <router-link to="/" class="flex shrink-0 items-center gap-3" :aria-label="siteName">
          <YLogo class="h-[28px] w-auto sm:h-[30px]" />
          <span class="y-site-badge">{{ $t('docs.badge') }}</span>
        </router-link>

        <div class="order-last w-full lg:order-none lg:w-auto lg:max-w-[420px] lg:flex-1">
          <DocsSearch :items="searchable" />
        </div>

        <div class="ml-auto flex items-center gap-1">
          <div class="mr-3 hidden items-center gap-[22px] text-sm font-medium lg:flex">
            <router-link to="/status" class="text-y-muted transition-colors hover:text-y-text">{{ $t('shell.status') }}</router-link>
            <a :href="GITHUB_URL" target="_blank" rel="noopener" class="text-y-muted transition-colors hover:text-y-text">{{ $t('shell.github') }}</a>
            <a :href="TELEGRAM_GROUP_URL" target="_blank" rel="noopener" class="text-y-muted transition-colors hover:text-y-text">{{ $t('shell.telegram') }}</a>
          </div>
          <div class="hidden items-center gap-1 sm:flex">
            <LanguageSwitcher />
            <ThemeToggle />
          </div>
          <router-link v-if="isAuthenticated" to="/domains" class="y-btn-ink ml-2 !min-h-[2.25rem] px-4">{{ $t('docs.dashboard') }}</router-link>
          <router-link v-else to="/login" class="y-btn-ink ml-2 !min-h-[2.25rem] px-4">{{ $t('shell.getStarted') }}</router-link>
        </div>
      </nav>

      <div class="border-t border-y-line lg:hidden">
        <button
          type="button"
          class="mx-auto flex min-h-[2.75rem] w-full max-w-[1360px] items-center gap-2 px-4 text-left text-sm sm:px-6"
          :aria-expanded="menuOpen"
          aria-controls="docs-mobile-nav"
          @click="menuOpen = !menuOpen"
        >
          <PanelLeft :size="16" :stroke-width="1.5" class="shrink-0 text-y-faint" />
          <span class="shrink-0 text-y-muted">{{ $t('docs.menu') }}</span>
          <template v-if="currentTitle">
            <ChevronRight :size="14" :stroke-width="1.5" class="shrink-0 text-y-faint" />
            <span class="truncate text-y-text">{{ currentTitle }}</span>
          </template>
          <ChevronDown :size="16" :stroke-width="1.5" class="ml-auto shrink-0 text-y-faint transition-transform" :class="{ 'rotate-180': menuOpen }" />
        </button>
        <div v-if="menuOpen" id="docs-mobile-nav" class="max-h-[70vh] overflow-y-auto border-t border-y-line px-2 py-5 sm:px-4">
          <DocsSidebar :groups="groups" :policies="policies" :loading="loading" :failed="failed" @retry="load(true)" />
          <div class="mt-5 flex flex-wrap items-center gap-x-5 gap-y-2 border-t border-y-line px-3 pt-4 text-sm">
            <router-link to="/status" class="text-y-muted hover:text-y-text">{{ $t('shell.status') }}</router-link>
            <a :href="GITHUB_URL" target="_blank" rel="noopener" class="text-y-muted hover:text-y-text">{{ $t('shell.github') }}</a>
            <a :href="TELEGRAM_GROUP_URL" target="_blank" rel="noopener" class="text-y-muted hover:text-y-text">{{ $t('shell.telegram') }}</a>
            <div class="ml-auto flex items-center gap-1 sm:hidden">
              <LanguageSwitcher />
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>
    </header>

    <div class="mx-auto flex w-full max-w-[1360px] flex-1 items-stretch gap-12 px-4 sm:px-6">
      <aside class="hidden w-[240px] shrink-0 lg:block">
        <div class="docs-sidebar sticky top-16 overflow-y-auto py-8">
          <DocsSidebar :groups="groups" :policies="policies" :loading="loading" :failed="failed" @retry="load(true)" />
        </div>
      </aside>
      <main class="min-w-0 flex-1">
        <slot />
      </main>
    </div>

    <Footer />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { PanelLeft, ChevronDown, ChevronRight } from 'lucide-vue-next'
import Footer from '../components/Footer.vue'
import YLogo from '../components/brand/YLogo.vue'
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
import ThemeToggle from '../components/ui/ThemeToggle.vue'
import DocsSearch from '../components/docs/DocsSearch.vue'
import DocsSidebar from '../components/docs/DocsSidebar.vue'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import { useDocsPages } from '../composables/useDocsPages'
import { GITHUB_URL, TELEGRAM_GROUP_URL } from '../config/links'

const route = useRoute()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const { groups, sidebarGroups, searchable, loading, failed, load } = useDocsPages()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const policies = computed(() => sidebarGroups.value[sidebarGroups.value.length - 1])
const menuOpen = ref(false)

const currentTitle = computed(() => searchable.value.find((item) => item.to === route.path)?.title || '')

watch(() => route.fullPath, () => {
  menuOpen.value = false
})

onMounted(() => {
  load()
})
</script>

<style scoped>
.docs-sidebar {
  max-height: calc(100vh - 4rem);
}
</style>
