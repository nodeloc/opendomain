<template>
  <header class="sticky top-0 z-50 border-b border-y-line-strong/70 bg-y-surface/90 backdrop-blur-md">
    <nav class="flex min-h-16 flex-wrap items-center gap-x-7 gap-y-3 px-5 py-3 sm:px-10" :aria-label="$t('shell.mainNav')">
      <div class="flex min-w-0 flex-1 items-center gap-5">
        <router-link to="/" class="shrink-0" :aria-label="siteName">
          <YLogo class="h-[30px] w-auto" />
        </router-link>
      </div>

      <div class="hidden items-center gap-6 text-sm font-semibold lg:flex">
        <template v-for="item in links" :key="item.label">
          <a v-if="item.href" :href="item.href" target="_blank" rel="noopener" class="text-y-muted transition-colors hover:text-y-text">{{ item.label }}</a>
          <router-link v-else :to="item.to" class="text-y-muted transition-colors hover:text-y-text" active-class="!text-y-text">{{ item.label }}</router-link>
        </template>
      </div>

      <div class="flex items-center gap-1">
        <div class="hidden items-center gap-1 sm:flex">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>

        <router-link v-if="!isAuthenticated" to="/login" class="y-btn-ink ml-2">{{ $t('shell.getStarted') }}</router-link>

        <template v-else>
          <router-link to="/domains" class="y-btn-ink ml-2 hidden sm:inline-flex">{{ $t('shell.openDashboard') }}</router-link>
          <div class="dropdown dropdown-end ml-1">
            <button tabindex="0" type="button" class="btn btn-ghost btn-sm btn-square" :aria-label="user?.username || $t('shell.account')">
              <UserAvatar :src="user?.avatar" :name="user?.username" :size="28" />
            </button>
            <ul tabindex="0" class="dropdown-content menu z-[60] mt-2 w-56 rounded-md p-1.5">
              <li class="menu-title"><span class="truncate">{{ user?.username }}</span></li>
              <li><router-link to="/domains"><Globe :size="16" :stroke-width="1.5" />{{ $t('nav.domains') }}</router-link></li>
              <li><router-link to="/profile"><User :size="16" :stroke-width="1.5" />{{ $t('nav.profile') }}</router-link></li>
              <li v-if="user?.is_admin"><router-link to="/admin"><ShieldCheck :size="16" :stroke-width="1.5" />{{ $t('nav.admin') }}</router-link></li>
              <li class="mt-1 border-t border-y-line pt-1">
                <a class="text-y-danger-text" @click="handleLogout"><LogOut :size="16" :stroke-width="1.5" />{{ $t('nav.logout') }}</a>
              </li>
            </ul>
          </div>
        </template>

        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square ml-1 lg:hidden"
          :aria-label="$t('shell.menu')"
          :aria-expanded="mobileOpen"
          @click="mobileOpen = !mobileOpen"
        >
          <X v-if="mobileOpen" :size="20" :stroke-width="1.5" />
          <Menu v-else :size="20" :stroke-width="1.5" />
        </button>
      </div>
    </nav>

    <div v-if="mobileOpen" class="border-t border-y-line bg-y-bg lg:hidden">
      <nav class="flex flex-col px-5 py-3">
        <template v-for="item in links" :key="item.label">
          <a v-if="item.href" :href="item.href" target="_blank" rel="noopener" class="rounded-sm px-2 py-2.5 text-[15px] text-y-muted hover:text-y-text">{{ item.label }}</a>
          <router-link v-else :to="item.to" class="rounded-sm px-2 py-2.5 text-[15px] text-y-muted hover:text-y-text">{{ item.label }}</router-link>
        </template>
        <div class="mt-2 flex items-center gap-1 border-t border-y-line pt-3 sm:hidden">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Menu, X, Globe, User, ShieldCheck, LogOut } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import { DOCS_PATH, GITHUB_URL } from '../config/links'
import YLogo from './brand/YLogo.vue'
import LanguageSwitcher from './LanguageSwitcher.vue'
import ThemeToggle from './ui/ThemeToggle.vue'
import UserAvatar from './ui/UserAvatar.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const user = computed(() => authStore.user)
const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const mobileOpen = ref(false)

const links = computed(() => [
  { to: DOCS_PATH, label: t('shell.docs') },
  { to: '/status', label: t('shell.status') },
  { to: '/announcements', label: t('nav.announcements') },
  { to: '/domain-health', label: t('nav.domainHealth') },
  { href: GITHUB_URL, label: t('shell.github') },
])

watch(() => route.fullPath, () => {
  mobileOpen.value = false
})

onMounted(() => {
  if (authStore.isAuthenticated && !authStore.user) authStore.fetchProfile()
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
