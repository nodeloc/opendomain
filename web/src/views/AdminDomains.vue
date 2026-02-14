<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">{{ $t('adminDomains.title') }}</h1>
      <div class="flex items-center gap-4">
        <button
          v-if="fossbillingEnabled"
          @click="syncFOSSBilling"
          :disabled="syncing"
          class="btn btn-primary btn-sm"
        >
          <span v-if="syncing" class="loading loading-spinner loading-sm"></span>
          {{ syncing ? $t('adminDomains.syncing') : $t('adminDomains.syncFOSSBilling') }}
        </button>
        <div class="text-sm opacity-70">{{ $t('adminDomains.totalCount', { count: pagination.total }) }}</div>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminDomains.totalDomains') }}</div>
        <div class="stat-value text-primary">{{ pagination.total }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminDomains.currentPage') }}</div>
        <div class="stat-value text-info">{{ domains.length }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminDomains.pages') }}</div>
        <div class="stat-value text-secondary">{{ pagination.total_pages }}</div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <div class="relative">
        <div class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="$t('adminDomains.searchPlaceholder')"
          class="input input-bordered w-full pl-12 pr-10"
          @input="debouncedSearch"
        />
        <button
          v-if="searchQuery"
          @click="clearSearch"
          class="absolute inset-y-0 right-0 flex items-center pr-4 text-gray-400 hover:text-gray-600"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <p v-if="searchQuery" class="text-sm text-gray-500 mt-2">
        {{ $t('adminDomains.searchResults', { count: domains.length, query: searchQuery }) }}
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Domains Table -->
    <div v-else class="overflow-x-auto bg-base-100 rounded-lg shadow">
      <table class="table w-full">
        <thead>
          <tr>
            <th>ID</th>
            <th>{{ $t('adminDomains.domain') }}</th>
            <th>{{ $t('adminDomains.owner') }}</th>
            <th>{{ $t('adminDomains.rootDomain') }}</th>
            <th>{{ $t('adminDomains.status') }}</th>
            <th>{{ $t('adminDomains.registeredAt') }}</th>
            <th>{{ $t('adminDomains.expiresAt') }}</th>
            <th>{{ $t('adminDomains.autoRenew') }}</th>
            <th>{{ $t('adminDomains.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="domain in domains" :key="domain.id">
            <td>{{ domain.id }}</td>
            <td>
              <div class="font-bold font-mono">{{ domain.full_domain }}</div>
              <div v-if="!domain.use_default_nameservers" class="text-xs opacity-50">
                {{ $t('adminDomains.customNS') }}
              </div>
            </td>
            <td>
              <div v-if="domain.user">
                <div class="font-semibold">{{ domain.user.username }}</div>
                <div class="text-xs opacity-50">{{ domain.user.email }}</div>
              </div>
              <span v-else class="text-gray-500">-</span>
            </td>
            <td>
              <span v-if="domain.root_domain" class="font-mono text-sm">
                {{ domain.root_domain.domain }}
              </span>
            </td>
            <td>
              <span :class="['px-2 py-0.5 rounded text-xs font-medium',
                domain.status === 'active' ? 'bg-green-100 text-green-800' :
                domain.status === 'suspended' ? 'bg-orange-100 text-orange-800' :
                domain.status === 'expired' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-800']">
                {{ domain.status }}
              </span>
            </td>
            <td>{{ formatDate(domain.registered_at) }}</td>
            <td>
              <div class="flex flex-col gap-1">
                <span :class="getExpiryClass(domain.expires_at)">
                  {{ formatDate(domain.expires_at) }}
                </span>
                <span v-if="getExpiryStatus(domain.expires_at)" :class="getExpiryBadgeClass(domain.expires_at)" class="badge badge-sm">
                  {{ getExpiryStatus(domain.expires_at) }}
                </span>
              </div>
            </td>
            <td>
              <span v-if="domain.auto_renew" class="text-success">✓</span>
              <span v-else class="text-gray-400">-</span>
            </td>
            <td>
              <div class="flex gap-1">
                <button @click="viewDomainDetails(domain)" class="btn btn-sm btn-ghost" :title="$t('adminDomains.detailsTooltip')">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                    <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                  </svg>
                </button>
                <button
                  @click="toggleDomainStatus(domain)"
                  class="btn btn-sm btn-ghost"
                  :title="domain.status === 'suspended' ? $t('adminDomains.activateTooltip') : $t('adminDomains.suspendTooltip')"
                >
                  <svg v-if="domain.status === 'suspended'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-warning" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
                <button @click="confirmDelete(domain)" class="btn btn-sm btn-ghost text-error" :title="$t('adminDomains.deleteTooltip')">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="domains.length === 0" class="text-center py-8 text-gray-500">
        {{ searchQuery ? $t('adminDomains.noSearchResults') : $t('adminDomains.noData') }}
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="!loading && pagination.total_pages > 1" class="flex justify-center mt-6">
      <div class="join">
        <button
          class="join-item btn"
          :disabled="pagination.page === 1"
          @click="goToPage(pagination.page - 1)"
        >
          «
        </button>
        <template v-for="(page, index) in visiblePages" :key="index">
          <button
            v-if="typeof page === 'number'"
            class="join-item btn"
            :class="{ 'btn-active': page === pagination.page }"
            @click="goToPage(page)"
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
          @click="goToPage(pagination.page + 1)"
        >
          »
        </button>
      </div>
    </div>

    <!-- Domain Details Modal -->
    <dialog class="modal" :class="{ 'modal-open': showDetailsModal }">
      <div class="modal-box max-w-3xl">
        <h3 class="font-bold text-lg mb-4">{{ $t('adminDomains.domainDetails') }}</h3>
        <div v-if="selectedDomain" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.domainId') }}</span></label>
              <div class="text-lg">{{ selectedDomain.id }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.fullDomain') }}</span></label>
              <div class="text-lg font-mono">{{ selectedDomain.full_domain }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.subdomain') }}</span></label>
              <div class="text-lg font-mono">{{ selectedDomain.subdomain }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.rootDomain') }}</span></label>
              <div class="text-lg font-mono">{{ selectedDomain.root_domain?.domain }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.status') }}</span></label>
              <div class="text-lg">{{ selectedDomain.status }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.owner') }}</span></label>
              <div class="text-lg">
                {{ selectedDomain.user?.username }} ({{ selectedDomain.user?.email }})
              </div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.registeredAt') }}</span></label>
              <div class="text-lg">{{ formatDate(selectedDomain.registered_at) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.expiresAt') }}</span></label>
              <div class="text-lg">{{ formatDate(selectedDomain.expires_at) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.autoRenew') }}</span></label>
              <div class="text-lg">{{ selectedDomain.auto_renew ? $t('common.yes') : $t('common.no') }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.dnsSynced') }}</span></label>
              <div class="text-lg">{{ selectedDomain.dns_synced ? $t('common.yes') : $t('common.no') }}</div>
            </div>
            <div class="col-span-2" v-if="selectedDomain.nameservers">
              <label class="label"><span class="label-text font-semibold">{{ $t('adminDomains.nameservers') }}</span></label>
              <div class="text-sm font-mono bg-base-200 p-2 rounded">
                {{ formatNameservers(selectedDomain.nameservers) }}
              </div>
            </div>
          </div>
        </div>
        <div class="modal-action">
          <button @click="showDetailsModal = false" class="btn">{{ $t('common.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showDetailsModal = false">
        <button>close</button>
      </form>
    </dialog>

    <!-- Delete Confirmation Modal -->
    <dialog class="modal" :class="{ 'modal-open': showDeleteModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ $t('adminDomains.confirmDelete') }}</h3>
        <p class="py-4">
          {{ $t('adminDomains.confirmDeleteMessage', { domain: domainToDelete?.full_domain }) }}
        </p>
        <div class="modal-action">
          <button @click="showDeleteModal = false" class="btn">{{ $t('common.cancel') }}</button>
          <button @click="handleDelete" class="btn btn-error" :disabled="deleting">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            <span v-else>{{ $t('common.delete') }}</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showDeleteModal = false">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSiteConfigStore } from '../stores/siteConfig'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const toast = useToast()
const siteConfig = useSiteConfigStore()

const allDomains = ref([]) // 存储所有域名
const domains = ref([]) // 显示的域名（可能经过过滤）
const searchQuery = ref('')
const loading = ref(true)
const showDetailsModal = ref(false)
const showDeleteModal = ref(false)
const selectedDomain = ref(null)
const domainToDelete = ref(null)
const deleting = ref(false)
const syncing = ref(false) // FOSSBilling同步状态
const pagination = ref({
  page: 1,
  page_size: 20,
  total: 0,
  total_pages: 0
})
let searchTimeout = null

const fossbillingEnabled = computed(() => siteConfig.fossbilling.enabled)

const fetchDomains = async (search = '') => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.page_size
    }
    if (search) {
      params.search = search
    }
    const response = await axios.get('/api/admin/domains', { params })
    domains.value = response.data.domains || []
    if (response.data.pagination) {
      pagination.value = response.data.pagination
    }
    // 首次加载时也获取所有域名用于统计（只在第一页时）
    if (pagination.value.page === 1 && !search) {
      allDomains.value = domains.value
    }
  } catch (error) {
    toast.error(t('adminDomains.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    fetchDomains(searchQuery.value)
  }, 500)
}

const clearSearch = () => {
  searchQuery.value = ''
  pagination.value.page = 1
  fetchDomains()
}

const goToPage = (page) => {
  if (page >= 1 && page <= pagination.value.total_pages) {
    pagination.value.page = page
    fetchDomains(searchQuery.value)
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

const viewDomainDetails = (domain) => {
  selectedDomain.value = domain
  showDetailsModal.value = true
}

const confirmDelete = (domain) => {
  domainToDelete.value = domain
  showDeleteModal.value = true
}

const handleDelete = async () => {
  deleting.value = true
  try {
    await axios.delete(`/api/admin/domains/${domainToDelete.value.id}`)
    toast.success(t('adminDomains.deleteSuccess'))
    showDeleteModal.value = false
    await fetchDomains(searchQuery.value)
  } catch (error) {
    toast.error(error.response?.data?.error || t('adminDomains.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

const toggleDomainStatus = async (domain) => {
  try {
    const newStatus = domain.status === 'suspended' ? 'active' : 'suspended'
    await axios.put(`/api/admin/domains/${domain.id}/status`, { status: newStatus })
    toast.success(t('adminDomains.statusUpdateSuccess'))
    await fetchDomains(searchQuery.value)
  } catch (error) {
    toast.error(error.response?.data?.error || t('adminDomains.statusUpdateFailed'))
  }
}

const syncFOSSBilling = async () => {
  if (!confirm(t('adminDomains.confirmSync'))) {
    return
  }

  syncing.value = true
  try {
    // 设置5分钟超时，因为同步大量域名需要时间
    const res = await axios.post('/api/admin/sync-fossbilling-domains', {}, {
      timeout: 300000 // 5分钟
    })
    const data = res.data

    // 显示同步结果
    const details = data.details || []
    const summary = `${t('adminDomains.syncComplete')}\n` +
      `${t('adminDomains.synced')}: ${data.synced_count}\n` +
      `${t('adminDomains.skipped')}: ${data.skipped_count}\n` +
      `${t('adminDomains.existing')}: ${data.existing_count}\n` +
      `${t('adminDomains.errors')}: ${data.error_count}`

    if (data.synced_count > 0 || data.skipped_count > 0) {
      toast.success(summary)
      // 刷新域名列表
      await fetchDomains(searchQuery.value)
    } else {
      toast.info(summary)
    }

    // 如果有详细信息，显示在控制台
    if (details.length > 0) {
      console.log('FOSSBilling Sync Details:', details)
    }
  } catch (error) {
    console.error('Sync error:', error)
    toast.error(error.response?.data?.error || t('adminDomains.syncFailed'))
  } finally {
    syncing.value = false
  }
}

const isExpiringSoon = (expiresAt) => {
  const now = new Date()
  const thirtyDaysFromNow = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000)
  const expiryDate = new Date(expiresAt)
  return expiryDate <= thirtyDaysFromNow && expiryDate > now
}

// 获取过期状态文本
const getExpiryStatus = (expiresAt) => {
  const now = new Date()
  const expiryDate = new Date(expiresAt)
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24))

  if (daysUntilExpiry > 30) {
    return null // 超过 30 天不显示状态
  } else if (daysUntilExpiry > 0 && daysUntilExpiry <= 7) {
    return t('adminDomains.expiringInDays', { days: daysUntilExpiry })
  } else if (daysUntilExpiry > 7 && daysUntilExpiry <= 30) {
    return t('adminDomains.expiringSoon')
  } else if (daysUntilExpiry <= 0 && daysUntilExpiry >= -20) {
    return t('adminDomains.expired', { days: Math.abs(daysUntilExpiry) })
  } else if (daysUntilExpiry < -20) {
    return t('adminDomains.willBeDeleted', { days: 30 - Math.abs(daysUntilExpiry) })
  }
  return null
}

// 获取过期状态样式类
const getExpiryBadgeClass = (expiresAt) => {
  const now = new Date()
  const expiryDate = new Date(expiresAt)
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24))

  if (daysUntilExpiry > 7 && daysUntilExpiry <= 30) {
    return 'badge-warning' // 即将过期 (7-30 天)
  } else if (daysUntilExpiry > 0 && daysUntilExpiry <= 7) {
    return 'badge-error' // 即将过期 (0-7 天)
  } else if (daysUntilExpiry <= 0 && daysUntilExpiry >= -20) {
    return 'badge-error' // 已过期 (0-20 天)
  } else if (daysUntilExpiry < -20) {
    return 'badge-error animate-pulse' // 即将删除 (>20 天)
  }
  return ''
}

// 获取过期日期文本样式
const getExpiryClass = (expiresAt) => {
  const now = new Date()
  const expiryDate = new Date(expiresAt)
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24))

  if (daysUntilExpiry <= 7 && daysUntilExpiry > 0) {
    return 'text-error font-semibold' // 7 天内过期
  } else if (daysUntilExpiry <= 30 && daysUntilExpiry > 7) {
    return 'text-warning font-semibold' // 30 天内过期
  } else if (daysUntilExpiry <= 0) {
    return 'text-error font-bold' // 已过期
  }
  return ''
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatNameservers = (nameservers) => {
  try {
    const ns = JSON.parse(nameservers)
    return Array.isArray(ns) ? ns.join(', ') : nameservers
  } catch {
    return nameservers
  }
}

onMounted(() => {
  fetchDomains()
})
</script>
