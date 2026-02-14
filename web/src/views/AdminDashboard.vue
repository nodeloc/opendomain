<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-4xl font-bold mb-2">{{ $t('admin.dashboard') }}</h1>
      <p class="text-lg opacity-70">{{ $t('admin.subtitle', { siteName: siteConfigStore.siteName }) }}</p>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="stat bg-base-100 shadow-xl rounded-box">
        <div class="stat-figure text-primary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
        <div class="stat-title">{{ $t('admin.totalUsers') }}</div>
        <div class="stat-value text-primary">{{ stats.total_users || 0 }}</div>
        <div class="stat-desc">{{ $t('admin.platformUsers') }}</div>
      </div>

      <div class="stat bg-base-100 shadow-xl rounded-box">
        <div class="stat-figure text-secondary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
          </svg>
        </div>
        <div class="stat-title">{{ $t('admin.totalDomains') }}</div>
        <div class="stat-value text-secondary">{{ stats.total_domains || 0 }}</div>
        <div class="stat-desc">{{ $t('admin.registeredDomains') }}</div>
      </div>

      <div class="stat bg-base-100 shadow-xl rounded-box">
        <div class="stat-figure text-accent">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="stat-title">{{ $t('admin.totalRevenue') }}</div>
        <div class="stat-value text-accent">{{ formatPrice(stats.total_revenue || 0) }}</div>
        <div class="stat-desc">{{ $t('admin.fromPaidDomains') }}</div>
      </div>
    </div>

    <!-- Admin Functions -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Root Domains Management -->
      <router-link to="/admin/root-domains" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-primary/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.rootDomains') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.rootDomainsDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Domains Management -->
      <router-link to="/admin/domains" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-purple-500/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.domains') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.domainsDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Pending Domains Management -->
      <router-link to="/admin/pending-domains" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-amber-500/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.pendingDomains') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.pendingDomainsDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Scan Status -->
      <router-link to="/admin/scan-status" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-info/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-info" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.scanStatus') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.scanStatusDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Coupons Management -->
      <router-link to="/admin/coupons" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-secondary/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.coupons') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.couponsDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Announcements Management -->
      <router-link to="/admin/announcements" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-accent/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.announcements') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.announcementsDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Pages Management -->
      <router-link to="/admin/pages" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-info/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-info" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.pages') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.pagesDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Users Management -->
      <router-link to="/admin/users" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-warning/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-warning" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.users') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.usersDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- Orders Management -->
      <router-link to="/admin/orders" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-success/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">{{ $t('admin.orders') }}</h2>
              <p class="text-sm opacity-70">{{ $t('admin.ordersDesc') }}</p>
            </div>
          </div>
        </div>
      </router-link>

      <!-- System Settings -->
      <router-link to="/admin/settings" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 border border-base-300">
        <div class="card-body">
          <div class="flex items-center gap-4">
            <div class="p-3 rounded-lg bg-error/10">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div>
              <h2 class="card-title">Settings</h2>
              <p class="text-sm opacity-70">System configuration</p>
            </div>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from '../utils/axios'
import { useSiteConfigStore } from '../stores/siteConfig'
import { useCurrency } from '../composables/useCurrency'

const siteConfigStore = useSiteConfigStore()
const { formatPrice } = useCurrency()

const stats = ref({
  total_users: 0,
  total_domains: 0,
  total_revenue: 0,
})

const fetchStats = async () => {
  try {
    const response = await axios.get('/api/admin/dashboard-stats')
    stats.value = response.data
  } catch (error) {
    console.error('Failed to fetch dashboard stats:', error)
  }
}

onMounted(async () => {
  await fetchStats()
})
</script>
