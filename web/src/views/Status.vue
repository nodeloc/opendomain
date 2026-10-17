<template>
  <div class="mx-auto flex w-full max-w-[960px] flex-col gap-8 px-5 py-12 sm:px-6 sm:py-16">
    <!-- Page title -->
    <div class="flex flex-wrap items-center gap-3">
      <h1 class="y-h1">{{ $t('status.pageTitle') }}</h1>
      <span class="y-site-badge">{{ $t('status.badge') }}</span>
      <button type="button" class="btn btn-outline btn-sm ml-auto gap-1.5" :disabled="loading" @click="loadAll">
        <RefreshCw :size="14" :stroke-width="1.5" :class="{ 'animate-spin': loading }" aria-hidden="true" />
        {{ $t('status.refresh') }}
      </button>
    </div>

    <!-- Overall hero card -->
    <section
      class="relative flex flex-wrap items-center gap-x-6 gap-y-4 overflow-hidden rounded-xl border p-7"
      :class="heroClass"
    >
      <div class="y-glow pointer-events-none absolute -top-40 right-[-120px] h-[420px] w-[420px]" aria-hidden="true"></div>
      <span
        class="relative flex h-11 w-11 shrink-0 items-center justify-center rounded-full"
        :class="overall.tone === 'ok' ? 'bg-y-mint' : overall.tone === 'muted' ? 'bg-y-surface-2' : overall.tone === 'warn' ? 'bg-y-warn' : 'bg-y-danger'"
        :style="overall.tone === 'ok' ? 'box-shadow: 0 0 24px rgb(var(--y-mint) / 0.6);' : ''"
        aria-hidden="true"
      >
        <Check v-if="overall.tone === 'ok'" :size="22" :stroke-width="2.6" class="text-y-on-mint" />
        <TriangleAlert v-else-if="overall.tone === 'warn'" :size="20" :stroke-width="2.4" class="text-y-bg" />
        <OctagonX v-else-if="overall.tone === 'danger'" :size="20" :stroke-width="2.4" class="text-y-bg" />
        <LoaderCircle v-else :size="20" :stroke-width="2" class="animate-spin text-y-muted" />
      </span>
      <div class="relative min-w-0 flex-1">
        <h2 class="m-0 text-3xl font-extrabold leading-9 tracking-[-0.03em] text-y-text">{{ overall.title }}</h2>
        <div class="mt-1 text-sm text-y-muted">
          {{ $t('status.overallSub') }}<template v-if="lastCheckedAt"> · {{ $t('status.lastCheck', { ago: lastCheckAgo }) }}</template>
        </div>
      </div>
    </section>

    <!-- Components -->
    <section class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <h2 class="m-0 text-lg font-bold text-y-text">{{ $t('status.componentsTitle') }}</h2>
      </div>
      <div
        v-for="c in componentRows"
        :key="c.key"
        class="flex flex-wrap items-center gap-x-4 gap-y-1 border-b border-y-line px-5 py-4 last:border-b-0"
      >
        <span class="font-semibold text-y-text">{{ c.label }}</span>
        <span v-if="c.host" class="font-mono text-[13px] text-y-faint">{{ c.host }}</span>
        <span class="ml-auto flex items-center gap-4">
          <span v-if="c.latencyMs !== null" class="font-mono text-[13px] text-y-muted">{{ $t('status.ms', { n: c.latencyMs }) }}</span>
          <StatusDot :tone="c.tone" :label="c.statusLabel" />
        </span>
      </div>
      <p v-if="apiError" class="m-0 border-t border-y-line px-5 py-3 text-sm text-y-warn">
        {{ $t('status.apiUnreachable') }}
      </p>
    </section>

    <!-- Nameservers -->
    <section v-if="nameservers.length || dohError" class="flex flex-col gap-4">
      <div>
        <h2 class="m-0 text-lg font-bold text-y-text">{{ $t('status.nameserversTitle') }}</h2>
        <p class="m-0 mt-1 text-sm text-y-muted">{{ $t('status.nameserversDesc') }}</p>
      </div>
      <p v-if="dohError" class="m-0 text-sm text-y-warn">{{ $t('status.dohUnavailable') }}</p>
      <div class="grid gap-4 sm:grid-cols-2">
        <div v-for="ns in nameservers" :key="ns.name" class="y-card flex flex-col gap-2.5 p-5">
          <div class="flex items-center justify-between gap-3">
            <span class="y-domain text-[17px]">{{ ns.name }}</span>
            <span
              v-if="ns.done"
              :class="ns.ok ? 'y-dot-live' : 'y-dot bg-y-danger'"
              role="img"
              :aria-label="ns.ok ? $t('status.resolves') : $t('status.notResolving')"
            ></span>
            <span v-else class="loading loading-spinner loading-xs"></span>
          </div>
          <div class="text-sm text-y-muted">PowerDNS</div>
          <div class="flex flex-col gap-1 font-mono text-[13px]">
            <span class="text-y-faint">
              {{ $t('status.ipv4') }}:
              <span :class="ns.v4.length ? 'text-y-text' : 'text-y-danger-text'">{{ ns.done ? (ns.v4.join(', ') || $t('status.none')) : '…' }}</span>
            </span>
            <span class="text-y-faint">
              {{ $t('status.ipv6') }}:
              <span :class="ns.v6.length ? 'text-y-text' : 'text-y-faint'">{{ ns.done ? (ns.v6.join(', ') || $t('status.none')) : '…' }}</span>
            </span>
          </div>
          <div v-if="ns.done && ns.ms !== null" class="text-[13px] text-y-faint">
            {{ $t('status.lookup') }} <span class="font-mono text-y-muted">{{ $t('status.ms', { n: ns.ms }) }}</span>
          </div>
          <div v-else-if="ns.done && ns.error" class="text-[13px] text-y-danger-text">{{ $t('status.lookupFailed') }}</div>
        </div>
      </div>
    </section>

    <!-- Root domains -->
    <section v-if="rootDomains.length" class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <div>
          <h2 class="m-0 text-lg font-bold text-y-text">{{ $t('status.rootsTitle') }}</h2>
          <p class="m-0 mt-1 text-sm text-y-muted">{{ $t('status.rootsDesc') }}</p>
        </div>
      </div>
      <div
        v-for="r in rootDomains"
        :key="r.domain"
        class="flex flex-wrap items-center gap-x-4 gap-y-2 border-b border-y-line px-5 py-4 last:border-b-0"
      >
        <span class="y-domain">{{ r.domain }}</span>
        <span class="ml-auto flex flex-wrap items-center justify-end gap-2">
          <template v-if="r.done">
            <template v-if="r.ns.length">
              <span v-for="n in r.ns" :key="n" class="y-chip font-mono">{{ n }}</span>
            </template>
            <span v-else class="text-sm text-y-warn">{{ $t('status.noNs') }}</span>
            <StatusDot :tone="r.ns.length ? 'ok' : 'warn'" :label="r.ns.length ? $t('status.resolves') : $t('status.notResolving')" />
          </template>
          <span v-else class="loading loading-spinner loading-xs"></span>
        </span>
      </div>
    </section>

    <!-- Past incidents -->
    <section class="flex flex-col gap-4">
      <h2 class="m-0 text-[28px] font-extrabold uppercase leading-[34px] tracking-[-0.03em] text-y-text">{{ $t('status.incidentsTitle') }}</h2>
      <div class="flex flex-col border-t border-y-line">
        <div class="flex flex-wrap gap-x-6 gap-y-2 border-b border-y-line py-4">
          <span class="w-[120px] text-sm text-y-faint">{{ $t('status.today') }}</span>
          <span class="min-w-0 flex-1 text-y-muted" style="flex-basis: 300px;">{{ $t('status.noIncidents') }}</span>
        </div>
      </div>
      <i18n-t keypath="status.maintenance" tag="div" class="text-sm text-y-faint">
        <template #link>
          <router-link to="/announcements" class="text-y-mint-text hover:text-y-mint-hover">{{ $t('status.maintenanceLink') }}</router-link>
        </template>
      </i18n-t>
      <div class="text-sm text-y-faint">{{ $t('status.footnote') }}</div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check, TriangleAlert, OctagonX, LoaderCircle, RefreshCw } from 'lucide-vue-next'
import axios from '../utils/axios'
import StatusDot from '../components/ui/StatusDot.vue'

const { t } = useI18n()

const loading = ref(false)
const apiError = ref(false)
const dohError = ref(false)
const overallStatus = ref('checking') // operational | degraded | outage | checking
const apiComponents = ref([])
const websiteMs = ref(null)
const websiteOk = ref(null)
const lastCheckedAt = ref(null)
const now = ref(Date.now())
const nameservers = ref([])
const rootDomains = ref([])

const toneFor = (status) => ({ operational: 'ok', degraded: 'warn', outage: 'danger' }[status] || 'muted')
const labelFor = (status) => ({
  operational: t('status.sOperational'),
  degraded: t('status.sDegraded'),
  outage: t('status.sOutage'),
  checking: t('status.sChecking'),
}[status] || t('status.sUnknown'))

const overall = computed(() => ({
  tone: overallStatus.value === 'checking' ? 'muted' : toneFor(overallStatus.value),
  title: {
    operational: t('status.overallOperational'),
    degraded: t('status.overallDegraded'),
    outage: t('status.overallOutage'),
    checking: t('status.overallChecking'),
  }[overallStatus.value] || t('status.sUnknown'),
}))

const heroClass = computed(() => ({
  ok: 'y-mint-soft',
  warn: 'y-warn-soft',
  danger: 'y-danger-soft',
  muted: 'bg-y-surface border-y-line',
}[overall.value.tone]))

const componentRows = computed(() => {
  const rows = [
    {
      key: 'website',
      label: t('status.cWebsite'),
      host: window.location.host,
      latencyMs: websiteMs.value,
      tone: websiteOk.value === null ? 'muted' : websiteOk.value ? 'ok' : 'danger',
      statusLabel: websiteOk.value === null ? t('status.sChecking') : websiteOk.value ? t('status.sOperational') : t('status.sOutage'),
    },
  ]
  const labels = { api: t('status.cApi'), database: t('status.cDatabase'), cache: t('status.cCache'), dns_api: t('status.cDnsApi') }
  if (apiComponents.value.length) {
    for (const c of apiComponents.value) {
      rows.push({
        key: c.key,
        label: labels[c.key] || c.key,
        host: '',
        latencyMs: c.latency_ms ?? null,
        tone: toneFor(c.status),
        statusLabel: labelFor(c.status),
      })
    }
  } else if (apiError.value) {
    for (const key of ['api', 'database', 'cache', 'dns_api']) {
      rows.push({ key, label: labels[key], host: '', latencyMs: null, tone: 'danger', statusLabel: t('status.sOutage') })
    }
  } else {
    for (const key of ['api', 'database', 'cache', 'dns_api']) {
      rows.push({ key, label: labels[key], host: '', latencyMs: null, tone: 'muted', statusLabel: t('status.sChecking') })
    }
  }
  return rows
})

const lastCheckAgo = computed(() => {
  if (!lastCheckedAt.value) return ''
  const s = Math.max(0, Math.round((now.value - lastCheckedAt.value) / 1000))
  if (s < 5) return t('status.justNow')
  if (s < 60) return t('status.secondsAgo', { n: s })
  return t('status.minutesAgo', { n: Math.floor(s / 60) })
})

const DOH = 'https://cloudflare-dns.com/dns-query'

const dohQuery = async (name, type) => {
  const started = performance.now()
  const res = await fetch(`${DOH}?name=${encodeURIComponent(name)}&type=${type}`, {
    headers: { accept: 'application/dns-json' },
  })
  const ms = Math.round(performance.now() - started)
  if (!res.ok) throw new Error(`DoH ${res.status}`)
  const json = await res.json()
  const typeNum = { A: 1, NS: 2, AAAA: 28 }[type]
  const answers = (json.Answer || []).filter((a) => a.type === typeNum).map((a) => String(a.data).replace(/\.$/, ''))
  return { answers, ms }
}

const checkNameservers = async (hosts) => {
  nameservers.value = hosts.map((name) => ({ name, v4: [], v6: [], ms: null, ok: false, done: false, error: false }))
  await Promise.all(
    nameservers.value.map(async (ns) => {
      try {
        const [a, aaaa] = await Promise.all([dohQuery(ns.name, 'A'), dohQuery(ns.name, 'AAAA').catch(() => ({ answers: [], ms: null }))])
        ns.v4 = a.answers
        ns.v6 = aaaa.answers
        ns.ms = a.ms
        ns.ok = a.answers.length > 0
      } catch {
        ns.error = true
      } finally {
        ns.done = true
      }
    })
  )
}

const checkRootDomains = async (domains) => {
  rootDomains.value = domains.map((domain) => ({ domain, ns: [], done: false }))
  await Promise.all(
    rootDomains.value.map(async (r) => {
      try {
        const { answers } = await dohQuery(r.domain, 'NS')
        r.ns = answers
      } catch {
        r.ns = []
      } finally {
        r.done = true
      }
    })
  )
}

const checkWebsite = async () => {
  const started = performance.now()
  try {
    const res = await fetch('/', { cache: 'no-store' })
    websiteMs.value = Math.round(performance.now() - started)
    websiteOk.value = res.ok
  } catch {
    websiteMs.value = null
    websiteOk.value = false
  }
}

const loadAll = async () => {
  loading.value = true
  const websitePromise = checkWebsite()
  try {
    const response = await axios.get('/api/public/status', { timeout: 8000 })
    const data = response.data || {}
    apiError.value = false
    apiComponents.value = Array.isArray(data.components) ? data.components : []
    overallStatus.value = ['operational', 'degraded', 'outage'].includes(data.status) ? data.status : 'checking'
    lastCheckedAt.value = data.checked_at ? new Date(data.checked_at).getTime() : Date.now()
    now.value = Date.now()
    const nsHosts = Array.isArray(data.nameservers) ? data.nameservers : []
    const roots = Array.isArray(data.root_domains) ? data.root_domains : []
    await Promise.all([checkNameservers(nsHosts), checkRootDomains(roots)])
  } catch (err) {
    console.error('Failed to fetch status:', err)
    apiError.value = true
    apiComponents.value = []
    overallStatus.value = 'outage'
    lastCheckedAt.value = Date.now()
    now.value = Date.now()
  } finally {
    await websitePromise
    loading.value = false
  }
}

let refreshTimer = null
let clockTimer = null

onMounted(() => {
  loadAll()
  refreshTimer = setInterval(loadAll, 60000)
  clockTimer = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onBeforeUnmount(() => {
  clearInterval(refreshTimer)
  clearInterval(clockTimer)
})
</script>
