<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <div>
      <h1 class="text-3xl font-bold">{{ $t('invitation.inviteFriends') }}</h1>
      <p class="text-lg opacity-70 mt-2">{{ $t('invitation.shareEarn') }}</p>
    </div>

    <!-- Invitation Stats Card -->
    <div class="card bg-gradient-to-br from-primary to-secondary text-primary-content shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-2xl">{{ $t('invitation.yourCode') }}</h2>
        <div class="flex items-center gap-4 my-4">
          <div class="font-mono text-4xl font-bold bg-base-100 text-base-content px-6 py-3 rounded-lg">
            {{ stats.invite_code || 'Loading...' }}
          </div>
          <button @click="copyInviteCode" class="btn btn-outline gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            {{ $t('invitation.copy') }}
          </button>
        </div>

        <div v-if="stats.invite_code" class="mt-2">
          <label class="label"><span class="label-text text-primary-content opacity-80">{{ $t('invitation.inviteUrl') }}</span></label>
          <div class="flex items-center gap-2">
            <input
              :value="inviteUrl"
              readonly
              class="input input-bordered bg-base-100 text-base-content flex-1 font-mono text-sm"
              @click="$event.target.select()"
            />
            <button @click="copyInviteUrl" class="btn btn-outline gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
              {{ $t('invitation.copyUrl') }}
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
          <div class="stat bg-base-100 text-base-content rounded-lg">
            <div class="stat-title">{{ $t('invitation.totalInvites') }}</div>
            <div class="stat-value text-primary">{{ stats.total_invites || 0 }}</div>
          </div>
          <div class="stat bg-base-100 text-base-content rounded-lg">
            <div class="stat-title">{{ $t('invitation.successful') }}</div>
            <div class="stat-value text-success">{{ stats.successful_invites || 0 }}</div>
          </div>
          <div class="stat bg-base-100 text-base-content rounded-lg">
            <div class="stat-title">{{ $t('invitation.totalRewards') }}</div>
            <div class="stat-value text-accent">{{ stats.total_rewards || 0 }}</div>
          </div>
          <div class="stat bg-base-100 text-base-content rounded-lg">
            <div class="stat-title">{{ $t('invitation.currentQuota') }}</div>
            <div class="stat-value">{{ stats.current_quota || 0 }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- How it Works -->
    <div class="card bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">{{ $t('invitation.howItWorks') }}</h2>
        <div class="grid md:grid-cols-3 gap-4 mt-4">
          <div class="flex flex-col items-center text-center p-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14 text-primary mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
            </svg>
            <h3 class="font-bold mb-2">{{ $t('invitation.step1Title') }}</h3>
            <p class="text-sm opacity-70">{{ $t('invitation.step1Desc') }}</p>
          </div>
          <div class="flex flex-col items-center text-center p-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14 text-secondary mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <h3 class="font-bold mb-2">{{ $t('invitation.step2Title') }}</h3>
            <p class="text-sm opacity-70">{{ $t('invitation.step2Desc') }}</p>
          </div>
          <div class="flex flex-col items-center text-center p-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14 text-accent mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
            </svg>
            <h3 class="font-bold mb-2">{{ $t('invitation.step3Title') }}</h3>
            <p class="text-sm opacity-70">{{ $t('invitation.step3Desc') }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Invitation History -->
    <div class="card bg-base-200 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">{{ $t('invitation.history') }}</h2>

        <div v-if="loadingInvitations" class="flex justify-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else-if="invitations.length === 0" class="text-center py-8 opacity-70">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto opacity-40 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
          </svg>
          <p>{{ $t('invitation.noInvitationsDesc') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="table table-zebra w-full">
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
                <td class="font-semibold">{{ invitation.invitee_name }}</td>
                <td>{{ formatDate(invitation.created_at) }}</td>
                <td>{{ invitation.reward_value }}</td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium" :class="invitation.reward_given ? 'bg-green-100 text-green-800' : 'bg-amber-100 text-amber-800'">
                    {{ invitation.reward_given ? $t('invitation.completed') : $t('invitation.pending') }}
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const toast = useToast()
const stats = ref({})
const invitations = ref([])
const loadingInvitations = ref(true)

const inviteUrl = computed(() => {
  if (!stats.value.invite_code) return ''
  return `${window.location.origin}/register?ref=${stats.value.invite_code}`
})

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
