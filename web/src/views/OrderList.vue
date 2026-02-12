<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl">
    <!-- 标题 -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold">{{ $t('order.myOrders') }}</h1>
        <p class="text-lg opacity-70 mt-2">{{ $t('order.viewManage') }}</p>
      </div>
    </div>

    <!-- 状态筛选 -->
    <div class="tabs tabs-boxed mb-6">
      <a
        class="tab"
        :class="{ 'tab-active': statusFilter === '' }"
        @click="statusFilter = ''; fetchOrders()"
      >
        {{ $t('order.all') }}
      </a>
      <a
        class="tab"
        :class="{ 'tab-active': statusFilter === 'pending' }"
        @click="statusFilter = 'pending'; fetchOrders()"
      >
        {{ $t('order.pending') }}
      </a>
      <a
        class="tab"
        :class="{ 'tab-active': statusFilter === 'paid' }"
        @click="statusFilter = 'paid'; fetchOrders()"
      >
        {{ $t('order.paid') }}
      </a>
      <a
        class="tab"
        :class="{ 'tab-active': statusFilter === 'cancelled' }"
        @click="statusFilter = 'cancelled'; fetchOrders()"
      >
        {{ $t('order.cancelled') }}
      </a>
      <a
        class="tab"
        :class="{ 'tab-active': statusFilter === 'expired' }"
        @click="statusFilter = 'expired'; fetchOrders()"
      >
        {{ $t('order.expired') }}
      </a>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- 空状态 -->
    <div v-else-if="orders.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 opacity-40 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="card-title text-2xl mb-2">{{ $t('order.noOrders') }}</h3>
        <p class="mb-4 opacity-70">{{ $t('order.noOrdersDesc') }}</p>
        <button class="btn btn-primary" @click="$router.push('/')">
          {{ $t('order.registerDomain') }}
        </button>
      </div>
    </div>

    <!-- 订单列表 -->
    <div v-else class="space-y-4">
      <div
        v-for="order in orders"
        :key="order.id"
        class="card bg-base-100 shadow-xl border border-base-300 hover:shadow-2xl transition-all"
      >
        <div class="card-body">
          <div class="flex justify-between items-start flex-wrap gap-4">
            <!-- 订单信息 -->
            <div class="flex-1 min-w-[250px]">
              <div class="flex items-center gap-3 mb-2">
                <h2 class="card-title text-xl font-mono">{{ order.full_domain }}</h2>
                <span class="px-2 py-0.5 rounded text-xs font-medium" :class="getStatusBadgeClass(order.status)">
                  {{ $t('order.' + order.status) }}
                </span>
              </div>
              <div class="text-sm opacity-70 space-y-1">
                <div>{{ $t('order.orderNumber') }}{{ order.order_number }}</div>
                <div>{{ $t('order.created') }}: {{ formatDate(order.created_at) }}</div>
                <div v-if="order.expires_at && order.status === 'pending'">
                  {{ $t('order.expires') }}: {{ formatDate(order.expires_at) }}
                </div>
                <div v-if="order.paid_at">
                  {{ $t('order.paidAt') }}: {{ formatDate(order.paid_at) }}
                </div>
              </div>
            </div>

            <!-- 价格信息 -->
            <div class="text-right">
              <div class="text-2xl font-bold font-mono text-primary">
                {{ formatPrice(order.final_price) }}
              </div>
              <div v-if="order.discount_amount > 0" class="text-sm text-success mt-1">
                {{ $t('order.discount') }} {{ formatPrice(order.discount_amount) }}
              </div>
              <div class="text-sm opacity-70 mt-1">
                {{ order.is_lifetime ? $t('order.lifetime') : `${order.years} ${$t('order.year')}` }}
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="card-actions justify-end mt-4">
            <button
              v-if="order.status === 'pending' && !isExpired(order)"
              class="btn btn-primary btn-sm"
              @click="payOrder(order)"
              :disabled="paying === order.id"
            >
              <span v-if="paying === order.id" class="loading loading-spinner loading-xs"></span>
              <span v-else>{{ $t('order.payNow') }}</span>
            </button>
            <button
              v-if="order.status === 'pending' || order.status === 'expired'"
              class="btn btn-error btn-sm btn-outline"
              @click="cancelOrder(order)"
              :disabled="cancelling === order.id"
            >
              <span v-if="cancelling === order.id" class="loading loading-spinner loading-xs"></span>
              <span v-else>{{ $t('common.cancel') }}</span>
            </button>
            <button class="btn btn-ghost btn-sm" @click="viewDetails(order)">
              {{ $t('order.viewDetails') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="flex justify-center mt-8">
      <div class="join">
        <button
          class="join-item btn"
          :disabled="page === 1"
          @click="changePage(page - 1)"
        >
          «
        </button>
        <button class="join-item btn">Page {{ page }}</button>
        <button
          class="join-item btn"
          :disabled="page * pageSize >= total"
          @click="changePage(page + 1)"
        >
          »
        </button>
      </div>
    </div>

    <!-- 订单详情弹窗 -->
    <dialog :class="{ 'modal': true, 'modal-open': showDetailsModal }">
      <div class="modal-box w-11/12 max-w-2xl" v-if="selectedOrder">
        <h3 class="font-bold text-2xl mb-4">Order Details</h3>

        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <div class="text-sm opacity-70">Order Number</div>
              <div class="font-mono">{{ selectedOrder.order_number }}</div>
            </div>
            <div>
              <div class="text-sm opacity-70">Status</div>
              <span class="px-2 py-0.5 rounded text-xs font-medium" :class="getStatusBadgeClass(selectedOrder.status)">
                {{ selectedOrder.status }}
              </span>
            </div>
            <div>
              <div class="text-sm opacity-70">Domain</div>
              <div class="font-mono font-bold">{{ selectedOrder.full_domain }}</div>
            </div>
            <div>
              <div class="text-sm opacity-70">Duration</div>
              <div>{{ selectedOrder.is_lifetime ? 'Lifetime' : `${selectedOrder.years} Years` }}</div>
            </div>
          </div>

          <div class="divider"></div>

          <div>
            <div class="text-sm opacity-70 mb-2">Price Breakdown</div>
            <div class="space-y-2">
              <div class="flex justify-between">
                <span>Base Price:</span>
                <span class="font-mono">{{ formatPrice(selectedOrder.base_price) }}</span>
              </div>
              <div v-if="selectedOrder.discount_amount > 0" class="flex justify-between text-success">
                <span>Discount:</span>
                <span class="font-mono">-{{ formatPrice(selectedOrder.discount_amount) }}</span>
              </div>
              <div class="divider my-2"></div>
              <div class="flex justify-between text-lg font-bold">
                <span>Total:</span>
                <span class="font-mono text-primary">{{ formatPrice(selectedOrder.final_price) }}</span>
              </div>
            </div>
          </div>

          <div class="divider"></div>

          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <div class="opacity-70">Created At</div>
              <div>{{ formatDateTime(selectedOrder.created_at) }}</div>
            </div>
            <div v-if="selectedOrder.paid_at">
              <div class="opacity-70">Paid At</div>
              <div>{{ formatDateTime(selectedOrder.paid_at) }}</div>
            </div>
            <div v-if="selectedOrder.expires_at && selectedOrder.status === 'pending'">
              <div class="opacity-70">Expires At</div>
              <div>{{ formatDateTime(selectedOrder.expires_at) }}</div>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button class="btn btn-ghost" @click="showDetailsModal = false">Close</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="showDetailsModal = false">close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'

const router = useRouter()
const toast = useToast()
const { formatPrice } = useCurrency()

const loading = ref(true)
const paying = ref(null)
const cancelling = ref(null)
const orders = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const statusFilter = ref('')

const showDetailsModal = ref(false)
const selectedOrder = ref(null)

onMounted(async () => {
  await fetchOrders()
})

const fetchOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }

    const response = await axios.get('/api/orders', { params })
    orders.value = response.data.orders || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to fetch orders:', error)
    toast.error('Failed to load orders')
  } finally {
    loading.value = false
  }
}

const payOrder = async (order) => {
  paying.value = order.id
  try {
    const response = await axios.post(`/api/payments/${order.id}/initiate`)
    const redirectURL = response.data.redirect_url
    window.location.href = redirectURL
  } catch (error) {
    console.error('Failed to initiate payment:', error)
    toast.error(error.response?.data?.error || 'Failed to initiate payment')
    paying.value = null
  }
}

const cancelOrder = async (order) => {
  if (!confirm(`Are you sure you want to cancel order #${order.order_number}?`)) {
    return
  }

  cancelling.value = order.id
  try {
    await axios.post(`/api/orders/${order.id}/cancel`)
    toast.success('Order cancelled successfully')
    await fetchOrders()
  } catch (error) {
    console.error('Failed to cancel order:', error)
    toast.error(error.response?.data?.error || 'Failed to cancel order')
  } finally {
    cancelling.value = null
  }
}

const viewDetails = (order) => {
  selectedOrder.value = order
  showDetailsModal.value = true
}

const changePage = (newPage) => {
  page.value = newPage
  fetchOrders()
}

const getStatusBadgeClass = (status) => {
  const classes = {
    pending: 'bg-amber-100 text-amber-800',
    paid: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-800',
    refunded: 'bg-blue-100 text-blue-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString()
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString()
}

const isExpired = (order) => {
  if (!order.expires_at) return false
  const now = new Date()
  const expiresAt = new Date(order.expires_at)
  return now > expiresAt
}
</script>
