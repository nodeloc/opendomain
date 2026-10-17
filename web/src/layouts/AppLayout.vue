<template>
  <div class="flex min-h-screen flex-col">
    <header class="sticky top-0 z-40 border-b border-y-line bg-y-bg/85 backdrop-blur-md">
      <nav class="mx-auto flex min-h-16 w-full max-w-[1280px] flex-wrap items-center gap-x-6 gap-y-2 px-5 py-3 sm:px-6" :aria-label="$t('shell.appNav')">
        <router-link to="/" class="shrink-0" :aria-label="siteName">
          <YLogo class="h-[26px] w-auto" />
        </router-link>
        <span aria-hidden="true" class="hidden text-y-line-strong md:inline">/</span>

        <div class="hidden flex-1 items-center gap-1 md:flex">
          <router-link
            v-for="item in primary"
            :key="item.to"
            :to="item.to"
            class="rounded-sm px-2.5 py-1.5 text-sm font-medium text-y-muted transition-colors hover:text-y-text"
            :class="{ '!text-y-text': isActive(item) }"
            :aria-current="isActive(item) ? 'page' : undefined"
          >
            {{ item.label }}
          </router-link>
          <span aria-hidden="true" class="mx-2 h-4 w-px bg-y-line-strong"></span>
          <router-link
            v-for="item in secondary"
            :key="item.to"
            :to="item.to"
            class="rounded-sm px-2.5 py-1.5 text-sm font-medium text-y-muted transition-colors hover:text-y-text"
          >
            {{ item.label }}
          </router-link>
        </div>

        <div class="ml-auto flex items-center gap-1">
          <div class="hidden items-center gap-1 sm:flex">
            <LanguageSwitcher />
            <ThemeToggle />
          </div>
          <div class="dropdown dropdown-end ml-1">
            <button
              tabindex="0"
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full border border-y-line-strong bg-y-surface-2"
              :aria-label="$t('shell.account')"
            >
              <UserAvatar :src="user?.avatar" :name="user?.username" :size="32" />
            </button>
            <ul tabindex="0" class="dropdown-content menu z-[60] mt-2 w-60 rounded-md p-1.5">
              <li class="menu-title">
                <span class="truncate normal-case tracking-normal">{{ user?.username || '—' }}</span>
              </li>
              <li v-if="user?.email" class="px-3 pb-2 text-xs text-y-faint"><span class="truncate p-0">{{ user.email }}</span></li>
              <li><router-link to="/profile"><User :size="16" :stroke-width="1.5" />{{ $t('nav.profile') }}</router-link></li>
              <li v-if="user?.is_admin"><router-link to="/admin"><ShieldCheck :size="16" :stroke-width="1.5" />{{ $t('nav.admin') }}</router-link></li>
              <li><router-link to="/"><ArrowUpRight :size="16" :stroke-width="1.5" />{{ $t('shell.viewSite') }}</router-link></li>
              <li class="mt-1 border-t border-y-line pt-1">
                <a class="text-y-danger-text" @click="handleLogout"><LogOut :size="16" :stroke-width="1.5" />{{ $t('nav.logout') }}</a>
              </li>
            </ul>
          </div>
          <button
            type="button"
            class="btn btn-ghost btn-sm btn-square md:hidden"
            :aria-label="$t('shell.menu')"
            :aria-expanded="mobileOpen"
            @click="mobileOpen = !mobileOpen"
          >
            <X v-if="mobileOpen" :size="20" :stroke-width="1.5" />
            <Menu v-else :size="20" :stroke-width="1.5" />
          </button>
        </div>
      </nav>

      <div v-if="mobileOpen" class="border-t border-y-line md:hidden">
        <nav class="flex flex-col px-5 py-3">
          <router-link
            v-for="item in [...primary, ...secondary]"
            :key="item.to"
            :to="item.to"
            class="rounded-sm px-2 py-2.5 text-[15px] text-y-muted hover:text-y-text"
            :class="{ '!text-y-text': isActive(item) }"
          >
            {{ item.label }}
          </router-link>
          <div class="mt-2 flex items-center gap-1 border-t border-y-line pt-3 sm:hidden">
            <LanguageSwitcher />
            <ThemeToggle />
          </div>
        </nav>
      </div>
    </header>

    <main class="flex-1">
      <slot />
    </main>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Menu, X, User, ShieldCheck, LogOut, ArrowUpRight } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import { DOCS_PATH } from '../config/links'
import YLogo from '../components/brand/YLogo.vue'
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
import ThemeToggle from '../components/ui/ThemeToggle.vue'
import UserAvatar from '../components/ui/UserAvatar.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()

const mobileOpen = ref(false)
const user = computed(() => authStore.user)
const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')

const primary = computed(() => [
  { to: '/domains', label: t('nav.domains'), match: ['/domains'] },
  { to: '/dashboard', label: t('shell.overview') },
  { to: '/orders', label: t('nav.orders'), match: ['/orders', '/checkout'] },
  { to: '/invitations', label: t('nav.invitations') },
  { to: '/coupons', label: t('nav.coupons') },
])

const secondary = computed(() => [
  { to: DOCS_PATH, label: t('shell.docs') },
  { to: '/status', label: t('shell.status') },
])

const isActive = (item) => {
  const prefixes = item.match || [item.to]
  return prefixes.some((p) => route.path === p || route.path.startsWith(p + '/'))
}

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
