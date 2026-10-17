<template>
  <div class="y-page">
    <PageHeader :eyebrow="$t('uiUser.orders.eyebrow')" :title="$t('uiUser.orders.title')" :description="$t('uiUser.orders.description')" />

    <!-- Status filter -->
    <div class="-mx-5 overflow-x-auto px-5 sm:mx-0 sm:px-0">
      <div role="tablist" class="tabs tabs-boxed inline-flex w-max">
        <button
          v-for="f in statusFilters"
          :key="f.value || 'all'"
          type="button"
          role="tab"
          class="tab whitespace-nowrap"
          :class="{ 'tab-active': statusFilter === f.value }"
          :aria-selected="statusFilter === f.value"
          @click="statusFilter = f.value; fetchOrders()"
        >
          {{ f.label }}
        </button>
      </div>
    </div>

    <section class="y-panel">
      <YLoader v-if="loading" :label="$t('uiUser.common.loading')" />

      <EmptyState v-else-if="orders.length === 0" :title="$t('uiUser.orders.emptyTitle')" :description="$t('uiUser.orders.emptyDesc')">
        <button type="button" class="btn btn-primary" @click="$router.push('/')">{{ $t('uiUser.orders.claim') }}</button>
      </EmptyState>

      <template v-else>
        <!-- Desktop table -->
        <div class="hidden overflow-x-auto md:block">
          <table class="y-table">
            <thead>
              <tr>
                <th>{{ $t('uiUser.orders.domain') }}</th>
                <th>{{ $t('uiUser.orders.number') }}</th>
                <th>{{ $t('uiUser.orders.status') }}</th>
                <th>{{ $t('uiUser.orders.term') }}</th>
                <th>{{ $t('uiUser.orders.created') }}</th>
                <th class="text-right">{{ $t('uiUser.orders.total') }}</th>
                <th class="text-right"><span class="sr-only">{{ $t('uiUser.orders.actions') }}</span></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in orders" :key="order.id">
                <td class="whitespace-nowrap font-mono text-[13px] font-medium text-y-text">{{ order.full_domain }}</td>
                <td class="whitespace-nowrap font-mono text-[13px] text-y-muted">{{ order.order_number }}</td>
                <td class="whitespace-nowrap">
                  <span :class="getStatusBadgeClass(order.status)">{{ $t('order.' + order.status) }}</span>
                  <p v-if="order.expires_at && order.status === 'pending'" class="mt-1 font-mono text-[11px] text-y-faint">
                    {{ $t('uiUser.orders.expiresOn', { date: formatDate(order.expires_at) }) }}
                  </p>
                </td>
                <td class="whitespace-nowrap text-y-muted">{{ termLabel(order) }}</td>
                <td class="whitespace-nowrap font-mono text-[13px] text-y-muted">{{ formatDate(order.created_at) }}</td>
                <td class="whitespace-nowrap text-right">
                  <span class="font-mono font-medium text-y-text">{{ formatPrice(order.final_price) }}</span>
                  <p v-if="order.discount_amount > 0" class="mt-0.5 font-mono text-[11px] text-y-mint-text">−{{ formatPrice(order.discount_amount) }}</p>
                </td>
                <td>
                  <div class="flex justify-end gap-1 whitespace-nowrap">
                    <button
                      v-if="order.status === 'pending' && !isExpired(order)"
                      type="button"
                      class="btn btn-primary btn-sm"
                      :disabled="paying === order.id"
                      @click="payOrder(order)"
                    >
                      <span v-if="paying === order.id" class="loading loading-spinner loading-xs"></span>
                      <span v-else>{{ $t('uiUser.orders.pay') }}</span>
                    </button>
                    <button
                      v-if="order.status === 'pending' || order.status === 'expired'"
                      type="button"
                      class="btn btn-ghost btn-sm text-y-danger-text"
                      :disabled="cancelling === order.id"
                      @click="cancelOrder(order)"
                    >
                      <span v-if="cancelling === order.id" class="loading loading-spinner loading-xs"></span>
                      <span v-else>{{ $t('uiUser.orders.cancel') }}</span>
                    </button>
                    <button type="button" class="btn btn-ghost btn-sm btn-square" :title="$t('uiUser.orders.details')" :aria-label="$t('uiUser.orders.details')" @click="viewDetails(order)">
                      <ChevronRight :size="16" :stroke-width="1.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Mobile list -->
        <ul class="divide-y divide-y-line md:hidden">
          <li v-for="order in orders" :key="order.id" class="space-y-3 px-5 py-4">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <p class="y-domain">{{ order.full_domain }}</p>
                <p class="mt-1 font-mono text-xs text-y-faint">{{ order.order_number }}</p>
              </div>
              <span class="shrink-0" :class="getStatusBadgeClass(order.status)">{{ $t('order.' + order.status) }}</span>
            </div>
            <div class="flex items-end justify-between gap-3">
              <div class="space-y-0.5 text-xs text-y-muted">
                <p>{{ termLabel(order) }} · <span class="font-mono">{{ formatDate(order.created_at) }}</span></p>
                <p v-if="order.expires_at && order.status === 'pending'" class="font-mono text-y-faint">
                  {{ $t('uiUser.orders.expiresOn', { date: formatDate(order.expires_at) }) }}
                </p>
              </div>
              <div class="text-right">
                <p class="font-mono text-lg font-medium text-y-text">{{ formatPrice(order.final_price) }}</p>
                <p v-if="order.discount_amount > 0" class="font-mono text-[11px] text-y-mint-text">−{{ formatPrice(order.discount_amount) }}</p>
              </div>
            </div>
            <div class="flex flex-wrap gap-2">
              <button
                v-if="order.status === 'pending' && !isExpired(order)"
                type="button"
                class="btn btn-primary btn-sm"
                :disabled="paying === order.id"
                @click="payOrder(order)"
              >
                <span v-if="paying === order.id" class="loading loading-spinner loading-xs"></span>
                <span v-else>{{ $t('uiUser.orders.pay') }}</span>
              </button>
              <button
                v-if="order.status === 'pending' || order.status === 'expired'"
                type="button"
                class="btn btn-ghost btn-sm text-y-danger-text"
                :disabled="cancelling === order.id"
                @click="cancelOrder(order)"
              >
                <span v-if="cancelling === order.id" class="loading loading-spinner loading-xs"></span>
                <span v-else>{{ $t('uiUser.orders.cancel') }}</span>
              </button>
              <button type="button" class="btn btn-ghost btn-sm ml-auto" @click="viewDetails(order)">{{ $t('uiUser.orders.details') }}</button>
            </div>
          </li>
        </ul>
      </template>
    </section>

    <!-- Pagination -->
    <div v-if="total > pageSize" class="flex items-center justify-center gap-2">
      <button type="button" class="btn btn-outline btn-sm btn-square" :disabled="page === 1" :aria-label="$t('common.previous')" @click="changePage(page - 1)">
        <ChevronLeft :size="16" :stroke-width="1.5" />
      </button>
      <span class="px-3 font-mono text-sm text-y-muted">{{ $t('uiUser.orders.page', { page }) }}</span>
      <button type="button" class="btn btn-outline btn-sm btn-square" :disabled="page * pageSize >= total" :aria-label="$t('common.next')" @click="changePage(page + 1)">
        <ChevronRight :size="16" :stroke-width="1.5" />
      </button>
    </div>

    <!-- Order details -->
    <dialog :class="{ 'modal': true, 'modal-open': showDetailsModal }">
      <div v-if="selectedOrder" class="modal-box w-11/12 max-w-xl">
        <div class="flex items-start justify-between gap-4">
          <div class="min-w-0">
            <h3 class="text-lg font-semibold text-y-text">{{ $t('uiUser.orders.detailsTitle') }}</h3>
            <p class="mt-1 font-mono text-sm text-y-muted">{{ selectedOrder.order_number }}</p>
          </div>
          <span class="shrink-0" :class="getStatusBadgeClass(selectedOrder.status)">{{ $t('order.' + selectedOrder.status) }}</span>
        </div>

        <dl class="mt-6 grid grid-cols-2 gap-x-6 gap-y-4">
          <div class="col-span-2">
            <dt class="y-label">{{ $t('uiUser.orders.domain') }}</dt>
            <dd class="y-domain mt-1.5">{{ selectedOrder.full_domain }}</dd>
          </div>
          <div>
            <dt class="y-label">{{ $t('uiUser.orders.term') }}</dt>
            <dd class="mt-1.5 text-sm text-y-text">{{ termLabel(selectedOrder) }}</dd>
          </div>
          <div>
            <dt class="y-label">{{ $t('uiUser.orders.created') }}</dt>
            <dd class="mt-1.5 font-mono text-[13px] text-y-text">{{ formatDateTime(selectedOrder.created_at) }}</dd>
          </div>
          <div v-if="selectedOrder.paid_at">
            <dt class="y-label">{{ $t('uiUser.orders.paidAt') }}</dt>
            <dd class="mt-1.5 font-mono text-[13px] text-y-text">{{ formatDateTime(selectedOrder.paid_at) }}</dd>
          </div>
          <div v-if="selectedOrder.expires_at && selectedOrder.status === 'pending'">
            <dt class="y-label">{{ $t('uiUser.orders.expiresAt') }}</dt>
            <dd class="mt-1.5 font-mono text-[13px] text-y-text">{{ formatDateTime(selectedOrder.expires_at) }}</dd>
          </div>
        </dl>

        <div class="y-inset mt-6 space-y-2 px-4 py-4 text-sm">
          <div class="flex justify-between gap-4">
            <span class="text-y-muted">{{ $t('uiUser.orders.basePrice') }}</span>
            <span class="font-mono text-y-text">{{ formatPrice(selectedOrder.base_price) }}</span>
          </div>
          <div v-if="selectedOrder.discount_amount > 0" class="flex justify-between gap-4">
            <span class="text-y-muted">{{ $t('uiUser.orders.discount') }}</span>
            <span class="font-mono text-y-mint-text">−{{ formatPrice(selectedOrder.discount_amount) }}</span>
          </div>
          <div class="y-divider !my-3"></div>
          <div class="flex items-baseline justify-between gap-4">
            <span class="font-medium text-y-text">{{ $t('uiUser.orders.total') }}</span>
            <span class="font-mono text-xl font-medium text-y-text">{{ formatPrice(selectedOrder.final_price) }}</span>
          </div>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDetailsModal = false">{{ $t('common.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="showDetailsModal = false">close</button>
      </form>
    </dialog>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { useI18n } from 'vue-i18n'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const router = useRouter()
const toast = useToast()
const { formatPrice } = useCurrency()
const { t } = useI18n()

const statusFilters = computed(() => [
  { value: '', label: t('order.all') },
  { value: 'pending', label: t('order.pending') },
  { value: 'paid', label: t('order.paid') },
  { value: 'cancelled', label: t('order.cancelled') },
  { value: 'expired', label: t('order.expired') },
])

const termLabel = (order) => (order.is_lifetime ? t('uiUser.orders.lifetime') : t('uiUser.orders.years', order.years))

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
    pending: 'y-chip-warn',
    paid: 'y-chip-mint',
    cancelled: 'y-chip',
    expired: 'y-chip',
    refunded: 'y-chip-info'
  }
  return classes[status] || 'y-chip'
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
