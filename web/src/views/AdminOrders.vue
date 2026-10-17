<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('uiAdmin2.eyebrow')"
      :title="$t('uiAdmin2.orders.title')"
      :description="$t('uiAdmin2.orders.description')"
    >
      <template #actions>
        <span class="font-mono text-sm text-y-faint">{{ $t('uiAdmin2.orders.totalCount', { count: pagination.total || 0 }) }}</span>
      </template>
    </PageHeader>

    <!-- Stats -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-5">
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.orders.stats.total') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ pagination.total || 0 }}</p>
      </div>
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.orders.stats.pending') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ getStatusCount('pending') }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.orders.stats.completed') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ getStatusCount('completed') }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.orders.stats.cancelled') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ getStatusCount('cancelled') }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
      <div class="y-panel col-span-2 p-5 lg:col-span-1">
        <p class="y-label">{{ $t('uiAdmin2.orders.stats.revenue') }}</p>
        <p class="mt-2 truncate font-mono text-3xl font-medium tracking-tight text-y-text">{{ formatPrice(totalRevenue) }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
    </div>

    <!-- Orders -->
    <section class="y-panel">
      <div class="y-panel-header">
        <div class="max-w-full overflow-x-auto">
          <div class="tabs tabs-boxed tabs-sm w-max" role="tablist" :aria-label="$t('uiAdmin2.orders.filterLabel')">
            <button
              v-for="option in statusFilters"
              :key="option.value || 'all'"
              type="button"
              role="tab"
              class="tab"
              :class="{ 'tab-active': filterStatus === option.value }"
              :aria-selected="filterStatus === option.value"
              @click="filterStatus = option.value"
            >
              {{ option.label }}
            </button>
          </div>
        </div>

        <label class="input input-bordered input-sm flex h-10 w-full items-center gap-2 sm:max-w-sm">
          <Search :size="16" :stroke-width="1.5" class="shrink-0 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('admin.orderManagement.searchPlaceholder')"
            class="min-w-0 grow bg-transparent text-sm"
            @input="debouncedSearch"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="btn btn-ghost btn-xs btn-square"
            :aria-label="$t('uiAdmin2.clearSearch')"
            :title="$t('uiAdmin2.clearSearch')"
            @click="clearSearch"
          >
            <X :size="14" :stroke-width="1.5" />
          </button>
        </label>
      </div>

      <p v-if="searchQuery" class="border-b border-y-line px-5 py-2.5 text-sm text-y-muted">
        {{ $t('admin.orderManagement.searchResults', { count: orders.length, query: searchQuery }) }}
      </p>

      <YLoader v-if="loading && orders.length === 0" :label="$t('uiAdmin2.loading')" />

      <EmptyState
        v-else-if="orders.length === 0"
        :title="$t('uiAdmin2.orders.emptyTitle')"
        :description="searchQuery || filterStatus ? $t('uiAdmin2.orders.emptyFilteredDescription') : $t('uiAdmin2.orders.emptyDescription')"
      >
        <button v-if="searchQuery" type="button" class="btn btn-outline btn-sm" @click="clearSearch">
          {{ $t('uiAdmin2.clearSearch') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('admin.orderManagement.table.orderNumber') }}</th>
              <th>{{ $t('admin.orderManagement.table.user') }}</th>
              <th>{{ $t('admin.orderManagement.table.domain') }}</th>
              <th>{{ $t('admin.orderManagement.table.rootDomain') }}</th>
              <th class="text-right">{{ $t('admin.orderManagement.table.amount') }}</th>
              <th>{{ $t('admin.orderManagement.table.status') }}</th>
              <th>{{ $t('admin.orderManagement.table.createdAt') }}</th>
              <th class="text-right">{{ $t('admin.orderManagement.table.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in orders" :key="order.id">
              <td class="whitespace-nowrap font-mono text-xs text-y-text">{{ order.order_number }}</td>
              <td>
                <template v-if="order.user">
                  <p class="font-medium text-y-text">{{ order.user.username }}</p>
                  <p class="font-mono text-xs text-y-faint">ID {{ order.user_id }}</p>
                </template>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td><span class="y-domain text-sm">{{ order.subdomain }}</span></td>
              <td>
                <span v-if="order.root_domain" class="font-mono text-sm text-y-muted">{{ order.root_domain.domain }}</span>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td class="whitespace-nowrap text-right">
                <p class="font-mono font-medium text-y-text">{{ formatPrice(Number(order.final_price || 0)) }}</p>
                <p v-if="order.base_price !== undefined && order.base_price !== null" class="font-mono text-xs text-y-faint">
                  {{ $t('admin.orderManagement.table.originalPrice') }} {{ formatPrice(Number(order.base_price || 0)) }}
                </p>
              </td>
              <td>
                <span :class="getStatusBadgeClass(order.status)">
                  {{ $t(`admin.orderManagement.status.${order.status}`) }}
                </span>
              </td>
              <td class="whitespace-nowrap text-y-muted">{{ formatDate(order.created_at) }}</td>
              <td class="text-right">
                <button
                  type="button"
                  class="btn btn-ghost btn-sm btn-square"
                  :title="$t('admin.orderManagement.table.viewDetails')"
                  :aria-label="$t('admin.orderManagement.table.viewDetails')"
                  @click="viewOrderDetails(order)"
                >
                  <Eye :size="16" :stroke-width="1.5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Pagination -->
    <div v-if="pagination.total_page > 1" class="flex flex-col items-center justify-between gap-3 sm:flex-row">
      <p class="font-mono text-xs text-y-faint">{{ $t('uiAdmin2.pageOf', { page: currentPage, total: pagination.total_page }) }}</p>
      <div class="join">
        <button
          type="button"
          class="join-item btn btn-outline btn-sm disabled:border-y-line disabled:bg-transparent disabled:text-y-faint"
          :disabled="currentPage === 1"
          :aria-label="$t('uiAdmin2.previousPage')"
          @click="goToPage(currentPage - 1)"
        >
          <ChevronLeft :size="16" :stroke-width="1.5" />
        </button>
        <template v-for="(page, index) in visiblePages" :key="index">
          <button
            v-if="typeof page === 'number'"
            type="button"
            class="join-item btn btn-outline btn-sm font-mono"
            :class="{ 'bg-y-surface-2 !text-y-mint-text': page === currentPage }"
            :aria-current="page === currentPage ? 'page' : undefined"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
          <button v-else type="button" class="join-item btn btn-outline btn-sm btn-disabled font-mono" tabindex="-1">…</button>
        </template>
        <button
          type="button"
          class="join-item btn btn-outline btn-sm disabled:border-y-line disabled:bg-transparent disabled:text-y-faint"
          :disabled="currentPage === pagination.total_page"
          :aria-label="$t('uiAdmin2.nextPage')"
          @click="goToPage(currentPage + 1)"
        >
          <ChevronRight :size="16" :stroke-width="1.5" />
        </button>
      </div>
    </div>

    <!-- Order Details Modal -->
    <div v-if="showDetailsModal" class="modal modal-open">
      <div class="modal-box relative max-w-2xl">
        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square absolute right-3 top-3"
          :aria-label="$t('uiAdmin2.close')"
          @click="showDetailsModal = false"
        >
          <X :size="16" :stroke-width="1.5" />
        </button>
        <h3 class="text-lg font-semibold text-y-text">{{ $t('admin.orderManagement.details.title') }}</h3>

        <div v-if="selectedOrder" class="mt-6 space-y-4">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <p class="font-mono text-sm text-y-text break-all">{{ selectedOrder.order_number }}</p>
            <span :class="getStatusBadgeClass(selectedOrder.status)">
              {{ $t(`admin.orderManagement.status.${selectedOrder.status}`) }}
            </span>
          </div>

          <dl class="grid grid-cols-1 gap-x-6 gap-y-4 rounded-md border border-y-line p-4 sm:grid-cols-2">
            <div>
              <dt class="y-label">{{ $t('admin.orderManagement.details.domain') }}</dt>
              <dd class="mt-1 y-domain text-sm">{{ selectedOrder.subdomain }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.orderManagement.details.rootDomain') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedOrder.root_domain?.domain || '—' }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.orderManagement.details.user') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ selectedOrder.user?.username || '—' }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.orderManagement.details.userId') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedOrder.user_id }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.orderManagement.details.years') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ $t('admin.orderManagement.details.yearsValue', { count: selectedOrder.years }) }}</dd>
            </div>
            <div v-if="selectedOrder.coupon_code">
              <dt class="y-label">{{ $t('admin.orderManagement.details.couponCode') }}</dt>
              <dd class="mt-1"><span class="y-tag">{{ selectedOrder.coupon_code }}</span></dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.orderManagement.details.createdAt') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ formatDate(selectedOrder.created_at) }}</dd>
            </div>
            <div v-if="selectedOrder.updated_at">
              <dt class="y-label">{{ $t('admin.orderManagement.details.updatedAt') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ formatDate(selectedOrder.updated_at) }}</dd>
            </div>
          </dl>

          <dl class="y-inset divide-y divide-y-line text-sm">
            <div class="flex items-center justify-between gap-4 px-4 py-3">
              <dt class="text-y-muted">{{ $t('admin.orderManagement.details.originalPrice') }}</dt>
              <dd class="font-mono text-y-text">{{ formatPrice(Number(selectedOrder.base_price || 0)) }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4 px-4 py-3">
              <dt class="text-y-muted">{{ $t('admin.orderManagement.details.discount') }}</dt>
              <dd class="font-mono text-y-text">-{{ formatPrice(Number(selectedOrder.discount_amount || 0)) }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4 px-4 py-3">
              <dt class="font-medium text-y-text">{{ $t('admin.orderManagement.details.finalAmount') }}</dt>
              <dd class="font-mono text-base font-medium text-y-mint-text">{{ formatPrice(Number(selectedOrder.final_price || 0)) }}</dd>
            </div>
          </dl>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDetailsModal = false">{{ $t('uiAdmin2.close') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { Search, X, Eye, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const { formatPrice } = useCurrency()

const orders = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(20)
const pagination = ref({
  page: 1,
  page_size: 20,
  total: 0,
  total_page: 0
})
const filterStatus = ref('')
const searchQuery = ref('')
const showDetailsModal = ref(false)
const selectedOrder = ref(null)
let searchTimeout = null

const totalRevenue = computed(() => {
  return orders.value
    .filter(o => o.status === 'completed')
    .reduce((sum, o) => sum + parseFloat(o.final_price || 0), 0)
})

const visiblePages = computed(() => {
  const current = currentPage.value
  const total = pagination.value.total_page
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

const getStatusCount = (status) => {
  return orders.value.filter(o => o.status === status).length
}

const getStatusBadgeClass = (status) => {
  const classMap = {
    pending: 'y-chip-warn',
    processing: 'y-chip-info',
    paid: 'y-chip-mint',
    completed: 'y-chip-mint',
    cancelled: 'y-chip-danger',
    failed: 'y-chip-danger'
  }
  return classMap[status] || 'y-chip'
}

const statusFilters = computed(() => [
  { value: '', label: t('admin.orderManagement.filters.all') },
  { value: 'pending', label: t('admin.orderManagement.filters.pending') },
  { value: 'processing', label: t('admin.orderManagement.filters.processing') },
  { value: 'completed', label: t('admin.orderManagement.filters.completed') },
  { value: 'cancelled', label: t('admin.orderManagement.filters.cancelled') }
])

const fetchOrders = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (searchQuery.value) {
      params.search = searchQuery.value
    }

    const response = await axios.get('/api/admin/orders', { params })
    orders.value = response.data.orders || []
    pagination.value = response.data.pagination || {}
  } catch (error) {
    toast.error(t('admin.orderManagement.fetchError'))
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchOrders()
  }, 500)
}

const clearSearch = () => {
  searchQuery.value = ''
  currentPage.value = 1
  fetchOrders()
}

const goToPage = (page) => {
  if (page >= 1 && page <= pagination.value.total_page) {
    currentPage.value = page
    fetchOrders()
  }
}

const viewOrderDetails = (order) => {
  selectedOrder.value = order
  showDetailsModal.value = true
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

watch(filterStatus, () => {
  currentPage.value = 1
  fetchOrders()
})

onMounted(() => {
  fetchOrders()
})
</script>
