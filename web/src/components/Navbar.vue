<template>
  <header class="sticky top-0 z-50 border-b border-y-line bg-y-bg/85 backdrop-blur-md">
    <div class="y-container flex h-16 items-center gap-6">
      <router-link to="/" class="shrink-0" :aria-label="siteName">
        <YLogo class="h-[26px] w-auto" />
      </router-link>

      <nav class="hidden items-center gap-1 lg:flex">
        <router-link
          v-for="item in links"
          :key="item.to"
          :to="item.to"
          class="rounded-sm px-3 py-1.5 text-sm text-y-muted transition-colors hover:text-y-text"
          active-class="!text-y-text"
          :exact-active-class="item.exact ? '!text-y-text' : ''"
        >
          {{ item.label }}
        </router-link>
      </nav>

      <div class="ml-auto flex items-center gap-1">
        <div class="hidden items-center gap-1 sm:flex">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>

        <template v-if="!isAuthenticated">
          <router-link to="/login" class="btn btn-ghost btn-sm hidden sm:inline-flex">{{ $t('shell.signIn') }}</router-link>
          <router-link to="/register" class="btn btn-primary btn-sm ml-1">{{ $t('shell.claimName') }}</router-link>
        </template>

        <template v-else>
          <router-link to="/dashboard" class="btn btn-primary btn-sm ml-1 hidden sm:inline-flex">
            {{ $t('shell.openDashboard') }}
          </router-link>
          <div class="dropdown dropdown-end ml-1">
            <button tabindex="0" type="button" class="btn btn-ghost btn-sm btn-square" :aria-label="user?.username || 'Account'">
              <UserAvatar :src="user?.avatar" :name="user?.username" :size="28" />
            </button>
            <ul tabindex="0" class="dropdown-content menu z-[60] mt-2 w-56 rounded-md p-1.5">
              <li class="menu-title"><span class="truncate">{{ user?.username }}</span></li>
              <li><router-link to="/dashboard"><LayoutDashboard :size="16" :stroke-width="1.5" />{{ $t('nav.dashboard') }}</router-link></li>
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
          class="btn btn-ghost btn-sm btn-square lg:hidden"
          :aria-label="$t('shell.menu')"
          :aria-expanded="mobileOpen"
          @click="mobileOpen = !mobileOpen"
        >
          <X v-if="mobileOpen" :size="20" :stroke-width="1.5" />
          <Menu v-else :size="20" :stroke-width="1.5" />
        </button>
      </div>
    </div>

    <div v-if="mobileOpen" class="border-t border-y-line bg-y-bg lg:hidden">
      <nav class="y-container flex flex-col py-3">
        <router-link
          v-for="item in links"
          :key="item.to"
          :to="item.to"
          class="rounded-sm px-2 py-2.5 text-[15px] text-y-muted hover:text-y-text"
        >
          {{ item.label }}
        </router-link>
        <router-link v-if="isAuthenticated" to="/dashboard" class="rounded-sm px-2 py-2.5 text-[15px] text-y-mint-text">
          {{ $t('shell.openDashboard') }}
        </router-link>
        <router-link v-else to="/login" class="rounded-sm px-2 py-2.5 text-[15px] text-y-muted hover:text-y-text">
          {{ $t('shell.signIn') }}
        </router-link>
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
import { Menu, X, LayoutDashboard, Globe, User, ShieldCheck, LogOut } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
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
  { to: '/', label: t('nav.home'), exact: true },
  { to: '/announcements', label: t('nav.announcements') },
  { to: '/domain-health', label: t('nav.domainHealth') },
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
