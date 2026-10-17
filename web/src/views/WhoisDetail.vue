<template>
  <div class="y-container py-14 sm:py-16">
    <div class="mx-auto max-w-3xl">
      <!-- Loading -->
      <YLoader v-if="loading" class="min-h-[300px]" :label="$t('common.loading')" />

      <!-- Not Found -->
      <div v-else-if="notFound" class="y-panel flex flex-col items-center px-6 py-14 text-center">
        <span class="y-icon-tile"><SearchX :size="20" :stroke-width="1.5" /></span>
        <h1 class="y-h3 mt-5">{{ $t('whois.notFound') }}</h1>
        <i18n-t keypath="whois.notFoundDesc" tag="p" class="mt-2 max-w-md text-sm leading-6 text-y-muted">
          <template #domain>
            <span class="y-domain">{{ domain }}</span>
          </template>
        </i18n-t>
        <router-link to="/" class="btn btn-primary mt-6">{{ $t('whois.searchAnother') }}</router-link>
      </div>

      <!-- WHOIS Info -->
      <div v-else-if="whois">
        <!-- Header -->
        <nav class="flex items-center gap-2 font-mono text-eyebrow uppercase" aria-label="Breadcrumb">
          <router-link to="/" class="text-y-faint transition-colors hover:text-y-text">{{ $t('nav.home') }}</router-link>
          <span class="text-y-faint" aria-hidden="true">/</span>
          <span class="text-y-mint-text">{{ $t('whois.breadcrumb') }}</span>
        </nav>
        <h1 class="mt-4 break-all font-mono text-[clamp(1.75rem,5vw,2.75rem)] font-medium leading-tight tracking-[-0.03em] text-y-text">
          {{ whois.full_domain }}
        </h1>
        <p class="y-lead mt-3">{{ $t('uiPublic.whois.registrationDesc') }}</p>

        <!-- Deletion Warning Banner -->
        <div
          v-if="whois.days_until_deletion !== null && whois.days_until_deletion !== undefined"
          class="alert mt-8 flex items-start text-left"
          :class="whois.days_until_deletion <= 3 ? 'alert-error' : 'alert-warning'"
        >
          <TriangleAlert :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0" />
          <div>
            <p class="font-semibold">{{ $t('whois.pendingDeletion') }}</p>
            <i18n-t keypath="whois.pendingDeletionDesc" tag="p" class="mt-0.5 text-sm">
              <template #days>
                <strong class="font-mono">{{ whois.days_until_deletion }} {{ whois.days_until_deletion !== 1 ? $t('whois.days') : $t('whois.day') }}</strong>
              </template>
            </i18n-t>
          </div>
        </div>

        <!-- Suspension Warning Banner (not yet suspended) -->
        <div
          v-else-if="whois.days_until_suspension !== null && whois.days_until_suspension !== undefined"
          class="alert alert-warning mt-8 flex items-start text-left"
        >
          <TriangleAlert :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0" />
          <div>
            <p class="font-semibold">{{ $t('whois.healthCheckFailed') }}</p>
            <i18n-t keypath="whois.healthCheckFailedDesc" tag="p" class="mt-0.5 text-sm">
              <template #days>
                <strong class="font-mono">{{ whois.days_until_suspension }} {{ whois.days_until_suspension !== 1 ? $t('whois.days') : $t('whois.day') }}</strong>
              </template>
            </i18n-t>
          </div>
        </div>

        <!-- Registration info -->
        <section class="y-panel mt-8">
          <div class="y-panel-header">
            <h2 class="text-sm font-semibold text-y-text">{{ $t('whois.registrationInfo') }}</h2>
          </div>
          <dl class="divide-y divide-y-line">
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.domain') }}</dt>
              <dd class="y-domain text-sm">{{ whois.full_domain }}</dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.status') }}</dt>
              <dd><StatusDot :tone="whoisStatusTone(whois.status)" :label="whois.status" /></dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.registered') }}</dt>
              <dd class="text-sm text-y-text">{{ formatDate(whois.registered_at) }}</dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.expires') }}</dt>
              <dd class="flex flex-wrap items-center gap-2 text-sm" :class="isExpiringSoon ? 'font-medium text-y-warn' : 'text-y-text'">
                {{ formatDate(whois.expires_at) }}
                <span v-if="isExpiringSoon" class="y-chip-warn">{{ $t('whois.expiringSoon') }}</span>
              </dd>
            </div>
            <div v-if="whois.suspended_at" class="whois-row">
              <dt class="y-label">{{ $t('whois.suspended') }}</dt>
              <dd class="text-sm text-y-warn">{{ formatDate(whois.suspended_at) }}</dd>
            </div>
            <div v-if="whois.suspend_reason" class="whois-row">
              <dt class="y-label">{{ $t('whois.suspendReason') }}</dt>
              <dd class="text-sm text-y-warn">{{ whois.suspend_reason }}</dd>
            </div>
            <div v-if="whois.first_failed_at" class="whois-row">
              <dt class="y-label">{{ $t('whois.firstFailed') }}</dt>
              <dd class="text-sm text-y-text">{{ formatDate(whois.first_failed_at) }}</dd>
            </div>
            <div v-if="whois.nameservers && whois.nameservers.length" class="whois-row">
              <dt class="y-label">{{ $t('whois.nameservers') }}</dt>
              <dd class="space-y-1">
                <div v-for="ns in whois.nameservers" :key="ns" class="break-all font-mono text-sm text-y-text">{{ ns }}</div>
              </dd>
            </div>
          </dl>
        </section>

        <!-- Scan Info -->
        <section v-if="whois.scan" class="y-panel mt-6">
          <div class="y-panel-header">
            <h2 class="text-sm font-semibold text-y-text">{{ $t('whois.healthMonitoring') }}</h2>
            <span class="text-xs text-y-faint">
              {{ $t('whois.lastScanned') }}: {{ whois.scan.last_scanned_at ? formatDate(whois.scan.last_scanned_at) : $t('whois.never') }}
            </span>
          </div>

          <!-- Overall Health + Uptime -->
          <div class="flex flex-col gap-4 border-b border-y-line px-5 py-4 sm:flex-row sm:items-center sm:gap-6">
            <StatusDot :tone="healthTone(whois.scan.overall_health)" :label="whois.scan.overall_health || 'unknown'" class="shrink-0 font-medium" />
            <div class="flex-1">
              <div class="mb-1.5 flex justify-between text-xs">
                <span class="y-label">{{ $t('whois.uptime') }}</span>
                <span class="font-mono text-y-muted">{{ whois.scan.uptime_percentage?.toFixed(1) }}%</span>
              </div>
              <div
                class="h-1.5 w-full overflow-hidden rounded-full bg-y-surface-2"
                role="progressbar"
                :aria-valuenow="whois.scan.uptime_percentage ?? 0"
                aria-valuemin="0"
                aria-valuemax="100"
              >
                <div
                  class="h-full rounded-full transition-all"
                  :class="uptimeBarClass(whois.scan.uptime_percentage)"
                  :style="{ width: `${whois.scan.uptime_percentage ?? 0}%` }"
                ></div>
              </div>
            </div>
          </div>

          <dl class="divide-y divide-y-line">
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.http') }}</dt>
              <dd><StatusDot :tone="statusTone(whois.scan.http_status, ['online'], ['offline'])" :label="whois.scan.http_status || 'unknown'" /></dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.dns') }}</dt>
              <dd><StatusDot :tone="statusTone(whois.scan.dns_status, ['resolved', 'success'], ['failed'])" :label="whois.scan.dns_status || 'unknown'" /></dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.ssl') }}</dt>
              <dd><StatusDot :tone="statusTone(whois.scan.ssl_status, ['valid'], ['invalid'])" :label="whois.scan.ssl_status || 'unknown'" /></dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.safeBrowsing') }}</dt>
              <dd><StatusDot :tone="statusTone(whois.scan.safe_browsing_status, ['safe'], ['unsafe', 'threat_detected'])" :label="whois.scan.safe_browsing_status || 'unknown'" /></dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.virusTotal') }}</dt>
              <dd><StatusDot :tone="statusTone(whois.scan.virustotal_status, ['clean'], ['malicious', 'suspicious'])" :label="whois.scan.virustotal_status || 'unknown'" /></dd>
            </div>
            <div class="whois-row">
              <dt class="y-label">{{ $t('whois.totalScans') }}</dt>
              <dd class="font-mono text-sm text-y-text">{{ $t('whois.totalScansValue', { total: whois.scan.total_scans, successful: whois.scan.successful_scans }) }}</dd>
            </div>
          </dl>
        </section>

        <!-- Footer -->
        <div class="mt-8 flex justify-center">
          <router-link to="/" class="btn btn-outline btn-sm gap-1.5">
            <Search :size="16" :stroke-width="1.5" />
            {{ $t('whois.searchAnother') }}
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Search, SearchX, TriangleAlert } from 'lucide-vue-next'
import axios from '../utils/axios'
import YLoader from '../components/brand/YLoader.vue'
import StatusDot from '../components/ui/StatusDot.vue'

const route = useRoute()
const domain = ref(route.params.domain)
const whois = ref(null)
const loading = ref(true)
const notFound = ref(false)

const isExpiringSoon = computed(() => {
  if (!whois.value?.expires_at) return false
  const diff = new Date(whois.value.expires_at) - new Date()
  return diff > 0 && diff < 30 * 24 * 60 * 60 * 1000
})

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleString()
}

const whoisStatusTone = (status) => {
  if (status === 'active') return 'ok'
  if (status === 'suspended') return 'warn'
  if (status === 'expired') return 'danger'
  return 'muted'
}

const healthTone = (health) => {
  if (health === 'healthy') return 'ok'
  if (health === 'degraded') return 'warn'
  if (health === 'down') return 'danger'
  return 'muted'
}

const statusTone = (value, okValues, badValues) => {
  if (okValues.includes(value)) return 'ok'
  if (badValues.includes(value)) return 'danger'
  return 'muted'
}

const uptimeBarClass = (pct) => {
  if (pct >= 90) return 'bg-y-mint'
  if (pct >= 50) return 'bg-y-warn'
  return 'bg-y-danger'
}

const fetchWhois = async () => {
  loading.value = true
  notFound.value = false
  try {
    const response = await axios.get(`/api/public/whois/${domain.value}`)
    whois.value = response.data
  } catch (error) {
    if (error.response?.status === 404) {
      notFound.value = true
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchWhois()
})
</script>

<style scoped>
.whois-row {
  @apply grid gap-1 px-5 py-3.5 sm:grid-cols-[11rem_1fr] sm:items-center sm:gap-4;
}
</style>
