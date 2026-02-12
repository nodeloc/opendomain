<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">{{ $t('admin.orderManagement.title') }}</h1>
      <div class="text-sm opacity-70">{{ $t('admin.orderManagement.totalCount', { count: pagination.total }) }}</div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('admin.orderManagement.stats.totalOrders') }}</div>
        <div class="stat-value text-sm">{{ pagination.total }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('admin.orderManagement.stats.pendingOrders') }}</div>
        <div class="stat-value text-sm text-warning">{{ getStatusCount('pending') }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('admin.orderManagement.stats.completedOrders') }}</div>
        <div class="stat-value text-sm text-success">{{ getStatusCount('completed') }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('admin.orderManagement.stats.cancelledOrders') }}</div>
        <div class="stat-value text-sm text-error">{{ getStatusCount('cancelled') }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('admin.orderManagement.stats.totalRevenue') }}</div>
        <div class="stat-value text-sm text-primary">{{ formatPrice(totalRevenue) }}</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="mb-4 flex gap-2">
      <button 
        @click="filterStatus = ''"
        :class="['px-4 py-2 rounded', filterStatus === '' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.orderManagement.filters.all') }}
      </button>
      <button 
        @click="filterStatus = 'pending'"
        :class="['px-4 py-2 rounded', filterStatus === 'pending' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.orderManagement.filters.pending') }}
      </button>
      <button 
        @click="filterStatus = 'processing'"
        :class="['px-4 py-2 rounded', filterStatus === 'processing' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.orderManagement.filters.processing') }}
      </button>
      <button 
        @click="filterStatus = 'completed'"
        :class="['px-4 py-2 rounded', filterStatus === 'completed' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.orderManagement.filters.completed') }}
      </button>
      <button 
        @click="filterStatus = 'cancelled'"
        :class="['px-4 py-2 rounded', filterStatus === 'cancelled' ? 'bg-primary text-white' : 'bg-base-200']"
      >
        {{ $t('admin.orderManagement.filters.cancelled') }}
      </button>
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
          :placeholder="$t('admin.orderManagement.searchPlaceholder')"
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
        {{ $t('admin.orderManagement.searchResults', { count: orders.length, query: searchQuery }) }}
      </p>
    </div>

    <!-- Orders Table -->
    <div class="overflow-x-auto bg-base-100 rounded-lg shadow">
      <table class="table w-full">
        <thead>
          <tr>
            <th>{{ $t('admin.orderManagement.table.orderNumber') }}</th>
            <th>{{ $t('admin.orderManagement.table.user') }}</th>
            <th>{{ $t('admin.orderManagement.table.domain') }}</th>
            <th>{{ $t('admin.orderManagement.table.rootDomain') }}</th>
            <th>{{ $t('admin.orderManagement.table.amount') }}</th>
            <th>{{ $t('admin.orderManagement.table.status') }}</th>
            <th>{{ $t('admin.orderManagement.table.createdAt') }}</th>
            <th>{{ $t('admin.orderManagement.table.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in orders" :key="order.id">
            <td>
              <div class="font-mono text-sm">{{ order.order_number }}</div>
            </td>
            <td>
              <div v-if="order.user">
                <div class="font-bold">{{ order.user.username }}</div>
                <div class="text-xs opacity-50">ID: {{ order.user_id }}</div>
              </div>
              <div v-else class="text-gray-500">-</div>
            </td>
            <td>
              <div class="font-medium">{{ order.subdomain }}</div>
            </td>
            <td>
              <div v-if="order.root_domain">{{ order.root_domain.domain }}</div>
              <div v-else class="text-gray-500">-</div>
            </td>
            <td>
              <div class="font-bold">{{ formatPrice(Number(order.final_price || 0)) }}</div>
              <div v-if="order.base_price !== undefined && order.base_price !== null" class="text-xs opacity-50">
                {{ $t('admin.orderManagement.table.originalPrice') }}: {{ formatPrice(Number(order.base_price || 0)) }}
              </div>
            </td>
            <td>
              <span :class="getStatusBadgeClass(order.status)">
                {{ $t(`admin.orderManagement.status.${order.status}`) }}
              </span>
            </td>
            <td>{{ formatDate(order.created_at) }}</td>
            <td>
              <button 
                @click="viewOrderDetails(order)" 
                class="btn btn-sm btn-ghost"
                :title="$t('admin.orderManagement.table.viewDetails')"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="orders.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('admin.orderManagement.noData') }}
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="pagination.total_page > 1" class="flex justify-center mt-6 gap-2">
      <button 
        @click="goToPage(currentPage - 1)" 
        :disabled="currentPage === 1"
        class="btn btn-sm"
      >
        {{ $t('admin.orderManagement.pagination.previous') }}
      </button>
      <button 
        v-for="page in visiblePages" 
        :key="page"
        @click="goToPage(page)"
        :class="['btn btn-sm', currentPage === page ? 'btn-primary' : '']"
      >
        {{ page }}
      </button>
      <button 
        @click="goToPage(currentPage + 1)" 
        :disabled="currentPage === pagination.total_page"
        class="btn btn-sm"
      >
        {{ $t('admin.orderManagement.pagination.next') }}
      </button>
    </div>

    <!-- Order Details Modal -->
    <div v-if="showDetailsModal" class="modal modal-open">
      <div class="modal-box max-w-2xl">
        <h3 class="font-bold text-lg mb-4">{{ $t('admin.orderManagement.details.title') }}</h3>
        <div v-if="selectedOrder" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.orderNumber') }}</span></label>
              <div class="text-lg font-mono">{{ selectedOrder.order_number }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.status') }}</span></label>
              <span :class="getStatusBadgeClass(selectedOrder.status)">
                {{ $t(`admin.orderManagement.status.${selectedOrder.status}`) }}
              </span>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.user') }}</span></label>
              <div class="text-lg">{{ selectedOrder.user?.username || '-' }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.userId') }}</span></label>
              <div class="text-lg">{{ selectedOrder.user_id }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.domain') }}</span></label>
              <div class="text-lg font-bold">{{ selectedOrder.subdomain }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.rootDomain') }}</span></label>
              <div class="text-lg">{{ selectedOrder.root_domain?.domain || '-' }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.years') }}</span></label>
              <div class="text-lg">{{ $t('admin.orderManagement.details.yearsValue', { count: selectedOrder.years }) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.originalPrice') }}</span></label>
              <div class="text-lg">{{ formatPrice(Number(selectedOrder.base_price || 0)) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.discount') }}</span></label>
              <div class="text-lg text-error">-{{ formatPrice(Number(selectedOrder.discount_amount || 0)) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.finalAmount') }}</span></label>
              <div class="text-lg font-bold text-primary">{{ formatPrice(Number(selectedOrder.final_price || 0)) }}</div>
            </div>
            <div v-if="selectedOrder.coupon_code">
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.couponCode') }}</span></label>
              <div class="text-lg font-mono">{{ selectedOrder.coupon_code }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.createdAt') }}</span></label>
              <div class="text-lg">{{ formatDate(selectedOrder.created_at) }}</div>
            </div>
            <div v-if="selectedOrder.updated_at">
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.orderManagement.details.updatedAt') }}</span></label>
              <div class="text-lg">{{ formatDate(selectedOrder.updated_at) }}</div>
            </div>
          </div>
        </div>
        <div class="modal-action">
          <button @click="showDetailsModal = false" class="btn">{{ $t('common.close') }}</button>
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

const { t } = useI18n()
const toast = useToast()
const { formatPrice } = useCurrency()

const orders = ref([])
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
  const total = pagination.value.total_page
  const current = currentPage.value
  const delta = 2
  const pages = []
  
  for (let i = Math.max(1, current - delta); i <= Math.min(total, current + delta); i++) {
    pages.push(i)
  }
  
  return pages
})

const getStatusCount = (status) => {
  return orders.value.filter(o => o.status === status).length
}

const getStatusBadgeClass = (status) => {
  const classMap = {
    pending: 'px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800',
    processing: 'px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800',
    completed: 'px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800',
    cancelled: 'px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800',
    failed: 'px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800'
  }
  return classMap[status] || 'px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800'
}

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
