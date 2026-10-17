<template>
  <div class="y-container py-14 sm:py-16">
    <div class="max-w-2xl">
      <p class="y-eyebrow">{{ $t('uiPublic.health.eyebrow') }}</p>
      <h1 class="y-h1 mt-4">{{ $t('health.title') }}</h1>
      <p class="y-lead mt-4">{{ $t('health.subtitle') }}</p>
    </div>

    <!-- Statistics -->
    <div class="mt-10 grid grid-cols-2 gap-px overflow-hidden rounded-lg border border-y-line bg-y-line lg:grid-cols-4">
      <div v-for="stat in statCards" :key="stat.key" class="bg-y-surface px-5 py-5 sm:px-6 sm:py-6">
        <p class="y-label flex items-center gap-2">
          <span v-if="stat.dot" class="y-dot" :class="stat.dot" aria-hidden="true"></span>
          {{ stat.label }}
        </p>
        <p class="mt-3 font-mono text-3xl font-medium tracking-tight text-y-text">{{ stat.value }}</p>
        <p class="mt-1 text-sm text-y-faint">{{ stat.desc }}</p>
      </div>
    </div>

    <!-- Health Reports -->
    <div class="y-panel mt-6 overflow-hidden">
      <div class="y-panel-header">
        <h2 class="y-h3">{{ $t('uiPublic.health.reportsTitle') }}</h2>
        <span v-if="!loading && pagination.total_pages > 1" class="font-mono text-xs text-y-faint">
          {{ $t('uiPublic.health.pageOf', { page: pagination.page, total: pagination.total_pages }) }}
        </span>
      </div>

      <YLoader v-if="loading" :label="$t('common.loading')" />

      <EmptyState
        v-else-if="healthReports.length === 0"
        :title="$t('health.noHealthData')"
        :description="$t('health.noHealthDataDesc')"
      />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('health.domain') }}</th>
              <th>{{ $t('health.overallHealth') }}</th>
              <th>{{ $t('health.http') }}</th>
              <th>{{ $t('health.dns') }}</th>
              <th>{{ $t('health.ssl') }}</th>
              <th>{{ $t('health.safeBrowsing') }}</th>
              <th>{{ $t('health.virusTotal') }}</th>
              <th>{{ $t('health.uptime') }}</th>
              <th>{{ $t('health.lastScanned') }}</th>
              <th>{{ $t('health.actionPending') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="report in healthReports" :key="report.domain_id">
              <td class="whitespace-nowrap font-mono font-medium">{{ maskDomain(report.domain_name) }}</td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="healthTone(report.overall_health)" :label="formatHealth(report.overall_health)" />
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="statusTone(report.http_status)" :label="formatHttpStatus(report.http_status)" />
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="statusTone(report.dns_status)" :label="formatDnsStatus(report.dns_status)" />
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="sslTone(report.ssl_status)" :label="formatSslStatus(report.ssl_status)" />
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="safeBrowsingTone(report.safe_browsing_status)" :label="formatSafeBrowsing(report.safe_browsing_status)" />
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="virusTotalTone(report.virustotal_status)" :label="formatVirusTotal(report.virustotal_status)" />
              </td>
              <td>
                <div class="flex items-center gap-2.5">
                  <div
                    class="h-1.5 w-20 overflow-hidden rounded-full bg-y-surface-2"
                    role="progressbar"
                    :aria-valuenow="report.uptime_percentage"
                    aria-valuemin="0"
                    aria-valuemax="100"
                  >
                    <div class="h-full rounded-full" :class="getUptimeClass(report.uptime_percentage)" :style="{ width: `${report.uptime_percentage ?? 0}%` }"></div>
                  </div>
                  <span class="font-mono text-xs text-y-muted">{{ report.uptime_percentage?.toFixed(1) }}%</span>
                </div>
              </td>
              <td class="whitespace-nowrap text-y-muted">{{ formatDate(report.last_scanned_at) }}</td>
              <td>
                <span
                  v-if="getPendingAction(report)"
                  class="whitespace-nowrap"
                  :class="getPendingActionClass(report)"
                  :title="getPendingActionTooltip(report)"
                >
                  {{ getPendingActionText(report) }}
                </span>
                <span v-else class="text-y-faint">—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pagination -->
    <nav v-if="!loading && pagination.total_pages > 1" class="mt-6 flex justify-center" :aria-label="$t('uiPublic.health.reportsTitle')">
      <div class="flex flex-wrap items-center gap-1">
        <button
          type="button"
          class="btn btn-ghost btn-sm px-2 disabled:border-transparent disabled:bg-transparent"
          :disabled="pagination.page === 1"
          :aria-label="$t('common.previous')"
          @click="goToPage(pagination.page - 1)"
        >
          <ChevronLeft :size="16" :stroke-width="1.5" />
        </button>
        <template v-for="(page, index) in visiblePages" :key="index">
          <button
            v-if="typeof page === 'number'"
            type="button"
            class="btn btn-sm min-w-[2.125rem] px-2 font-mono"
            :class="page === pagination.page ? 'btn-outline bg-y-surface-2' : 'btn-ghost text-y-muted'"
            :aria-current="page === pagination.page ? 'page' : undefined"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
          <span v-else class="px-2 font-mono text-sm text-y-faint">…</span>
        </template>
        <button
          type="button"
          class="btn btn-ghost btn-sm px-2 disabled:border-transparent disabled:bg-transparent"
          :disabled="pagination.page === pagination.total_pages"
          :aria-label="$t('common.next')"
          @click="goToPage(pagination.page + 1)"
        >
          <ChevronRight :size="16" :stroke-width="1.5" />
        </button>
      </div>
    </nav>

    <!-- Details Modal -->
    <dialog class="modal" :class="{ 'modal-open': selectedReport }">
      <div class="modal-box w-11/12 max-w-4xl">
        <div v-if="selectedReport">
          <h3 class="break-all font-mono text-lg font-semibold text-y-text">{{ maskDomain(selectedReport.domain_name) }}</h3>

          <div class="mt-5 grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3">
            <div class="y-inset px-4 py-3">
              <p class="y-label">{{ $t('health.overallHealth') }}</p>
              <StatusDot class="mt-2" :tone="healthTone(selectedReport.overall_health)" :label="formatHealth(selectedReport.overall_health)" />
            </div>
            <div class="y-inset px-4 py-3">
              <p class="y-label">{{ $t('health.httpStatus') }}</p>
              <StatusDot class="mt-2" :tone="statusTone(selectedReport.http_status)" :label="formatHttpStatus(selectedReport.http_status)" />
            </div>
            <div class="y-inset px-4 py-3">
              <p class="y-label">{{ $t('health.dnsStatus') }}</p>
              <StatusDot class="mt-2" :tone="statusTone(selectedReport.dns_status)" :label="formatDnsStatus(selectedReport.dns_status)" />
            </div>
            <div class="y-inset px-4 py-3">
              <p class="y-label">{{ $t('health.safeBrowsing') }}</p>
              <StatusDot class="mt-2" :tone="safeBrowsingTone(selectedReport.safe_browsing_status)" :label="formatSafeBrowsing(selectedReport.safe_browsing_status)" />
            </div>
            <div class="y-inset px-4 py-3">
              <p class="y-label">{{ $t('health.virusTotal') }}</p>
              <StatusDot class="mt-2" :tone="virusTotalTone(selectedReport.virustotal_status)" :label="formatVirusTotal(selectedReport.virustotal_status)" />
            </div>
            <div class="y-inset px-4 py-3">
              <p class="y-label">{{ $t('health.uptime') }}</p>
              <p class="mt-1.5 font-mono text-lg font-medium text-y-text">{{ selectedReport.uptime_percentage?.toFixed(1) }}%</p>
            </div>
          </div>

          <YLoader v-if="loadingScans" class="!py-8" />

          <div v-else-if="scans.length > 0" class="mt-6">
            <h4 class="text-sm font-semibold text-y-text">{{ $t('health.recentScans') }}</h4>
            <div class="y-panel mt-3 overflow-x-auto">
              <table class="y-table">
                <thead>
                  <tr>
                    <th>{{ $t('health.scanType') }}</th>
                    <th>{{ $t('health.status') }}</th>
                    <th>{{ $t('health.responseTime') }}</th>
                    <th>{{ $t('health.scanDetails') }}</th>
                    <th>{{ $t('health.scannedAt') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="scan in scans" :key="scan.id">
                    <td><span class="y-tag">{{ scan.scan_type.toUpperCase() }}</span></td>
                    <td class="whitespace-nowrap">
                      <StatusDot :tone="scan.status === 'success' ? 'ok' : 'danger'" :label="scan.status" />
                    </td>
                    <td class="whitespace-nowrap font-mono text-xs text-y-muted">{{ scan.response_time ? scan.response_time + 'ms' : '—' }}</td>
                    <td class="max-w-xs text-sm text-y-muted">
                      <div class="truncate" :title="getScanDetailsTooltip(scan)">
                        <span v-if="scan.http_status_code" class="font-mono">HTTP {{ scan.http_status_code }}</span>
                        <span v-if="scan.ssl_valid !== null" class="ml-2">SSL: {{ scan.ssl_valid ? $t('health.valid') : $t('health.invalid') }}</span>
                        <span v-if="scan.scan_type === 'safebrowsing' && scan.scan_details" class="ml-2">{{ formatScanDetails(scan.scan_details) }}</span>
                        <span v-if="scan.scan_type === 'virustotal' && scan.scan_details" class="ml-2 font-mono">{{ formatScanDetails(scan.scan_details) }}</span>
                        <span v-if="scan.error_message" class="ml-2 text-y-danger-text">{{ truncate(scan.error_message, 50) }}</span>
                      </div>
                    </td>
                    <td class="whitespace-nowrap text-xs text-y-muted">{{ formatDateTime(scan.scanned_at) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="closeDetails">{{ $t('health.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeDetails">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>
<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import axios from '../utils/axios'
import YLoader from '../components/brand/YLoader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'

const { t } = useI18n()

const healthReports = ref([])
const stats = ref({})
const loading = ref(true)
const selectedReport = ref(null)
const scans = ref([])
const loadingScans = ref(false)
const pagination = ref({
  page: 1,
  page_size: 20,
  total: 0,
  total_pages: 0,
})

onMounted(async () => {
  await Promise.all([fetchHealthReports(), fetchStatistics()])
})

const fetchHealthReports = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.page_size,
    }
    const response = await axios.get('/api/public/domain-health', { params })
    healthReports.value = response.data.health_reports || []
    if (response.data.pagination) {
      pagination.value = response.data.pagination
    }
  } catch (error) {
    console.error('Failed to fetch health reports:', error)
  } finally {
    loading.value = false
  }
}

const goToPage = (page) => {
  if (page >= 1 && page <= pagination.value.total_pages) {
    pagination.value.page = page
    fetchHealthReports()
    // Scroll to top of the page
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const visiblePages = computed(() => {
  const current = pagination.value.page
  const total = pagination.value.total_pages
  const pages = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages.filter(p => p !== '...' || pages.indexOf(p) === pages.lastIndexOf(p))
})

const fetchStatistics = async () => {
  try {
    const response = await axios.get('/api/public/health-statistics')
    stats.value = response.data
  } catch (error) {
    console.error('Failed to fetch statistics:', error)
  }
}

const viewDetails = async (report) => {
  selectedReport.value = report
  loadingScans.value = true
  try {
    const response = await axios.get(`/api/public/domain-health/${report.domain_id}/scans`)
    scans.value = response.data.scans
  } catch (error) {
    console.error('Failed to fetch scans:', error)
  } finally {
    loadingScans.value = false
  }
}

const statCards = computed(() => [
  { key: 'total', label: t('health.totalDomains'), value: stats.value.total_domains || 0, desc: t('health.beingMonitored') },
  { key: 'healthy', label: t('health.healthy'), value: stats.value.healthy_domains || 0, desc: t('health.allOperational'), dot: 'bg-y-mint' },
  { key: 'degraded', label: t('health.degraded'), value: stats.value.degraded_domains || 0, desc: t('health.partialIssues'), dot: 'bg-y-warn' },
  { key: 'down', label: t('health.down'), value: stats.value.down_domains || 0, desc: t('health.serviceUnavailable'), dot: 'bg-y-danger' },
])

const closeDetails = () => {
  selectedReport.value = null
  scans.value = []
}

const healthTone = (health) => {
  const tones = { healthy: 'ok', degraded: 'warn', down: 'danger' }
  return tones[health] || 'muted'
}

const statusTone = (status) => {
  if (status === 'online' || status === 'resolved') return 'ok'
  if (status === 'offline' || status === 'failed') return 'danger'
  return 'muted'
}

const sslTone = (status) => {
  if (status === 'valid') return 'ok'
  if (status === 'invalid') return 'danger'
  return 'muted'
}

const safeBrowsingTone = (status) => {
  if (status === 'safe') return 'ok'
  if (status === 'unsafe') return 'danger'
  return 'muted'
}

const virusTotalTone = (status) => {
  if (status === 'clean') return 'ok'
  if (status === 'malicious') return 'danger'
  if (status === 'suspicious') return 'warn'
  return 'muted'
}
const formatSafeBrowsing = (status) => {
  const key = `health.${status || 'unknown'}`
  return t(key) !== key ? t(key) : (status || t('health.unknown'))
}

const formatVirusTotal = (status) => {
  const key = `health.${status || 'unknown'}`
  return t(key) !== key ? t(key) : (status || t('health.unknown'))
}

const formatHttpStatus = (status) => {
  const key = `health.${status || 'unknown'}`
  return t(key) !== key ? t(key) : (status || t('health.unknown'))
}

const formatDnsStatus = (status) => {
  const key = `health.${status || 'unknown'}`
  return t(key) !== key ? t(key) : (status || t('health.unknown'))
}

const formatSslStatus = (status) => {
  const key = `health.${status || 'unknown'}`
  return t(key) !== key ? t(key) : (status || t('health.unknown'))
}

const formatScanDetails = (details) => {
  try {
    const parsed = JSON.parse(details)
    if (parsed.safe !== undefined) {
      return parsed.safe ? t('health.safe') : t('health.unsafe')
    }
    if (parsed.malicious !== undefined) {
      return `M:${parsed.malicious} S:${parsed.suspicious} H:${parsed.harmless}`
    }
    return details
  } catch {
    return details
  }
}

const getUptimeClass = (uptime) => {
  if (uptime >= 99) return 'bg-y-mint'
  if (uptime >= 95) return 'bg-y-warn'
  return 'bg-y-danger'
}
const formatHealth = (health) => {
  const key = `health.${health}`
  return t(key) !== key ? t(key) : health
}

const formatDate = (dateString) => {
  if (!dateString) return t('health.never')
  return new Date(dateString).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const formatDateTime = (dateString) => {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const truncate = (str, length) => {
  return str && str.length > length ? str.substring(0, length) + '...' : str
}

const maskDomain = (domain) => {
  if (!domain) return ''

  // Split domain into subdomain and root domain
  const parts = domain.split('.')
  if (parts.length < 2) return domain

  // Get the subdomain (first part)
  const subdomain = parts[0]

  // If subdomain is too short (1 char), don't mask
  if (subdomain.length <= 1) {
    return domain
  }

  // For 2-character subdomains: show first, mask rest
  if (subdomain.length === 2) {
    const maskedSubdomain = subdomain.substring(0, 1) + '*'
    parts[0] = maskedSubdomain
    return parts.join('.')
  }

  // For 3-character subdomains: show first and last, mask middle
  if (subdomain.length === 3) {
    const maskedSubdomain = subdomain.substring(0, 1) + '*' + subdomain.substring(2)
    parts[0] = maskedSubdomain
    return parts.join('.')
  }

  // For longer subdomains: show first 2 and last 1, mask the middle
  const maskedSubdomain = subdomain.substring(0, 2) + '****' + subdomain.substring(subdomain.length - 1)

  // Reconstruct the domain
  parts[0] = maskedSubdomain
  return parts.join('.')
}

const getScanDetailsTooltip = (scan) => {
  const details = []

  if (scan.http_status_code) {
    details.push(`HTTP ${scan.http_status_code}`)
  }

  if (scan.ssl_valid !== null) {
    details.push(`SSL: ${scan.ssl_valid ? t('health.valid') : t('health.invalid')}`)
  }

  if (scan.scan_details) {
    if (scan.scan_type === 'safebrowsing' || scan.scan_type === 'virustotal') {
      details.push(formatScanDetails(scan.scan_details))
    }
  }

  if (scan.error_message) {
    details.push(`Error: ${scan.error_message}`)
  }

  return details.join(' | ')
}

const getPendingAction = (report) => {
  // Check if already suspended
  if (report.is_suspended) {
    return 'suspended'
  }

  // Auto-suspend if malicious content detected
  if (report.safe_browsing_status === 'unsafe' || report.virustotal_status === 'malicious') {
    return 'auto-suspend'
  }

  // Check if domain is down (HTTP/DNS issues)
  const isDown = (report.http_status === 'offline' || report.dns_status === 'failed')

  if (isDown) {
    if (report.first_failed_at) {
      const firstFailedDate = new Date(report.first_failed_at)
      const now = new Date()
      const daysSinceFailure = Math.floor((now - firstFailedDate) / (1000 * 60 * 60 * 24))

      if (daysSinceFailure >= 30) {
        return 'delete'
      } else if (daysSinceFailure >= 7) {
        return 'suspend'
      } else if (daysSinceFailure >= 0) {
        return 'warning'
      }
    } else {
      // Domain is down but not yet tracked - show monitoring status
      return 'monitoring'
    }
  }

  return null
}

const getPendingActionClass = (report) => {
  const action = getPendingAction(report)
  if (action === 'delete' || action === 'auto-suspend') return 'y-chip-danger'
  if (action === 'suspend' || action === 'warning' || action === 'monitoring') return 'y-chip-warn'
  if (action === 'suspended') return 'y-chip'
  return ''
}
const getPendingActionText = (report) => {
  const action = getPendingAction(report)

  if (action === 'suspended') {
    return t('health.suspended')
  }

  if (action === 'auto-suspend') {
    return t('health.autoSuspended')
  }

  if (action === 'monitoring') {
    return t('health.actionPending')
  }

  if (report.first_failed_at) {
    const firstFailedDate = new Date(report.first_failed_at)
    const now = new Date()
    const daysSinceFailure = Math.floor((now - firstFailedDate) / (1000 * 60 * 60 * 24))

    if (action === 'delete') {
      return t('health.willDelete') + ' ' + t('health.daysRemaining', { days: Math.max(0, 30 - daysSinceFailure) })
    }

    if (action === 'suspend') {
      return t('health.willSuspend') + ' ' + t('health.daysRemaining', { days: Math.max(0, 7 - daysSinceFailure) })
    }

    if (action === 'warning') {
      const daysUntilSuspend = 7 - daysSinceFailure
      return t('health.willSuspend') + ' ' + t('health.daysRemaining', { days: daysUntilSuspend })
    }
  }

  return ''
}

const getPendingActionTooltip = (report) => {
  const action = getPendingAction(report)

  if (action === 'monitoring') {
    return 'Will suspend in 7 days if not resolved'
  }

  return ''
}
</script>
