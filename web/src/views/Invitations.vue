<template>
  <div class="y-page">
    <PageHeader :eyebrow="$t('uiUser.invitations.eyebrow')" :title="$t('uiUser.invitations.title')" :description="$t('uiUser.invitations.description')" />

    <!-- Share -->
    <section class="y-panel overflow-hidden">
      <div class="grid gap-px bg-y-line lg:grid-cols-[1fr_1.6fr]">
        <div class="space-y-3 bg-y-surface px-5 py-5">
          <p class="y-label">{{ $t('uiUser.invitations.codeLabel') }}</p>
          <div class="flex items-center gap-2">
            <div class="y-inset min-w-0 flex-1 truncate px-4 py-2.5 font-mono text-xl font-medium tracking-tight text-y-text">
              {{ stats.invite_code || '—' }}
            </div>
            <button type="button" class="btn btn-outline gap-1.5" :disabled="!stats.invite_code" @click="copyInviteCode">
              <Copy :size="16" :stroke-width="1.5" />
              {{ $t('uiUser.common.copy') }}
            </button>
          </div>
        </div>
        <div class="space-y-3 bg-y-surface px-5 py-5">
          <p class="y-label">{{ $t('uiUser.invitations.linkLabel') }}</p>
          <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
            <input
              :value="inviteUrl || '—'"
              readonly
              class="y-inset h-11 min-w-0 truncate sm:flex-1 px-4 font-mono text-sm text-y-text focus:border-y-mint/70 focus-visible:outline-none"
              :aria-label="$t('uiUser.invitations.linkLabel')"
              @click="$event.target.select()"
            />
            <button type="button" class="btn btn-primary gap-1.5" :disabled="!inviteUrl" @click="copyInviteUrl">
              <Link2 :size="16" :stroke-width="1.5" />
              {{ $t('uiUser.invitations.copyLink') }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats -->
    <section class="y-panel grid grid-cols-2 gap-px overflow-hidden bg-y-line lg:grid-cols-4">
      <div v-for="stat in statCells" :key="stat.label" class="bg-y-surface px-5 py-6">
        <p class="y-label">{{ stat.label }}</p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">{{ stat.value }}</p>
      </div>
    </section>

    <!-- How it works -->
    <section class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <h2 class="y-h3">{{ $t('uiUser.invitations.howTitle') }}</h2>
      </div>
      <ol class="grid gap-px bg-y-line md:grid-cols-3">
        <li v-for="(step, i) in steps" :key="step.title" class="flex gap-4 bg-y-surface px-5 py-5">
          <span class="font-mono text-sm text-y-mint-text">0{{ i + 1 }}</span>
          <div>
            <h3 class="text-[15px] font-semibold text-y-text">{{ step.title }}</h3>
            <p class="mt-1 text-sm leading-6 text-y-muted">{{ step.desc }}</p>
          </div>
        </li>
      </ol>
    </section>

    <!-- History -->
    <section class="y-panel">
      <div class="y-panel-header">
        <h2 class="y-h3">{{ $t('uiUser.invitations.historyTitle') }}</h2>
      </div>

      <YLoader v-if="loadingInvitations" :label="$t('uiUser.common.loading')" />

      <EmptyState v-else-if="invitations.length === 0" :title="$t('uiUser.invitations.emptyTitle')" :description="$t('uiUser.invitations.emptyDesc')" />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('invitation.user') }}</th>
              <th>{{ $t('invitation.joined') }}</th>
              <th>{{ $t('invitation.reward') }}</th>
              <th>{{ $t('invitation.status') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="invitation in invitations" :key="invitation.id">
              <td class="font-medium">{{ invitation.invitee_name }}</td>
              <td class="whitespace-nowrap font-mono text-[13px] text-y-muted">{{ formatDate(invitation.created_at) }}</td>
              <td class="font-mono text-[13px]">{{ invitation.reward_value }}</td>
              <td class="whitespace-nowrap">
                <StatusDot
                  :tone="invitation.reward_given ? 'ok' : 'warn'"
                  :label="invitation.reward_given ? $t('invitation.completed') : $t('invitation.pending')"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Copy, Link2 } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const stats = ref({})
const invitations = ref([])
const loadingInvitations = ref(true)

const inviteUrl = computed(() => {
  if (!stats.value.invite_code) return ''
  return `${window.location.origin}/register?ref=${stats.value.invite_code}`
})

const statCells = computed(() => [
  { label: t('uiUser.invitations.totalInvites'), value: stats.value.total_invites || 0 },
  { label: t('uiUser.invitations.successful'), value: stats.value.successful_invites || 0 },
  { label: t('uiUser.invitations.totalRewards'), value: stats.value.total_rewards || 0 },
  { label: t('uiUser.invitations.currentQuota'), value: stats.value.current_quota || 0 },
])

const steps = computed(() => [
  { title: t('uiUser.invitations.step1Title'), desc: t('uiUser.invitations.step1Desc') },
  { title: t('uiUser.invitations.step2Title'), desc: t('uiUser.invitations.step2Desc') },
  { title: t('uiUser.invitations.step3Title'), desc: t('uiUser.invitations.step3Desc') },
])

onMounted(async () => {
  await Promise.all([fetchStats(), fetchInvitations()])
})

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/invitations/stats')
    stats.value = response.data
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

const fetchInvitations = async () => {
  loadingInvitations.value = true
  try {
    const response = await axios.get('/api/invitations/my')
    invitations.value = response.data.invitations
  } catch (error) {
    console.error('Failed to fetch invitations:', error)
  } finally {
    loadingInvitations.value = false
  }
}

const copyInviteCode = () => {
  if (stats.value.invite_code) {
    navigator.clipboard.writeText(stats.value.invite_code)
    toast.success(t('user.inviteCodeCopied'))
  }
}

const copyInviteUrl = () => {
  if (inviteUrl.value) {
    navigator.clipboard.writeText(inviteUrl.value)
    toast.success(t('invitation.urlCopied'))
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
