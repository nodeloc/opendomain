<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold mb-6">{{ $t('scanStatus.title') }}</h1>

    <!-- API 配额状态 -->
    <div class="card bg-base-100 shadow-xl mb-6">
      <div class="card-body">
        <h2 class="card-title">{{ $t('scanStatus.quotaStatus') }}</h2>
        
        <div v-if="loadingQuota" class="flex justify-center py-4">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
          <!-- Google Safe Browsing -->
          <div class="stat bg-base-200 rounded-lg">
            <div class="stat-title">Google Safe Browsing</div>
            <div class="stat-value text-2xl">{{ quota.google_safe_browsing?.used || 0 }} / {{ quota.google_safe_browsing?.limit || 10000 }}</div>
            <div class="stat-desc">
              {{ $t('scanStatus.remaining') }}: {{ quota.google_safe_browsing?.remaining || 10000 }}
            </div>
            <progress 
              class="progress progress-primary w-full mt-2" 
              :value="quota.google_safe_browsing?.used || 0" 
              :max="quota.google_safe_browsing?.limit || 10000"
            ></progress>
          </div>

          <!-- VirusTotal -->
          <div class="stat bg-base-200 rounded-lg">
            <div class="stat-title">VirusTotal</div>
            <div class="stat-value text-2xl">{{ quota.virustotal?.used || 0 }} / {{ quota.virustotal?.limit || 500 }}</div>
            <div class="stat-desc">
              {{ $t('scanStatus.remaining') }}: {{ quota.virustotal?.remaining || 500 }}
            </div>
            <progress 
              class="progress progress-secondary w-full mt-2" 
              :value="quota.virustotal?.used || 0" 
              :max="quota.virustotal?.limit || 500"
            ></progress>
          </div>
        </div>
        
        <div class="flex justify-end mt-4">
          <button @click="fetchQuota" class="btn btn-sm btn-ghost">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {{ $t('scanStatus.refresh') }}
          </button>
        </div>
      </div>
    </div>

    <!-- 域名扫描摘要 -->
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title">{{ $t('scanStatus.scanSummaries') }}</h2>
        
        <!-- 搜索和筛选 -->
        <div class="flex gap-4 mb-4">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('scanStatus.searchPlaceholder')"
            class="input input-bordered flex-1"
            @keyup.enter="fetchSummaries"
          />
          <select v-model="statusFilter" class="select select-bordered" @change="fetchSummaries">
            <option value="">{{ $t('scanStatus.allStatus') }}</option>
            <option value="healthy">{{ $t('scanStatus.healthy') }}</option>
            <option value="degraded">{{ $t('scanStatus.degraded') }}</option>
            <option value="down">{{ $t('scanStatus.down') }}</option>
          </select>
          <button @click="fetchSummaries" class="btn btn-primary">
            {{ $t('scanStatus.search') }}
          </button>
        </div>

        <div v-if="loadingSummaries" class="flex justify-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else-if="summaries.length === 0" class="text-center py-8 opacity-60">
          {{ $t('scanStatus.noData') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="table table-zebra w-full">
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
                <th>{{ $t('scanStatus.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="summary in summaries" :key="summary.domain_id">
                <td class="font-mono">{{ summary.domain_name }}</td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getHealthBadgeClass(summary.overall_health)">
                    {{ summary.overall_health }}
                  </div>
                </td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getStatusBadgeClass(summary.http_status)">
                    {{ summary.http_status }}
                  </div>
                </td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getStatusBadgeClass(summary.dns_status)">
                    {{ summary.dns_status }}
                  </div>
                </td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getStatusBadgeClass(summary.ssl_status)">
                    {{ summary.ssl_status }}
                  </div>
                </td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getSafeBrowsingBadgeClass(summary.safe_browsing_status)">
                    {{ summary.safe_browsing_status }}
                  </div>
                </td>
                <td>
                  <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getVirusTotalBadgeClass(summary.virustotal_status)">
                    {{ summary.virustotal_status }}
                  </div>
                </td>
                <td>
                  <span :class="getUptimeClass(summary.uptime_percentage)">
                    {{ summary.uptime_percentage ? summary.uptime_percentage.toFixed(2) + '%' : '-' }}
                  </span>
                </td>
                <td>{{ formatDate(summary.last_scanned_at) }}</td>
                <td>
                  <button @click="viewDetails(summary)" class="btn btn-sm btn-ghost">
                    {{ $t('scanStatus.viewDetails') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- 分页 -->
          <div v-if="summaryPagination.total_pages > 1" class="flex justify-center mt-6">
            <div class="join">
              <button
                class="join-item btn btn-sm"
                :disabled="summaryPagination.page === 1"
                @click="goToSummaryPage(summaryPagination.page - 1)"
              >
                «
              </button>
              <button class="join-item btn btn-sm">
                {{ $t('scanStatus.page') }} {{ summaryPagination.page }} / {{ summaryPagination.total_pages }}
              </button>
              <button
                class="join-item btn btn-sm"
                :disabled="summaryPagination.page === summaryPagination.total_pages"
                @click="goToSummaryPage(summaryPagination.page + 1)"
              >
                »
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 扫描记录详情 Modal -->
    <dialog class="modal" :class="{ 'modal-open': selectedDomain }">
      <div class="modal-box w-11/12 max-w-5xl">
        <h3 class="font-bold text-2xl mb-4">
          {{ $t('scanStatus.scanRecords') }}: {{ selectedDomain?.domain_name }}
        </h3>

        <div v-if="loadingRecords" class="flex justify-center py-8">
          <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else-if="scanRecords.length === 0" class="text-center py-8 opacity-60">
          {{ $t('scanStatus.noRecords') }}
        </div>

        <div v-else>
          <div class="overflow-x-auto">
            <table class="table table-sm">
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
                  <td>
                    <span class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center border border-gray-300">{{ record.scan_type.toUpperCase() }}</span>
                  </td>
                  <td>
                    <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getScanStatusClass(record.status)">
                      {{ record.status }}
                    </div>
                  </td>
                  <td>{{ record.response_time ? record.response_time + 'ms' : '-' }}</td>
                  <td>
                    <div class="max-w-xs">
                      <div v-if="record.http_status_code" class="text-xs">
                        HTTP: {{ record.http_status_code }}
                      </div>
                      <div v-if="record.ssl_valid !== null" class="text-xs">
                        SSL: {{ record.ssl_valid ? '✓ Valid' : '✗ Invalid' }}
                      </div>
                      <div v-if="record.error_message" class="text-xs text-error truncate" :title="record.error_message">
                        {{ record.error_message }}
                      </div>
                      <div v-if="record.scan_details" class="text-xs truncate" :title="record.scan_details">
                        {{ formatScanDetails(record.scan_details) }}
                      </div>
                    </div>
                  </td>
                  <td class="text-xs">{{ formatDateTime(record.scanned_at) }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 记录分页 -->
          <div v-if="recordPagination.total_pages > 1" class="flex justify-center mt-4">
            <div class="join">
              <button
                class="join-item btn btn-sm"
                :disabled="recordPagination.page === 1"
                @click="goToRecordPage(recordPagination.page - 1)"
              >
                «
              </button>
              <button class="join-item btn btn-sm">
                {{ recordPagination.page }} / {{ recordPagination.total_pages }}
              </button>
              <button
                class="join-item btn btn-sm"
                :disabled="recordPagination.page === recordPagination.total_pages"
                @click="goToRecordPage(recordPagination.page + 1)"
              >
                »
              </button>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button @click="closeDetails" class="btn">{{ $t('scanStatus.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeDetails">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const toast = useToast()

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

onMounted(() => {
  fetchQuota()
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

const getHealthBadgeClass = (health) => {
  const classes = {
    healthy: 'bg-green-100 text-green-800',
    degraded: 'bg-yellow-100 text-yellow-800',
    down: 'bg-red-100 text-red-800',
  }
  return classes[health] || 'bg-gray-100 text-gray-800'
}

const getStatusBadgeClass = (status) => {
  if (status === 'success' || status === 'online' || status === 'resolved' || status === 'valid' || status === 'safe' || status === 'clean') {
    return 'bg-green-100 text-green-800'
  }
  if (status === 'failed' || status === 'offline' || status === 'invalid' || status === 'unsafe' || status === 'malicious') {
    return 'bg-red-100 text-red-800'
  }
  return 'bg-gray-100 text-gray-800'
}

const getScanStatusClass = (status) => {
  if (status === 'success') return 'bg-green-100 text-green-800'
  if (status === 'threat_detected') return 'bg-red-100 text-red-800'
  if (status === 'failed') return 'bg-yellow-100 text-yellow-800'
  if (status === 'quota_exceeded') return 'bg-blue-100 text-blue-800'
  return 'bg-gray-100 text-gray-800'
}

const getSafeBrowsingBadgeClass = (status) => {
  if (status === 'safe') return 'bg-green-100 text-green-800'
  if (status === 'unsafe') return 'bg-red-100 text-red-800'
  return 'bg-gray-100 text-gray-800'
}

const getVirusTotalBadgeClass = (status) => {
  if (status === 'clean') return 'bg-green-100 text-green-800'
  if (status === 'malicious') return 'bg-red-100 text-red-800'
  if (status === 'suspicious') return 'bg-yellow-100 text-yellow-800'
  return 'bg-gray-100 text-gray-800'
}

const getUptimeClass = (uptime) => {
  if (!uptime) return 'text-gray-400'
  if (uptime >= 99) return 'text-green-600 font-semibold'
  if (uptime >= 95) return 'text-green-500'
  if (uptime >= 90) return 'text-yellow-600'
  return 'text-red-600 font-semibold'
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
