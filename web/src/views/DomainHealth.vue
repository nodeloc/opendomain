<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <div>
      <h1 class="text-3xl font-bold">{{ $t('health.title') }}</h1>
      <p class="text-lg opacity-70 mt-2">{{ $t('health.subtitle') }}</p>
    </div>

    <!-- Statistics -->
    <div class="stats stats-vertical lg:stats-horizontal shadow w-full">
      <div class="stat">
        <div class="stat-title">{{ $t('health.totalDomains') }}</div>
        <div class="stat-value">{{ stats.total_domains || 0 }}</div>
        <div class="stat-desc">{{ $t('health.beingMonitored') }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">{{ $t('health.healthy') }}</div>
        <div class="stat-value text-success">{{ stats.healthy_domains || 0 }}</div>
        <div class="stat-desc">{{ $t('health.allOperational') }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">{{ $t('health.degraded') }}</div>
        <div class="stat-value text-warning">{{ stats.degraded_domains || 0 }}</div>
        <div class="stat-desc">{{ $t('health.partialIssues') }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">{{ $t('health.down') }}</div>
        <div class="stat-value text-error">{{ stats.down_domains || 0 }}</div>
        <div class="stat-desc">{{ $t('health.serviceUnavailable') }}</div>
      </div>
    </div>

    <!-- Health Reports -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="healthReports.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <div class="text-6xl mb-4">🏥</div>
        <h3 class="card-title text-2xl mb-2">{{ $t('health.noHealthData') }}</h3>
        <p>{{ $t('health.noHealthDataDesc') }}</p>
      </div>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="table table-zebra w-full">
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
            <th>{{ $t('health.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="report in healthReports" :key="report.domain_id">
            <td class="font-mono font-semibold">{{ maskDomain(report.domain_name) }}</td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getHealthBadgeClass(report.overall_health)">
                {{ formatHealth(report.overall_health) }}
              </div>
            </td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getStatusBadgeClass(report.http_status)">
                {{ formatHttpStatus(report.http_status) }}
              </div>
            </td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getStatusBadgeClass(report.dns_status)">
                {{ formatDnsStatus(report.dns_status) }}
              </div>
            </td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getSSLBadgeClass(report.ssl_status)">
                {{ formatSslStatus(report.ssl_status) }}
              </div>
            </td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getSafeBrowsingBadgeClass(report.safe_browsing_status)">
                {{ formatSafeBrowsing(report.safe_browsing_status) }}
              </div>
            </td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="getVirusTotalBadgeClass(report.virustotal_status)">
                {{ formatVirusTotal(report.virustotal_status) }}
              </div>
            </td>
            <td>
              <div class="flex items-center gap-2">
                <progress
                  class="progress w-20"
                  :class="getUptimeClass(report.uptime_percentage)"
                  :value="report.uptime_percentage"
                  max="100"
                ></progress>
                <span class="text-sm">{{ report.uptime_percentage?.toFixed(1) }}%</span>
              </div>
            </td>
            <td>{{ formatDate(report.last_scanned_at) }}</td>
            <td>
              <div v-if="getPendingAction(report)" class="text-xs">
                <div
                  class="px-2 py-0.5 rounded font-medium inline-flex items-center"
                  :class="getPendingActionClass(report)"
                  :title="getPendingActionTooltip(report)"
                >
                  {{ getPendingActionText(report) }}
                </div>
              </div>
              <span v-else class="text-xs text-gray-500">-</span>
            </td>
            <td>
              <button @click="viewDetails(report)" class="btn btn-sm btn-ghost">
                {{ $t('health.details') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Details Modal -->
    <dialog class="modal" :class="{ 'modal-open': selectedReport }">
      <div class="modal-box w-11/12 max-w-4xl">
        <div v-if="selectedReport">
          <h3 class="font-bold text-2xl mb-4">{{ maskDomain(selectedReport.domain_name) }}</h3>

          <div class="grid grid-cols-2 md:grid-cols-3 gap-4 mb-6">
            <div class="stat bg-base-200 rounded-lg">
              <div class="stat-title">{{ $t('health.overallHealth') }}</div>
              <div class="stat-value text-2xl" :class="getHealthTextClass(selectedReport.overall_health)">
                {{ formatHealth(selectedReport.overall_health) }}
              </div>
            </div>
            <div class="stat bg-base-200 rounded-lg">
              <div class="stat-title">{{ $t('health.httpStatus') }}</div>
              <div class="stat-value text-2xl">{{ selectedReport.http_status }}</div>
            </div>
            <div class="stat bg-base-200 rounded-lg">
              <div class="stat-title">{{ $t('health.dnsStatus') }}</div>
              <div class="stat-value text-2xl">{{ selectedReport.dns_status }}</div>
            </div>
            <div class="stat bg-base-200 rounded-lg">
              <div class="stat-title">{{ $t('health.safeBrowsing') }}</div>
              <div class="stat-value text-2xl" :class="getSafeBrowsingTextClass(selectedReport.safe_browsing_status)">
                {{ formatSafeBrowsing(selectedReport.safe_browsing_status) }}
              </div>
            </div>
            <div class="stat bg-base-200 rounded-lg">
              <div class="stat-title">{{ $t('health.virusTotal') }}</div>
              <div class="stat-value text-2xl" :class="getVirusTotalTextClass(selectedReport.virustotal_status)">
                {{ formatVirusTotal(selectedReport.virustotal_status) }}
              </div>
            </div>
            <div class="stat bg-base-200 rounded-lg">
              <div class="stat-title">{{ $t('health.uptime') }}</div>
              <div class="stat-value text-2xl">{{ selectedReport.uptime_percentage?.toFixed(1) }}%</div>
            </div>
          </div>

          <div v-if="loadingScans" class="flex justify-center py-8">
            <span class="loading loading-spinner loading-lg"></span>
          </div>

          <div v-else-if="scans.length > 0">
            <h4 class="font-bold text-lg mb-3">{{ $t('health.recentScans') }}</h4>
            <div class="overflow-x-auto max-h-96">
              <table class="table table-sm">
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
                    <td>{{ scan.scan_type.toUpperCase() }}</td>
                    <td>
                      <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center" :class="scan.status === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                        {{ scan.status }}
                      </div>
                    </td>
                    <td>{{ scan.response_time ? scan.response_time + 'ms' : '-' }}</td>
                    <td class="text-sm max-w-xs">
                      <div class="truncate" :title="getScanDetailsTooltip(scan)">
                        <span v-if="scan.http_status_code">HTTP {{ scan.http_status_code }}</span>
                        <span v-if="scan.ssl_valid !== null" class="ml-2">SSL: {{ scan.ssl_valid ? $t('health.valid') : $t('health.invalid') }}</span>
                        <span v-if="scan.scan_type === 'safebrowsing' && scan.scan_details" class="ml-2">{{ formatScanDetails(scan.scan_details) }}</span>
                        <span v-if="scan.scan_type === 'virustotal' && scan.scan_details" class="ml-2">{{ formatScanDetails(scan.scan_details) }}</span>
                        <span v-if="scan.error_message" class="text-error ml-2">{{ truncate(scan.error_message, 50) }}</span>
                      </div>
                    </td>
                    <td class="text-xs whitespace-nowrap">{{ formatDateTime(scan.scanned_at) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button @click="closeDetails" class="btn">{{ $t('health.close') }}</button>
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

const { t } = useI18n()

const healthReports = ref([])
const stats = ref({})
const loading = ref(true)
const selectedReport = ref(null)
const scans = ref([])
const loadingScans = ref(false)

onMounted(async () => {
  await Promise.all([fetchHealthReports(), fetchStatistics()])
})

const fetchHealthReports = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/public/domain-health')
    // Only show domains with issues (not healthy)
    healthReports.value = (response.data.health_reports || []).filter(report =>
      report.overall_health !== 'healthy'
    )
  } catch (error) {
    console.error('Failed to fetch health reports:', error)
  } finally {
    loading.value = false
  }
}

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

const closeDetails = () => {
  selectedReport.value = null
  scans.value = []
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
  if (status === 'online' || status === 'resolved') return 'bg-green-100 text-green-800'
  if (status === 'offline' || status === 'failed') return 'bg-red-100 text-red-800'
  return 'bg-gray-100 text-gray-800'
}

const getSSLBadgeClass = (status) => {
  if (status === 'valid') return 'bg-green-100 text-green-800'
  if (status === 'invalid') return 'bg-red-100 text-red-800'
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

const getSafeBrowsingTextClass = (status) => {
  if (status === 'safe') return 'text-success'
  if (status === 'unsafe') return 'text-error'
  return ''
}

const getVirusTotalTextClass = (status) => {
  if (status === 'clean') return 'text-success'
  if (status === 'malicious') return 'text-error'
  if (status === 'suspicious') return 'text-warning'
  return ''
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

const getHealthTextClass = (health) => {
  const classes = {
    healthy: 'text-success',
    degraded: 'text-warning',
    down: 'text-error',
  }
  return classes[health] || ''
}

const getUptimeClass = (uptime) => {
  if (uptime >= 99) return 'progress-success'
  if (uptime >= 95) return 'progress-warning'
  return 'progress-error'
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

  // If subdomain is too short, don't mask
  if (subdomain.length <= 3) {
    return domain
  }

  // Show first 2 characters and last character, mask the middle
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
  if (action === 'delete') return 'bg-red-100 text-red-800'
  if (action === 'suspend' || action === 'auto-suspend') return 'bg-orange-100 text-orange-800'
  if (action === 'warning' || action === 'monitoring') return 'bg-yellow-100 text-yellow-800'
  if (action === 'suspended') return 'bg-gray-100 text-gray-800'
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
