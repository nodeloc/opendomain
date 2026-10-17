<template>
  <div class="y-container py-14 sm:py-16">
    <div class="max-w-2xl">
      <p class="y-eyebrow">{{ $t('uiPublic.pending.eyebrow') }}</p>
      <h1 class="y-h1 mt-4">{{ $t('pendingDomains.title') }}</h1>
      <p class="y-lead mt-4">{{ $t('pendingDomains.description') }}</p>
    </div>

    <!-- Stats -->
    <div class="mt-10 grid grid-cols-1 gap-px overflow-hidden rounded-lg border border-y-line bg-y-line sm:grid-cols-3">
      <div class="bg-y-surface px-6 py-6">
        <p class="y-label">{{ $t('pendingDomains.total') }}</p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">{{ pagination.total }}</p>
      </div>
      <div class="bg-y-surface px-6 py-6">
        <p class="y-label flex items-center gap-2"><span class="y-dot bg-y-mint" aria-hidden="true"></span>{{ $t('pendingDomains.healthy') }}</p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">{{ healthyCount }}</p>
      </div>
      <div class="bg-y-surface px-6 py-6">
        <p class="y-label flex items-center gap-2"><span class="y-dot bg-y-danger" aria-hidden="true"></span>{{ $t('pendingDomains.unhealthy') }}</p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">{{ unhealthyCount }}</p>
        <p class="mt-1 text-sm text-y-faint">{{ $t('pendingDomains.willBeDeleted') }}</p>
      </div>
    </div>

    <!-- Domains -->
    <div class="y-panel mt-6 overflow-hidden">
      <div class="border-b border-y-line p-3 sm:p-4">
        <label class="relative block">
          <Search :size="16" :stroke-width="1.5" class="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('pendingDomains.searchPlaceholder')"
            class="input input-bordered w-full pl-10 font-mono"
            autocomplete="off"
            autocapitalize="off"
            spellcheck="false"
          />
        </label>
      </div>

      <YLoader v-if="loading" :label="$t('common.loading')" />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
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
            <tr v-for="domain in domains" :key="domain.id">
              <td class="font-mono font-medium">{{ domain.full_domain }}</td>
              <td class="font-mono text-y-muted">{{ domain.root_domain?.domain || 'N/A' }}</td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="statusTone(domain.status)" :label="getStatusText(domain.status)" />
              </td>
              <td class="whitespace-nowrap font-mono text-xs text-y-muted">{{ formatDate(domain.expires_at) }}</td>
              <td class="whitespace-nowrap">
                <div v-if="domain.first_failed_at">
                  <span class="font-mono text-xs text-y-danger-text">{{ formatDate(domain.first_failed_at) }}</span>
                  <p class="mt-1 text-xs text-y-muted">
                    {{ $t('pendingDomains.daysAgo', { days: getDaysSince(domain.first_failed_at) }) }}
                    <span v-if="getDaysSince(domain.first_failed_at) >= 30" class="font-medium text-y-danger-text">
                      · {{ $t('pendingDomains.willBeDeleted') }}
                    </span>
                  </p>
                </div>
                <span v-else class="text-y-faint">—</span>
              </td>
            </tr>
            <tr v-if="domains.length === 0">
              <td colspan="5" class="!p-0">
                <EmptyState :title="$t('pendingDomains.noData')" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pagination -->
    <nav v-if="!loading && pagination.total_pages > 1" class="mt-6 flex items-center justify-center gap-2">
      <button
        type="button"
        class="btn btn-ghost btn-sm px-2 disabled:border-transparent disabled:bg-transparent"
        :disabled="pagination.page === 1"
        :aria-label="$t('common.previous')"
        @click="changePage(pagination.page - 1)"
      >
        <ChevronLeft :size="16" :stroke-width="1.5" />
      </button>
      <span class="px-2 font-mono text-sm text-y-muted">
        {{ $t('pendingDomains.page') }} {{ pagination.page }} / {{ pagination.total_pages }}
      </span>
      <button
        type="button"
        class="btn btn-ghost btn-sm px-2 disabled:border-transparent disabled:bg-transparent"
        :disabled="pagination.page === pagination.total_pages"
        :aria-label="$t('common.next')"
        @click="changePage(pagination.page + 1)"
      >
        <ChevronRight :size="16" :stroke-width="1.5" />
      </button>
    </nav>

    <!-- Info Box -->
    <div class="y-panel mt-6 flex items-start gap-3 px-5 py-4">
      <Info :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-info" />
      <p class="text-sm leading-6 text-y-muted">{{ $t('pendingDomains.infoText') }}</p>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ChevronLeft, ChevronRight, Info, Search } from 'lucide-vue-next'
import axios from '../utils/axios'
import YLoader from '../components/brand/YLoader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
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

let searchTimer = null

const healthyCount = computed(() => domains.value.filter(d => d.status === 'healthy').length)
const unhealthyCount = computed(() => domains.value.filter(d => d.status === 'unhealthy').length)

const fetchDomains = async (page = 1) => {
  loading.value = true
  try {
    const params = {
      page: page,
      per_page: pagination.value.per_page
    }
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    const res = await axios.get('/api/public/pending-domains', { params })
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

watch(searchQuery, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    fetchDomains(1)
  }, 400)
})

const changePage = (page) => {
  if (page >= 1 && page <= pagination.value.total_pages) {
    fetchDomains(page)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const getStatusText = (status) => {
  return t(`pendingDomains.status_${status}`)
}

const statusTone = (status) => {
  const tones = { pending: 'info', healthy: 'ok', unhealthy: 'danger' }
  return tones[status] || 'muted'
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
