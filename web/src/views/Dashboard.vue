<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-8">
    <div class="text-center py-8">
      <h1 class="text-4xl font-bold mb-2">{{ $t('dashboard.welcome') }}{{ user?.username }}!</h1>
      <p class="text-xl opacity-70">{{ $t('dashboard.manageText') }}</p>
    </div>

    <!-- Stats -->
    <div class="stats stats-vertical lg:stats-horizontal shadow w-full">
      <div class="stat">
        <div class="stat-title">{{ $t('dashboard.myDomains') }}</div>
        <div class="stat-value">{{ domainCount }}</div>
        <div class="stat-desc">{{ user?.domain_quota - domainCount }} {{ $t('dashboard.slotsRemaining') }}</div>
      </div>

      <div class="stat">
        <div class="stat-title">{{ $t('dashboard.accountLevel') }}</div>
        <div class="stat-value text-primary">{{ user?.user_level }}</div>
        <div class="stat-desc">{{ user?.email_verified ? $t('dashboard.verified') : $t('dashboard.notVerified') }}</div>
      </div>

      <div class="stat">
        <div class="stat-title">{{ $t('dashboard.inviteCode') }}</div>
        <div class="stat-value text-secondary text-2xl">{{ user?.invite_code }}</div>
        <div class="stat-desc">{{ $t('dashboard.shareRewards') }}</div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <router-link to="/domains" class="card bg-base-200 shadow-xl card-hover">
        <div class="card-body items-center text-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14 text-primary mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
          </svg>
          <h3 class="card-title">{{ $t('dashboard.myDomainsCard') }}</h3>
          <p>{{ $t('dashboard.viewManageDomains') }}</p>
        </div>
      </router-link>

      <router-link to="/profile" class="card bg-base-200 shadow-xl card-hover">
        <div class="card-body items-center text-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14 text-secondary mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <h3 class="card-title">{{ $t('dashboard.profileCard') }}</h3>
          <p>{{ $t('dashboard.updateSettings') }}</p>
        </div>
      </router-link>

      <div class="card bg-base-200 shadow-xl card-hover">
        <div class="card-body items-center text-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-14 w-14 text-accent mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
          </svg>
          <h3 class="card-title">{{ $t('dashboard.invitationsCard') }}</h3>
          <p>{{ $t('dashboard.inviteFriends') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import axios from '../utils/axios'

const authStore = useAuthStore()
const user = computed(() => authStore.user)
const domainCount = ref(0)

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
