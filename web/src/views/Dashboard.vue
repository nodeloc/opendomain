<template>
  <div class="y-page">
    <PageHeader :eyebrow="$t('uiUser.dashboard.eyebrow')" :description="$t('uiUser.dashboard.description')">
      <template #title>
        {{ user?.username ? $t('uiUser.dashboard.greeting', { name: user.username }) : $t('uiUser.dashboard.greetingFallback') }}
      </template>
      <template #actions>
        <router-link to="/" class="btn btn-primary gap-1.5">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('uiUser.dashboard.claim') }}
        </router-link>
      </template>
    </PageHeader>

    <!-- Stats -->
    <section class="y-panel grid grid-cols-1 gap-px overflow-hidden bg-y-line sm:grid-cols-2 lg:grid-cols-4">
      <div class="bg-y-surface px-5 py-6">
        <p class="y-label">{{ $t('uiUser.dashboard.domains') }}</p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">
          {{ domainCount }}<span class="text-lg text-y-faint"> / {{ user?.domain_quota ?? '—' }}</span>
        </p>
        <p class="mt-2 text-sm text-y-muted">{{ $t('uiUser.dashboard.slotsLeft', { n: slotsLeft }) }}</p>
      </div>

      <div class="bg-y-surface px-5 py-6">
        <p class="y-label">{{ $t('uiUser.dashboard.level') }}</p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">{{ levelLabel }}</p>
        <p class="mt-2 text-sm text-y-muted">{{ $t('uiUser.dashboard.levelHint') }}</p>
      </div>

      <div class="bg-y-surface px-5 py-6">
        <p class="y-label">{{ $t('uiUser.dashboard.email') }}</p>
        <p class="mt-3 truncate text-[15px] text-y-text" :title="user?.email">{{ user?.email || '—' }}</p>
        <div class="mt-2">
          <StatusDot
            :tone="user?.email_verified ? 'ok' : 'warn'"
            :label="user?.email_verified ? $t('uiUser.dashboard.verified') : $t('uiUser.dashboard.notVerified')"
          />
        </div>
      </div>

      <div class="bg-y-surface px-5 py-6">
        <p class="y-label">{{ $t('uiUser.dashboard.inviteCode') }}</p>
        <div class="mt-3 flex items-center gap-2">
          <span class="min-w-0 truncate font-mono text-2xl font-medium tracking-tight text-y-text">{{ user?.invite_code || '—' }}</span>
          <button
            v-if="user?.invite_code"
            type="button"
            class="btn btn-ghost btn-sm btn-square shrink-0"
            :title="$t('uiUser.dashboard.copyCode')"
            :aria-label="$t('uiUser.dashboard.copyCode')"
            @click="copyInviteCode"
          >
            <Check v-if="codeCopied" :size="16" :stroke-width="1.5" class="text-y-mint-text" />
            <Copy v-else :size="16" :stroke-width="1.5" />
          </button>
        </div>
        <p class="mt-2 text-sm text-y-muted">{{ $t('uiUser.dashboard.inviteHint') }}</p>
      </div>
    </section>

    <!-- Shortcuts -->
    <section class="space-y-4">
      <h2 class="y-label">{{ $t('uiUser.dashboard.shortcuts') }}</h2>
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <router-link v-for="tile in tiles" :key="tile.to" :to="tile.to" class="y-link-tile group p-5">
          <div class="flex items-start justify-between">
            <span class="y-icon-tile"><component :is="tile.icon" :size="20" :stroke-width="1.5" /></span>
            <ArrowUpRight :size="16" :stroke-width="1.5" class="text-y-faint transition-colors group-hover:text-y-text" />
          </div>
          <h3 class="mt-5 text-[15px] font-semibold text-y-text">{{ tile.title }}</h3>
          <p class="mt-1 text-sm leading-6 text-y-muted">{{ tile.desc }}</p>
        </router-link>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, Copy, Check, ArrowUpRight, Globe, User, Gift, Ticket } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useToast } from '../composables/useToast'
import axios from '../utils/axios'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'

const { t, te } = useI18n()
const toast = useToast()
const authStore = useAuthStore()
const user = computed(() => authStore.user)
const domainCount = ref(0)
const codeCopied = ref(false)

const slotsLeft = computed(() => Math.max(0, (user.value?.domain_quota || 0) - domainCount.value))

const levelLabel = computed(() => {
  const level = user.value?.user_level
  if (!level) return '—'
  return te(`home.levelNames.${level}`) ? t(`home.levelNames.${level}`) : level
})

const tiles = computed(() => [
  { to: '/domains', icon: Globe, title: t('uiUser.dashboard.tiles.domains'), desc: t('uiUser.dashboard.tiles.domainsDesc') },
  { to: '/profile', icon: User, title: t('uiUser.dashboard.tiles.profile'), desc: t('uiUser.dashboard.tiles.profileDesc') },
  { to: '/invitations', icon: Gift, title: t('uiUser.dashboard.tiles.invitations'), desc: t('uiUser.dashboard.tiles.invitationsDesc') },
  { to: '/coupons', icon: Ticket, title: t('uiUser.dashboard.tiles.coupons'), desc: t('uiUser.dashboard.tiles.couponsDesc') },
])

const copyInviteCode = async () => {
  if (!user.value?.invite_code) return
  try {
    await navigator.clipboard.writeText(user.value.invite_code)
    codeCopied.value = true
    toast.success(t('uiUser.common.codeCopied'))
    setTimeout(() => { codeCopied.value = false }, 2000)
  } catch {
    toast.error(t('uiUser.common.copyFailed'))
  }
}

onMounted(async () => {
  if (!authStore.user) {
    await authStore.fetchProfile()
  }
  await fetchDomainCount()
})

const fetchDomainCount = async () => {
  try {
    const response = await axios.get('/api/domains')
    domainCount.value = response.data.domains.length
  } catch (error) {
    console.error('Failed to fetch domain count:', error)
  }
}
</script>
