<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('uiAdmin1.pending.title')"
      :description="$t('adminPendingDomains.description')"
    >
      <template #actions>
        <span class="y-chip font-mono">{{ $t('adminPendingDomains.totalCount', { count: pagination.total }) }}</span>
      </template>
    </PageHeader>

    <section class="y-panel grid grid-cols-1 divide-y divide-y-line sm:grid-cols-3 sm:divide-x sm:divide-y-0">
      <div v-for="stat in statCells" :key="stat.key" class="space-y-2 p-5">
        <p class="y-label flex items-center gap-2">
          <span class="y-dot" :class="stat.dot" aria-hidden="true"></span>
          {{ stat.label }}
        </p>
        <p class="font-mono text-3xl font-medium tracking-tight text-y-text">{{ stat.value }}</p>
      </div>
    </section>

    <section class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <label class="relative block w-full sm:max-w-sm">
          <span class="sr-only">{{ $t('common.search') }}</span>
          <Search :size="16" :stroke-width="1.5" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('adminPendingDomains.searchPlaceholder')"
            class="input input-bordered h-10 min-h-0 w-full pl-9 font-mono text-sm"
          />
        </label>
      </div>

      <YLoader v-if="loading" :label="$t('common.loading')" />

      <EmptyState v-else-if="domains.length === 0" :title="$t('adminPendingDomains.noData')" />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('adminPendingDomains.domain') }}</th>
              <th>{{ $t('adminPendingDomains.rootDomain') }}</th>
              <th>{{ $t('adminPendingDomains.status') }}</th>
              <th>{{ $t('adminPendingDomains.registeredAt') }}</th>
              <th>{{ $t('adminPendingDomains.expiresAt') }}</th>
              <th>{{ $t('adminPendingDomains.firstFailedAt') }}</th>
              <th class="text-right">{{ $t('adminPendingDomains.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="domain in domains" :key="domain.id">
              <td><span class="y-domain">{{ domain.full_domain }}</span></td>
              <td>
                <span v-if="domain.root_domain?.domain" class="font-mono text-sm text-y-muted">{{ domain.root_domain.domain }}</span>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td>
                <StatusDot :tone="pendingStatusTone(domain.status)" :label="getStatusText(domain.status)" />
              </td>
              <td class="whitespace-nowrap text-sm text-y-muted">{{ formatDate(domain.registered_at) }}</td>
              <td class="whitespace-nowrap text-sm text-y-muted">{{ formatDate(domain.expires_at) }}</td>
              <td class="whitespace-nowrap">
                <template v-if="domain.first_failed_at">
                  <span class="block text-sm text-y-danger-text">{{ formatDate(domain.first_failed_at) }}</span>
                  <span class="block text-xs text-y-faint">{{ $t('adminPendingDomains.daysAgo', { days: getDaysSince(domain.first_failed_at) }) }}</span>
                </template>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td class="text-right">
                <button
                  type="button"
                  class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                  :title="$t('adminPendingDomains.delete')"
                  :aria-label="$t('adminPendingDomains.delete')"
                  @click="deleteDomain(domain)"
                >
                  <Trash2 :size="16" :stroke-width="1.5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div
        v-if="!loading && pagination.total_pages > 1"
        class="flex flex-col items-center justify-between gap-3 border-t border-y-line px-5 py-3 sm:flex-row"
      >
        <p class="font-mono text-xs text-y-faint">{{ $t('uiAdmin1.pageOf', { page: pagination.page, total: pagination.total_pages }) }}</p>
        <div class="join">
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="pagination.page === 1"
            :aria-label="$t('common.previous')"
            @click="changePage(pagination.page - 1)"
          >
            <ChevronLeft :size="16" :stroke-width="1.5" />
          </button>
          <template v-for="(page, index) in visiblePages" :key="index">
            <button
              v-if="typeof page === 'number'"
              type="button"
              class="join-item btn btn-ghost btn-sm min-w-[2.125rem] px-2 font-mono font-medium"
              :class="page === pagination.page ? 'bg-y-surface-2 text-y-mint-text' : 'text-y-muted'"
              :aria-current="page === pagination.page ? 'page' : undefined"
              @click="changePage(page)"
            >
              {{ page }}
            </button>
            <span v-else class="join-item flex h-[2.125rem] items-center px-2 font-mono text-xs text-y-faint">…</span>
          </template>
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="pagination.page === pagination.total_pages"
            :aria-label="$t('common.next')"
            @click="changePage(pagination.page + 1)"
          >
            <ChevronRight :size="16" :stroke-width="1.5" />
          </button>
        </div>
      </div>
    </section>

    <!-- Delete Confirmation Modal -->
    <dialog ref="deleteModal" class="modal">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('adminPendingDomains.confirmDelete') }}</h3>
        <p class="mt-3 text-sm leading-6 text-y-muted">
          {{ $t('adminPendingDomains.confirmDeleteMessage', { domain: domainToDelete?.full_domain }) }}
        </p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="closeDeleteModal">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-error" :disabled="deleting" @click="confirmDelete">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            {{ $t('common.delete') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeDeleteModal">{{ $t('common.close') }}</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Search, Trash2, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()

const domains = ref([])
const searchQuery = ref('')
const loading = ref(true)
const deleteModal = ref(null)
const domainToDelete = ref(null)
const deleting = ref(false)
const pagination = ref({
  page: 1,
  per_page: 50,
  total: 0,
  total_pages: 0
})

let searchTimer = null

const pendingCount = computed(() => domains.value.filter(d => d.status === 'pending').length)
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
    const res = await axios.get('/api/admin/pending-domains', { params })
    domains.value = res.data.pending_domains || []
    if (res.data.pagination) {
      pagination.value = res.data.pagination
    }
  } catch (error) {
    toast.error(t('adminPendingDomains.fetchFailed') + ': ' + (error.response?.data?.error || error.message))
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

const statCells = computed(() => [
  { key: 'pending', dot: 'bg-y-info', label: t('adminPendingDomains.pending'), value: pendingCount.value },
  { key: 'healthy', dot: 'bg-y-mint', label: t('adminPendingDomains.healthy'), value: healthyCount.value },
  { key: 'unhealthy', dot: 'bg-y-danger', label: t('adminPendingDomains.unhealthy'), value: unhealthyCount.value },
])

const pendingStatusTone = (status) => ({ pending: 'info', healthy: 'ok', unhealthy: 'danger' }[status] || 'muted')

const getStatusText = (status) => {
  return t(`adminPendingDomains.status_${status}`)
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

const deleteDomain = (domain) => {
  domainToDelete.value = domain
  deleteModal.value?.showModal()
}

const closeDeleteModal = () => {
  deleteModal.value?.close()
  domainToDelete.value = null
}

const confirmDelete = async () => {
  if (!domainToDelete.value) return

  deleting.value = true
  try {
    await axios.delete(`/api/admin/pending-domains/${domainToDelete.value.id}`)
    toast.success(t('adminPendingDomains.deleteSuccess'))
    closeDeleteModal()
    await fetchDomains()
  } catch (error) {
    toast.error(t('adminPendingDomains.deleteFailed') + ': ' + (error.response?.data?.error || error.message))
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchDomains()
})
</script>
