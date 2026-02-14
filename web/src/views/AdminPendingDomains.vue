<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold">{{ $t('adminPendingDomains.title') }}</h1>
        <p class="text-sm opacity-70 mt-2">{{ $t('adminPendingDomains.description') }}</p>
      </div>
      <div class="text-sm opacity-70">{{ $t('adminPendingDomains.totalCount', { count: pagination.total }) }}</div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminPendingDomains.pending') }}</div>
        <div class="stat-value text-primary">{{ pendingCount }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminPendingDomains.healthy') }}</div>
        <div class="stat-value text-success">{{ healthyCount }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminPendingDomains.unhealthy') }}</div>
        <div class="stat-value text-error">{{ unhealthyCount }}</div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <input
        v-model="searchQuery"
        type="text"
        :placeholder="$t('adminPendingDomains.searchPlaceholder')"
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
            <th>{{ $t('adminPendingDomains.domain') }}</th>
            <th>{{ $t('adminPendingDomains.rootDomain') }}</th>
            <th>{{ $t('adminPendingDomains.status') }}</th>
            <th>{{ $t('adminPendingDomains.registeredAt') }}</th>
            <th>{{ $t('adminPendingDomains.expiresAt') }}</th>
            <th>{{ $t('adminPendingDomains.firstFailedAt') }}</th>
            <th>{{ $t('adminPendingDomains.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="domain in filteredDomains" :key="domain.id">
            <td class="font-mono">{{ domain.full_domain }}</td>
            <td>{{ domain.root_domain?.domain || 'N/A' }}</td>
            <td>
              <span
                class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center"
                :class="{
                  'bg-blue-100 text-blue-800': domain.status === 'pending',
                  'bg-green-100 text-green-800': domain.status === 'healthy',
                  'bg-red-100 text-red-800': domain.status === 'unhealthy'
                }"
              >
                {{ getStatusText(domain.status) }}
              </span>
            </td>
            <td>{{ formatDate(domain.registered_at) }}</td>
            <td>{{ formatDate(domain.expires_at) }}</td>
            <td>
              <span v-if="domain.first_failed_at" class="text-error">
                {{ formatDate(domain.first_failed_at) }}
                <br />
                <span class="text-xs">({{ $t('adminPendingDomains.daysAgo', { days: getDaysSince(domain.first_failed_at) }) }})</span>
              </span>
              <span v-else class="opacity-50">-</span>
            </td>
            <td>
              <button
                @click="deleteDomain(domain)"
                class="btn btn-ghost btn-sm text-error"
                :title="$t('adminPendingDomains.delete')"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </td>
          </tr>
          <tr v-if="filteredDomains.length === 0">
            <td colspan="7" class="text-center py-8 opacity-50">
              {{ $t('adminPendingDomains.noData') }}
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
        <template v-for="(page, index) in visiblePages" :key="index">
          <button
            v-if="typeof page === 'number'"
            class="join-item btn"
            :class="{ 'btn-active': page === pagination.page }"
            @click="changePage(page)"
          >
            {{ page }}
          </button>
          <button v-else class="join-item btn btn-disabled">
            ...
          </button>
        </template>
        <button
          class="join-item btn"
          :disabled="pagination.page === pagination.total_pages"
          @click="changePage(pagination.page + 1)"
        >
          »
        </button>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <dialog ref="deleteModal" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">{{ $t('adminPendingDomains.confirmDelete') }}</h3>
        <p class="py-4">
          {{ $t('adminPendingDomains.confirmDeleteMessage', { domain: domainToDelete?.full_domain }) }}
        </p>
        <div class="modal-action">
          <button @click="closeDeleteModal" class="btn">{{ $t('common.cancel') }}</button>
          <button @click="confirmDelete" class="btn btn-error" :disabled="deleting">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            {{ $t('common.delete') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button @click="closeDeleteModal">close</button>
      </form>
    </dialog>
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
const deleteModal = ref(null)
const domainToDelete = ref(null)
const deleting = ref(false)
const pagination = ref({
  page: 1,
  per_page: 50,
  total: 0,
  total_pages: 0
})

const pendingCount = computed(() => domains.value.filter(d => d.status === 'pending').length)
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
    const res = await axios.get('/api/admin/pending-domains', {
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
    toast.error(t('adminPendingDomains.fetchFailed') + ': ' + (error.response?.data?.error || error.message))
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
