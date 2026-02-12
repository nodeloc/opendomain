<template>
  <div class="container mx-auto px-4 py-8">
    <div class="mb-6">
      <h1 class="text-3xl font-bold">{{ $t('pendingDomains.title') }}</h1>
      <p class="text-sm opacity-70 mt-2">{{ $t('pendingDomains.description') }}</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('pendingDomains.total') }}</div>
        <div class="stat-value text-primary">{{ pagination.total }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('pendingDomains.healthy') }}</div>
        <div class="stat-value text-success">{{ healthyCount }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('pendingDomains.unhealthy') }}</div>
        <div class="stat-value text-error">{{ unhealthyCount }}</div>
        <div class="stat-desc">{{ $t('pendingDomains.willBeDeleted') }}</div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <input
        v-model="searchQuery"
        type="text"
        :placeholder="$t('pendingDomains.searchPlaceholder')"
        class="input input-bordered w-full"
      />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Domains Table -->
    <div v-else class="overflow-x-auto bg-base-100 rounded-lg shadow">
      <table class="table">
        <thead>
          <tr>
            <th>{{ $t('pendingDomains.domain') }}</th>
            <th>{{ $t('pendingDomains.rootDomain') }}</th>
            <th>{{ $t('pendingDomains.status') }}</th>
            <th>{{ $t('pendingDomains.expiresAt') }}</th>
            <th>{{ $t('pendingDomains.firstFailedAt') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="domain in filteredDomains" :key="domain.id">
            <td class="font-mono">{{ domain.full_domain }}</td>
            <td>{{ domain.root_domain?.domain || 'N/A' }}</td>
            <td>
              <span
                class="badge"
                :class="{
                  'badge-primary': domain.status === 'pending',
                  'badge-success': domain.status === 'healthy',
                  'badge-error': domain.status === 'unhealthy'
                }"
              >
                {{ getStatusText(domain.status) }}
              </span>
            </td>
            <td>{{ formatDate(domain.expires_at) }}</td>
            <td>
              <span v-if="domain.first_failed_at" class="text-error">
                {{ formatDate(domain.first_failed_at) }}
                <br />
                <span class="text-xs">
                  {{ $t('pendingDomains.daysAgo', { days: getDaysSince(domain.first_failed_at) }) }}
                  <span v-if="getDaysSince(domain.first_failed_at) >= 30" class="font-bold">
                    - {{ $t('pendingDomains.willBeDeleted') }}
                  </span>
                </span>
              </span>
              <span v-else class="opacity-50">-</span>
            </td>
          </tr>
          <tr v-if="filteredDomains.length === 0">
            <td colspan="5" class="text-center py-8 opacity-50">
              {{ $t('pendingDomains.noData') }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="!loading && pagination.total_pages > 1" class="flex justify-center mt-6">
      <div class="join">
        <button
          class="join-item btn"
          :disabled="pagination.page === 1"
          @click="changePage(pagination.page - 1)"
        >
          «
        </button>
        <button class="join-item btn">
          {{ $t('pendingDomains.page') }} {{ pagination.page }} / {{ pagination.total_pages }}
        </button>
        <button
          class="join-item btn"
          :disabled="pagination.page === pagination.total_pages"
          @click="changePage(pagination.page + 1)"
        >
          »
        </button>
      </div>
    </div>

    <!-- Info Box -->
    <div class="alert alert-info mt-6">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <span>{{ $t('pendingDomains.infoText') }}</span>
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

const domains = ref([])
const searchQuery = ref('')
const loading = ref(true)
const pagination = ref({
  page: 1,
  per_page: 50,
  total: 0,
  total_pages: 0
})

const healthyCount = computed(() => domains.value.filter(d => d.status === 'healthy').length)
const unhealthyCount = computed(() => domains.value.filter(d => d.status === 'unhealthy').length)

const filteredDomains = computed(() => {
  if (!searchQuery.value) return domains.value
  const query = searchQuery.value.toLowerCase()
  return domains.value.filter(d =>
    d.full_domain.toLowerCase().includes(query) ||
    d.subdomain.toLowerCase().includes(query)
  )
})

const fetchDomains = async (page = 1) => {
  loading.value = true
  try {
    const res = await axios.get('/api/public/pending-domains', {
      params: {
        page: page,
        per_page: pagination.value.per_page
      }
    })
    domains.value = res.data.pending_domains || []
    if (res.data.pagination) {
      pagination.value = res.data.pagination
    }
  } catch (error) {
    toast.error(t('pendingDomains.fetchFailed') + ': ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const changePage = (page) => {
  if (page >= 1 && page <= pagination.value.total_pages) {
    fetchDomains(page)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const getStatusText = (status) => {
  return t(`pendingDomains.status_${status}`)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getDaysSince = (date) => {
  if (!date) return 0
  const now = new Date()
  const then = new Date(date)
  return Math.floor((now - then) / (1000 * 60 * 60 * 24))
}

onMounted(() => {
  fetchDomains()
})
</script>
