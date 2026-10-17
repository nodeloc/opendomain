<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('abuseAdmin.title')"
      :description="$t('abuseAdmin.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-outline gap-2" :disabled="loading" @click="fetchReports">
          <span v-if="loading" class="loading loading-spinner loading-sm"></span>
          <RefreshCw v-else :size="16" :stroke-width="1.5" />
          {{ $t('abuseAdmin.refresh') }}
        </button>
      </template>
    </PageHeader>

    <!-- Status filter tabs with counts -->
    <div class="y-seg w-fit max-w-full overflow-x-auto" role="group" :aria-label="$t('abuseAdmin.columns.status')">
      <button
        v-for="tab in statusTabs"
        :key="tab.key"
        type="button"
        :aria-pressed="statusFilter === tab.key"
        @click="filterByStatus(tab.key)"
      >
        {{ tab.label }}
        <span class="ml-1.5 font-mono text-xs opacity-80">{{ tab.count }}</span>
      </button>
    </div>

    <section class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <label class="relative block w-full sm:max-w-md">
          <span class="sr-only">{{ $t('common.search') }}</span>
          <Search :size="16" :stroke-width="1.5" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('abuseAdmin.searchPlaceholder')"
            class="input input-bordered h-10 min-h-0 w-full pl-9 pr-10 text-sm"
            @input="debouncedSearch"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="btn btn-ghost btn-xs btn-square absolute right-2 top-1/2 -translate-y-1/2 text-y-faint"
            :aria-label="$t('uiAdmin1.clearSearch')"
            @click="clearSearch"
          >
            <X :size="14" :stroke-width="1.5" />
          </button>
        </label>
        <span class="font-mono text-xs text-y-faint">{{ $t('abuseAdmin.totalCount', { count: total }) }}</span>
      </div>

      <YLoader v-if="loading" :label="$t('common.loading')" />

      <EmptyState
        v-else-if="reports.length === 0"
        :title="searchQuery || statusFilter ? $t('abuseAdmin.noResults') : $t('abuseAdmin.noReports')"
        :description="searchQuery || statusFilter ? '' : $t('abuseAdmin.noReportsDesc')"
      />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('abuseAdmin.columns.created') }}</th>
              <th>{{ $t('abuseAdmin.columns.subdomain') }}</th>
              <th>{{ $t('abuseAdmin.columns.category') }}</th>
              <th>{{ $t('abuseAdmin.columns.reporter') }}</th>
              <th>{{ $t('abuseAdmin.columns.status') }}</th>
              <th class="text-right">{{ $t('abuseAdmin.columns.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="report in reports" :key="report.id">
              <td class="whitespace-nowrap text-sm text-y-muted">{{ formatDateTime(report.created_at) }}</td>
              <td>
                <button type="button" class="y-domain block text-left" @click="openDetail(report)">
                  {{ report.subdomain }}
                </button>
                <span
                  v-if="report.domain_id"
                  class="mt-1 inline-block"
                  :class="domainStatusChipClass(report.domain_status)"
                >{{ domainStatusLabel(report.domain_status) }}</span>
                <span v-else class="y-tag mt-1">{{ $t('abuseAdmin.notRegistered') }}</span>
              </td>
              <td>
                <span class="y-tag">{{ categoryLabel(report.category) }}</span>
              </td>
              <td>
                <span v-if="report.reporter_email" class="break-all text-sm text-y-muted">{{ report.reporter_email }}</span>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td class="whitespace-nowrap">
                <span :class="statusChipClass(report.status)">{{ statusLabel(report.status) }}</span>
              </td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('abuseAdmin.actions.view')"
                    :aria-label="$t('abuseAdmin.actions.view')"
                    @click="openDetail(report)"
                  >
                    <Eye :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('common.delete')"
                    :aria-label="$t('common.delete')"
                    @click="confirmDelete(report)"
                  >
                    <Trash2 :size="16" :stroke-width="1.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div
        v-if="!loading && totalPages > 1"
        class="flex flex-col items-center justify-between gap-3 border-t border-y-line px-5 py-3 sm:flex-row"
      >
        <p class="font-mono text-xs text-y-faint">{{ $t('uiAdmin1.pageOf', { page, total: totalPages }) }}</p>
        <div class="join">
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="page === 1"
            :aria-label="$t('common.previous')"
            @click="goToPage(page - 1)"
          >
            <ChevronLeft :size="16" :stroke-width="1.5" />
          </button>
          <span class="join-item flex h-[2.125rem] items-center px-3 font-mono text-xs text-y-muted">{{ page }}</span>
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="page >= totalPages"
            :aria-label="$t('common.next')"
            @click="goToPage(page + 1)"
          >
            <ChevronRight :size="16" :stroke-width="1.5" />
          </button>
        </div>
      </div>
    </section>

    <!-- Detail modal -->
    <dialog class="modal" :class="{ 'modal-open': showDetail }">
      <div class="modal-box w-11/12 max-w-2xl">
        <div v-if="selected" class="space-y-5">
          <div class="flex flex-wrap items-center gap-3">
            <h3 class="text-lg font-semibold text-y-text">{{ $t('abuseAdmin.detail.title') }}</h3>
            <span class="y-domain">{{ selected.subdomain }}</span>
            <span :class="statusChipClass(selected.status)">{{ statusLabel(selected.status) }}</span>
          </div>

          <dl class="grid grid-cols-1 gap-x-6 gap-y-4 sm:grid-cols-2">
            <div>
              <dt class="y-label">{{ $t('abuseAdmin.detail.id') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selected.id }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('abuseAdmin.detail.reportedAt') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ formatDateTime(selected.created_at) }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('abuseAdmin.detail.category') }}</dt>
              <dd class="mt-1"><span class="y-tag">{{ categoryLabel(selected.category) }}</span></dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('abuseAdmin.detail.reporterEmail') }}</dt>
              <dd class="mt-1 break-all text-sm text-y-text">{{ selected.reporter_email || '—' }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('abuseAdmin.detail.reporterIp') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selected.reporter_ip || '—' }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('abuseAdmin.detail.matchedDomain') }}</dt>
              <dd class="mt-1 text-sm text-y-text">
                <template v-if="selected.domain_id">
                  <span class="font-mono">#{{ selected.domain_id }}</span>
                  <span class="ml-2" :class="domainStatusChipClass(selected.domain_status)">
                    {{ domainStatusLabel(selected.domain_status) }}
                  </span>
                </template>
                <span v-else class="text-y-faint">{{ $t('abuseAdmin.notRegistered') }}</span>
              </dd>
            </div>
            <div v-if="selected.resolved_at">
              <dt class="y-label">{{ $t('abuseAdmin.detail.resolvedAt') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ formatDateTime(selected.resolved_at) }}</dd>
            </div>
          </dl>

          <div>
            <p class="y-label">{{ $t('abuseAdmin.detail.url') }}</p>
            <div v-if="selected.url" class="y-inset mt-1.5 flex items-start gap-2 px-3 py-2">
              <span class="min-w-0 flex-1 break-all font-mono text-xs text-y-text">{{ selected.url }}</span>
              <button
                type="button"
                class="btn btn-ghost btn-xs btn-square shrink-0"
                :title="$t('abuseAdmin.detail.copyUrl')"
                :aria-label="$t('abuseAdmin.detail.copyUrl')"
                @click="copyUrl(selected.url)"
              >
                <Copy :size="14" :stroke-width="1.5" />
              </button>
            </div>
            <p v-if="selected.url" class="mt-1.5 text-xs text-y-warn">{{ $t('abuseAdmin.detail.urlWarning') }}</p>
            <p v-else class="mt-1.5 text-sm text-y-faint">{{ $t('abuseAdmin.detail.noUrl') }}</p>
          </div>

          <div>
            <p class="y-label">{{ $t('abuseAdmin.detail.details') }}</p>
            <div class="y-inset mt-1.5 max-h-48 overflow-y-auto whitespace-pre-wrap break-words px-3 py-2 text-sm leading-6 text-y-text">{{ selected.details }}</div>
          </div>

          <div class="y-field">
            <label class="y-field-label" for="abuse-admin-note">{{ $t('abuseAdmin.detail.adminNote') }}</label>
            <textarea
              id="abuse-admin-note"
              v-model="draftNote"
              class="textarea textarea-bordered h-24 w-full"
              :placeholder="$t('abuseAdmin.detail.adminNotePlaceholder')"
            ></textarea>
          </div>

          <div class="y-field">
            <span class="y-field-label">{{ $t('abuseAdmin.detail.setStatus') }}</span>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="opt in statusOptions"
                :key="opt.value"
                type="button"
                class="btn btn-sm"
                :class="draftStatus === opt.value ? 'btn-primary' : 'btn-outline'"
                :aria-pressed="draftStatus === opt.value"
                @click="draftStatus = opt.value"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>
        </div>

        <div class="modal-action justify-between">
          <button type="button" class="btn btn-error btn-outline" @click="confirmDelete(selected)">
            {{ $t('common.delete') }}
          </button>
          <div class="flex gap-2">
            <button type="button" class="btn btn-ghost" @click="closeDetail">{{ $t('common.cancel') }}</button>
            <button type="button" class="btn btn-primary" :disabled="saving || !dirty" @click="saveDetail">
              <span v-if="saving" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ $t('abuseAdmin.actions.save') }}</span>
            </button>
          </div>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeDetail">
        <button>{{ $t('common.close') }}</button>
      </form>
    </dialog>

    <!-- Delete confirmation modal -->
    <dialog class="modal" :class="{ 'modal-open': showDelete }">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('abuseAdmin.confirmDeleteTitle') }}</h3>
        <p class="mt-3 text-sm leading-6 text-y-muted">
          {{ $t('abuseAdmin.confirmDeleteMessage', { subdomain: deleteTarget?.subdomain }) }}
        </p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDelete = false">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-error" :disabled="deleting" @click="handleDelete">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            <span v-else>{{ $t('common.delete') }}</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showDelete = false">
        <button>{{ $t('common.close') }}</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { formatDateTime } from '../utils/time'
import { RefreshCw, Search, X, Eye, Trash2, ChevronLeft, ChevronRight, Copy } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(false)
const reports = ref([])
const counts = ref({ open: 0, reviewing: 0, resolved: 0, dismissed: 0 })
const total = ref(0)
const page = ref(1)
const pageSize = 20
const statusFilter = ref('')
const searchQuery = ref('')

const showDetail = ref(false)
const selected = ref(null)
const draftStatus = ref('open')
const draftNote = ref('')
const saving = ref(false)

const showDelete = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

const statusKeys = ['open', 'reviewing', 'resolved', 'dismissed']

const statusTabs = computed(() => {
  const allCount = statusKeys.reduce((sum, k) => sum + (counts.value[k] || 0), 0)
  return [
    { key: '', label: t('abuseAdmin.tabs.all'), count: allCount },
    ...statusKeys.map((k) => ({ key: k, label: t(`abuseAdmin.tabs.${k}`), count: counts.value[k] || 0 })),
  ]
})

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

const statusOptions = computed(() => [
  ...statusKeys.map((k) => ({ value: k, label: t(`abuseAdmin.status.${k}`) })),
])

const dirty = computed(() => {
  if (!selected.value) return false
  return draftStatus.value !== selected.value.status || draftNote.value !== (selected.value.admin_note || '')
})

const fetchReports = async () => {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize }
    if (statusFilter.value) params.status = statusFilter.value
    if (searchQuery.value.trim()) params.search = searchQuery.value.trim()
    const { data } = await axios.get('/api/admin/abuse-reports', { params })
    reports.value = data.reports || []
    total.value = data.total || 0
    counts.value = { open: 0, reviewing: 0, resolved: 0, dismissed: 0, ...(data.counts || {}) }
  } catch (error) {
    toast.error(error.response?.data?.error || t('abuseAdmin.toasts.fetchFailed'))
  } finally {
    loading.value = false
  }
}

let searchTimeout = null
const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchReports()
  }, 400)
}

const clearSearch = () => {
  searchQuery.value = ''
  page.value = 1
  fetchReports()
}

const filterByStatus = (key) => {
  statusFilter.value = key
  page.value = 1
  fetchReports()
}

const goToPage = (p) => {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  fetchReports()
}

const openDetail = (report) => {
  selected.value = report
  draftStatus.value = report.status
  draftNote.value = report.admin_note || ''
  showDetail.value = true
}

const closeDetail = () => {
  showDetail.value = false
  selected.value = null
}

const saveDetail = async () => {
  if (!selected.value) return
  saving.value = true
  try {
    const { data } = await axios.put(`/api/admin/abuse-reports/${selected.value.id}`, {
      status: draftStatus.value,
      admin_note: draftNote.value,
    })
    const updated = data.report
    const idx = reports.value.findIndex((r) => r.id === updated.id)
    if (idx !== -1) reports.value[idx] = updated
    selected.value = updated
    draftStatus.value = updated.status
    draftNote.value = updated.admin_note || ''
    toast.success(t('abuseAdmin.toasts.saveSuccess'))
    fetchReports()
  } catch (error) {
    toast.error(error.response?.data?.error || t('abuseAdmin.toasts.saveFailed'))
  } finally {
    saving.value = false
  }
}

const confirmDelete = (report) => {
  if (!report) return
  deleteTarget.value = report
  showDelete.value = true
}

const handleDelete = async () => {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await axios.delete(`/api/admin/abuse-reports/${deleteTarget.value.id}`)
    toast.success(t('abuseAdmin.toasts.deleteSuccess'))
    showDelete.value = false
    if (selected.value?.id === deleteTarget.value.id) closeDetail()
    deleteTarget.value = null
    if (reports.value.length === 1 && page.value > 1) page.value -= 1
    fetchReports()
  } catch (error) {
    toast.error(error.response?.data?.error || t('abuseAdmin.toasts.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

const copyUrl = async (url) => {
  try {
    await navigator.clipboard.writeText(url)
    toast.success(t('abuseAdmin.toasts.copySuccess'))
  } catch {
    toast.error(t('abuseAdmin.toasts.copyFailed'))
  }
}

const statusLabel = (s) => (statusKeys.includes(s) ? t(`abuseAdmin.status.${s}`) : s)
const categoryLabel = (c) => t(`abuseAdmin.categories.${c}`, c)
const domainStatusLabel = (s) => (s ? t(`abuseAdmin.domainStatus.${s}`, s) : t('abuseAdmin.domainStatus.unknown'))

const statusChipClass = (s) => ({
  open: 'y-chip-danger',
  reviewing: 'y-chip-warn',
  resolved: 'y-chip-mint',
  dismissed: 'y-chip',
}[s] || 'y-chip')

const domainStatusChipClass = (s) => ({
  active: 'y-chip-mint',
  suspended: 'y-chip-danger',
  expired: 'y-chip-warn',
}[s] || 'y-chip')

onMounted(fetchReports)
</script>
