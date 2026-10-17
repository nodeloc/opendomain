<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('uiAdmin3.scan.title')"
      :description="$t('uiAdmin3.scan.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-outline btn-sm gap-1.5" :disabled="loadingQuota" @click="fetchQuota">
          <RefreshCw :size="16" :stroke-width="1.5" :class="{ 'animate-spin': loadingQuota }" />
          {{ $t('uiAdmin3.scan.refreshQuota') }}
        </button>
      </template>
    </PageHeader>

    <!-- API quotas -->
    <section class="y-panel">
      <div class="y-panel-header">
        <h2 class="text-base font-semibold text-y-text">{{ $t('uiAdmin3.scan.quotaTitle') }}</h2>
      </div>

      <YLoader v-if="loadingQuota" :size="32" />

      <div v-else class="grid grid-cols-1 md:grid-cols-2">
        <div
          v-for="(meter, index) in quotaMeters"
          :key="meter.key"
          class="space-y-3 px-5 py-5"
          :class="{ 'border-t border-y-line md:border-l md:border-t-0': index > 0 }"
        >
          <div class="flex items-center justify-between gap-3">
            <p class="y-label">{{ meter.name }}</p>
            <span class="font-mono text-xs text-y-faint">{{ meter.pct }}%</span>
          </div>
          <p class="font-mono text-3xl font-medium tracking-tight text-y-text">
            {{ meter.used }}<span class="text-lg text-y-faint"> / {{ meter.limit }}</span>
          </p>
          <div class="h-1 w-full overflow-hidden rounded-full bg-y-surface-2" role="progressbar" :aria-valuenow="meter.used" aria-valuemin="0" :aria-valuemax="meter.limit">
            <div class="h-full rounded-full transition-[width] duration-300" :class="meter.barClass" :style="{ width: meter.pct + '%' }"></div>
          </div>
          <p class="text-sm text-y-muted">
            {{ $t('scanStatus.remaining') }}: <span class="font-mono text-y-text">{{ meter.remaining }}</span>
          </p>
        </div>
      </div>
    </section>

    <!-- Summaries and suspend history -->
    <section class="y-panel">
      <div class="y-panel-header">
        <div role="tablist" class="tabs tabs-boxed w-full sm:w-fit">
          <button
            type="button"
            role="tab"
            class="tab flex-1 sm:flex-none"
            :class="{ 'tab-active': activeTab === 'summaries' }"
            :aria-selected="activeTab === 'summaries'"
            @click="activeTab = 'summaries'"
          >
            {{ $t('uiAdmin3.scan.summariesTab') }}
          </button>
          <button
            type="button"
            role="tab"
            class="tab flex-1 sm:flex-none"
            :class="{ 'tab-active': activeTab === 'history' }"
            :aria-selected="activeTab === 'history'"
            @click="activeTab = 'history'"
          >
            {{ $t('uiAdmin3.scan.historyTab') }}
          </button>
        </div>
      </div>

      <!-- Tab: domain scan summaries -->
      <div v-show="activeTab === 'summaries'">
        <div class="flex flex-col gap-3 border-b border-y-line px-5 py-4 sm:flex-row sm:items-center">
          <label class="input input-bordered flex min-w-0 flex-1 items-center gap-2">
            <Search :size="16" :stroke-width="1.5" class="shrink-0 text-y-faint" />
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="$t('scanStatus.searchPlaceholder')"
              class="min-w-0 grow bg-transparent font-mono text-sm"
              autocomplete="off"
              spellcheck="false"
              @keyup.enter="fetchSummaries"
            />
          </label>
          <select v-model="statusFilter" class="select select-bordered sm:w-44" @change="fetchSummaries">
            <option value="">{{ $t('scanStatus.allStatus') }}</option>
            <option value="healthy">{{ $t('scanStatus.healthy') }}</option>
            <option value="degraded">{{ $t('scanStatus.degraded') }}</option>
            <option value="down">{{ $t('scanStatus.down') }}</option>
          </select>
          <button type="button" class="btn btn-primary" @click="fetchSummaries">
            {{ $t('scanStatus.search') }}
          </button>
        </div>

        <YLoader v-if="loadingSummaries" :label="$t('common.loading')" />

        <EmptyState
          v-else-if="summaries.length === 0"
          :title="$t('scanStatus.noData')"
          :description="$t('uiAdmin3.scan.noDataDesc')"
        />

        <div v-else>
          <div class="overflow-x-auto">
            <table class="y-table">
              <thead>
                <tr>
                  <th>{{ $t('scanStatus.domain') }}</th>
                  <th>{{ $t('scanStatus.overallHealth') }}</th>
                  <th>{{ $t('scanStatus.httpStatus') }}</th>
                  <th>{{ $t('scanStatus.dnsStatus') }}</th>
                  <th>{{ $t('scanStatus.sslStatus') }}</th>
                  <th>{{ $t('scanStatus.safeBrowsing') }}</th>
                  <th>{{ $t('scanStatus.virusTotal') }}</th>
                  <th>{{ $t('scanStatus.uptime') }}</th>
                  <th>{{ $t('scanStatus.lastScanned') }}</th>
                  <th class="text-right">{{ $t('scanStatus.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="summary in summaries" :key="summary.domain_id">
                  <td><span class="y-domain whitespace-nowrap">{{ summary.domain_name }}</span></td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getHealthTone(summary.overall_health)" :label="healthLabel(summary.overall_health)" />
                  </td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getStatusTone(summary.http_status)" :label="summary.http_status || '-'" />
                  </td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getStatusTone(summary.dns_status)" :label="summary.dns_status || '-'" />
                  </td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getStatusTone(summary.ssl_status)" :label="summary.ssl_status || '-'" />
                  </td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getSafeBrowsingTone(summary.safe_browsing_status)" :label="summary.safe_browsing_status || '-'" />
                  </td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getVirusTotalTone(summary.virustotal_status)" :label="summary.virustotal_status || '-'" />
                  </td>
                  <td class="whitespace-nowrap font-mono text-sm">
                    <span :class="getUptimeClass(summary.uptime_percentage)">
                      {{ summary.uptime_percentage ? summary.uptime_percentage.toFixed(2) + '%' : '-' }}
                    </span>
                  </td>
                  <td class="whitespace-nowrap font-mono text-xs text-y-muted">{{ formatDate(summary.last_scanned_at) }}</td>
                  <td class="text-right">
                    <button type="button" class="btn btn-ghost btn-xs gap-1" @click="viewDetails(summary)">
                      {{ $t('scanStatus.viewDetails') }}
                      <ChevronRight :size="14" :stroke-width="1.5" />
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="summaryPagination.total_pages > 1" class="flex items-center justify-between gap-3 border-t border-y-line px-5 py-3">
            <span class="font-mono text-xs text-y-faint">
              {{ $t('uiAdmin3.pageOf', { page: summaryPagination.page, total: summaryPagination.total_pages }) }}
            </span>
            <div class="flex items-center gap-1">
              <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                :disabled="summaryPagination.page === 1"
                :aria-label="$t('uiAdmin3.prevPage')"
                @click="goToSummaryPage(summaryPagination.page - 1)"
              >
                <ChevronLeft :size="16" :stroke-width="1.5" />
              </button>
              <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                :disabled="summaryPagination.page === summaryPagination.total_pages"
                :aria-label="$t('uiAdmin3.nextPage')"
                @click="goToSummaryPage(summaryPagination.page + 1)"
              >
                <ChevronRight :size="16" :stroke-width="1.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tab: suspend history -->
      <div v-show="activeTab === 'history'">
        <div class="flex flex-col gap-3 border-b border-y-line px-5 py-4 sm:flex-row sm:items-center">
          <label class="input input-bordered flex min-w-0 flex-1 items-center gap-2">
            <Search :size="16" :stroke-width="1.5" class="shrink-0 text-y-faint" />
            <input
              v-model="historySearchQuery"
              type="text"
              :placeholder="$t('scanStatus.searchDomain')"
              class="min-w-0 grow bg-transparent font-mono text-sm"
              autocomplete="off"
              spellcheck="false"
              @keyup.enter="fetchSuspendHistory"
            />
          </label>
          <button type="button" class="btn btn-primary" @click="fetchSuspendHistory">
            {{ $t('scanStatus.search') }}
          </button>
        </div>

        <YLoader v-if="loadingHistory" :label="$t('common.loading')" />

        <EmptyState
          v-else-if="suspendHistories.length === 0"
          :title="$t('scanStatus.noHistory')"
          :description="$t('uiAdmin3.scan.noHistoryDesc')"
        />

        <div v-else>
          <div class="overflow-x-auto">
            <table class="y-table">
              <thead>
                <tr>
                  <th>{{ $t('scanStatus.domain') }}</th>
                  <th>{{ $t('scanStatus.reason') }}</th>
                  <th>{{ $t('scanStatus.details') }}</th>
                  <th>{{ $t('scanStatus.suspendedAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="history in suspendHistories" :key="history.id">
                  <td><span class="y-domain whitespace-nowrap">{{ history.domain_name }}</span></td>
                  <td class="whitespace-nowrap">
                    <span class="y-chip-danger">{{ history.reason }}</span>
                  </td>
                  <td>
                    <p class="min-w-[16rem] max-w-md text-sm text-y-muted">{{ history.details }}</p>
                  </td>
                  <td class="whitespace-nowrap font-mono text-xs text-y-muted">{{ formatDateTime(history.created_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="historyPagination.total_pages > 1" class="flex items-center justify-between gap-3 border-t border-y-line px-5 py-3">
            <span class="font-mono text-xs text-y-faint">
              {{ $t('uiAdmin3.pageOf', { page: historyPagination.page, total: historyPagination.total_pages }) }}
            </span>
            <div class="flex items-center gap-1">
              <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                :disabled="historyPagination.page === 1"
                :aria-label="$t('uiAdmin3.prevPage')"
                @click="goToHistoryPage(historyPagination.page - 1)"
              >
                <ChevronLeft :size="16" :stroke-width="1.5" />
              </button>
              <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                :disabled="historyPagination.page === historyPagination.total_pages"
                :aria-label="$t('uiAdmin3.nextPage')"
                @click="goToHistoryPage(historyPagination.page + 1)"
              >
                <ChevronRight :size="16" :stroke-width="1.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Scan records modal -->
    <dialog class="modal" :class="{ 'modal-open': selectedDomain }">
      <div class="modal-box w-11/12 max-w-5xl">
        <div class="flex items-start justify-between gap-4">
          <div class="min-w-0 space-y-1">
            <h3 class="text-lg font-semibold text-y-text">{{ $t('uiAdmin3.scan.recordsTitle') }}</h3>
            <p class="y-domain text-sm">{{ selectedDomain?.domain_name }}</p>
          </div>
          <button type="button" class="btn btn-ghost btn-sm btn-square" :aria-label="$t('scanStatus.close')" @click="closeDetails">
            <X :size="18" :stroke-width="1.5" />
          </button>
        </div>

        <YLoader v-if="loadingRecords" :label="$t('common.loading')" />

        <EmptyState
          v-else-if="scanRecords.length === 0"
          :title="$t('scanStatus.noRecords')"
          :description="$t('uiAdmin3.scan.noRecordsDesc')"
        />

        <div v-else class="y-panel mt-5">
          <div class="overflow-x-auto">
            <table class="y-table">
              <thead>
                <tr>
                  <th>{{ $t('scanStatus.scanType') }}</th>
                  <th>{{ $t('scanStatus.status') }}</th>
                  <th>{{ $t('scanStatus.responseTime') }}</th>
                  <th>{{ $t('scanStatus.details') }}</th>
                  <th>{{ $t('scanStatus.scannedAt') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in scanRecords" :key="record.id">
                  <td><span class="y-tag">{{ record.scan_type }}</span></td>
                  <td class="whitespace-nowrap">
                    <StatusDot :tone="getScanStatusTone(record.status)" :label="record.status" />
                  </td>
                  <td class="whitespace-nowrap font-mono text-sm text-y-muted">{{ record.response_time ? record.response_time + 'ms' : '-' }}</td>
                  <td>
                    <div class="max-w-xs space-y-0.5 text-xs">
                      <p v-if="record.http_status_code" class="font-mono text-y-muted">
                        HTTP <span class="text-y-text">{{ record.http_status_code }}</span>
                      </p>
                      <p v-if="record.ssl_valid !== null" class="inline-flex items-center gap-1 font-mono text-y-muted">
                        SSL
                        <span :class="record.ssl_valid ? 'text-y-mint-text' : 'text-y-danger-text'" class="inline-flex items-center gap-0.5">
                          <Check v-if="record.ssl_valid" :size="12" :stroke-width="2" />
                          <X v-else :size="12" :stroke-width="2" />
                          {{ record.ssl_valid ? $t('uiAdmin3.scan.sslValid') : $t('uiAdmin3.scan.sslInvalid') }}
                        </span>
                      </p>
                      <p v-if="record.error_message" class="truncate text-y-danger-text" :title="record.error_message">
                        {{ record.error_message }}
                      </p>
                      <p v-if="record.scan_details" class="truncate font-mono text-y-muted" :title="record.scan_details">
                        {{ formatScanDetails(record.scan_details) }}
                      </p>
                    </div>
                  </td>
                  <td class="whitespace-nowrap font-mono text-xs text-y-muted">{{ formatDateTime(record.scanned_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="recordPagination.total_pages > 1" class="flex items-center justify-between gap-3 border-t border-y-line px-5 py-3">
            <span class="font-mono text-xs text-y-faint">
              {{ $t('uiAdmin3.pageOf', { page: recordPagination.page, total: recordPagination.total_pages }) }}
            </span>
            <div class="flex items-center gap-1">
              <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                :disabled="recordPagination.page === 1"
                :aria-label="$t('uiAdmin3.prevPage')"
                @click="goToRecordPage(recordPagination.page - 1)"
              >
                <ChevronLeft :size="16" :stroke-width="1.5" />
              </button>
              <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                :disabled="recordPagination.page === recordPagination.total_pages"
                :aria-label="$t('uiAdmin3.nextPage')"
                @click="goToRecordPage(recordPagination.page + 1)"
              >
                <ChevronRight :size="16" :stroke-width="1.5" />
              </button>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="closeDetails">{{ $t('scanStatus.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeDetails">
        <button>{{ $t('scanStatus.close') }}</button>
      </form>
    </dialog>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RefreshCw, Search, ChevronLeft, ChevronRight, Check, X } from 'lucide-vue-next'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()

// Tab state
const activeTab = ref('summaries')

const quota = ref({})
const loadingQuota = ref(false)

const summaries = ref([])
const loadingSummaries = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const summaryPagination = ref({
  page: 1,
  page_size: 50,
  total: 0,
  total_pages: 0,
})

const selectedDomain = ref(null)
const scanRecords = ref([])
const loadingRecords = ref(false)
const recordPagination = ref({
  page: 1,
  page_size: 50,
  total: 0,
  total_pages: 0,
})

const suspendHistories = ref([])
const loadingHistory = ref(false)
const historySearchQuery = ref('')
const historyPagination = ref({
  page: 1,
  page_size: 20,
  total: 0,
  total_pages: 0,
})

onMounted(() => {
  fetchQuota()
  fetchSuspendHistory()
  fetchSummaries()
})

const fetchQuota = async () => {
  loadingQuota.value = true
  try {
    const response = await axios.get('/api/admin/api-quota')
    quota.value = response.data
  } catch (error) {
    console.error('Failed to fetch quota:', error)
    toast.error(t('scanStatus.fetchQuotaFailed'))
  } finally {
    loadingQuota.value = false
  }
}

const fetchSummaries = async () => {
  loadingSummaries.value = true
  try {
    const params = {
      page: summaryPagination.value.page,
      page_size: summaryPagination.value.page_size,
    }
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    const response = await axios.get('/api/admin/scan-summaries', { params })
    summaries.value = response.data.summaries || []
    if (response.data.pagination) {
      summaryPagination.value = response.data.pagination
    }
  } catch (error) {
    console.error('Failed to fetch summaries:', error)
    toast.error(t('scanStatus.fetchSummariesFailed'))
  } finally {
    loadingSummaries.value = false
  }
}

const goToSummaryPage = (page) => {
  summaryPagination.value.page = page
  fetchSummaries()
}

const fetchSuspendHistory = async () => {
  loadingHistory.value = true
  try {
    const params = {
      page: historyPagination.value.page,
      page_size: historyPagination.value.page_size,
    }
    if (historySearchQuery.value) {
      params.search = historySearchQuery.value
    }
    const response = await axios.get('/api/admin/suspend-history', { params })
    suspendHistories.value = response.data.histories || []
    if (response.data.pagination) {
      historyPagination.value = response.data.pagination
    }
  } catch (error) {
    console.error('Failed to fetch suspend history:', error)
    toast.error(t('scanStatus.fetchHistoryFailed'))
  } finally {
    loadingHistory.value = false
  }
}

const goToHistoryPage = (page) => {
  historyPagination.value.page = page
  fetchSuspendHistory()
}

const viewDetails = async (summary) => {
  selectedDomain.value = summary
  recordPagination.value.page = 1
  await fetchScanRecords()
}

const fetchScanRecords = async () => {
  if (!selectedDomain.value) return
  
  loadingRecords.value = true
  try {
    const params = {
      domain_id: selectedDomain.value.domain_id,
      page: recordPagination.value.page,
      page_size: recordPagination.value.page_size,
    }
    const response = await axios.get('/api/admin/scan-records', { params })
    scanRecords.value = response.data.scans || []
    if (response.data.pagination) {
      recordPagination.value = response.data.pagination
    }
  } catch (error) {
    console.error('Failed to fetch scan records:', error)
    toast.error(t('scanStatus.fetchRecordsFailed'))
  } finally {
    loadingRecords.value = false
  }
}

const goToRecordPage = (page) => {
  recordPagination.value.page = page
  fetchScanRecords()
}

const closeDetails = () => {
  selectedDomain.value = null
  scanRecords.value = []
}

const quotaMeter = (key, name, data, fallbackLimit) => {
  const used = data?.used || 0
  const limit = data?.limit || fallbackLimit
  const remaining = data?.remaining ?? Math.max(limit - used, 0)
  const pct = limit > 0 ? Math.min(100, Math.round((used / limit) * 100)) : 0
  const barClass = pct >= 100 ? 'bg-y-danger' : pct >= 85 ? 'bg-y-warn' : 'bg-y-mint'
  return { key, name, used, limit, remaining, pct, barClass }
}

const quotaMeters = computed(() => [
  quotaMeter('google_safe_browsing', 'Google Safe Browsing', quota.value.google_safe_browsing, 10000),
  quotaMeter('virustotal', 'VirusTotal', quota.value.virustotal, 500),
])

const healthLabel = (health) => {
  if (['healthy', 'degraded', 'down'].includes(health)) return t(`scanStatus.${health}`)
  return health || '-'
}

const getHealthTone = (health) => {
  const tones = {
    healthy: 'ok',
    degraded: 'warn',
    down: 'danger',
  }
  return tones[health] || 'muted'
}

const getStatusTone = (status) => {
  if (status === 'success' || status === 'online' || status === 'resolved' || status === 'valid' || status === 'safe' || status === 'clean') {
    return 'ok'
  }
  if (status === 'failed' || status === 'offline' || status === 'invalid' || status === 'unsafe' || status === 'malicious') {
    return 'danger'
  }
  return 'muted'
}

const getScanStatusTone = (status) => {
  if (status === 'success') return 'ok'
  if (status === 'threat_detected') return 'danger'
  if (status === 'failed') return 'warn'
  if (status === 'quota_exceeded') return 'info'
  return 'muted'
}

const getSafeBrowsingTone = (status) => {
  if (status === 'safe') return 'ok'
  if (status === 'unsafe') return 'danger'
  return 'muted'
}

const getVirusTotalTone = (status) => {
  if (status === 'clean') return 'ok'
  if (status === 'malicious') return 'danger'
  if (status === 'suspicious') return 'warn'
  return 'muted'
}

const getUptimeClass = (uptime) => {
  if (!uptime) return 'text-y-faint'
  if (uptime >= 99) return 'text-y-mint-text font-medium'
  if (uptime >= 95) return 'text-y-text'
  if (uptime >= 90) return 'text-y-warn'
  return 'text-y-danger-text font-medium'
}
const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

const formatScanDetails = (details) => {
  if (!details) return '-'
  try {
    const parsed = JSON.parse(details)
    if (parsed.safe !== undefined) {
      return parsed.safe ? 'Safe' : 'Unsafe'
    }
    if (parsed.malicious !== undefined) {
      return `M:${parsed.malicious} S:${parsed.suspicious} H:${parsed.harmless}`
    }
    return details.substring(0, 50)
  } catch {
    return details.substring(0, 50)
  }
}
</script>
