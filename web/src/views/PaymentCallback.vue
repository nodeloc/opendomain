<template>
  <section class="y-container flex min-h-[calc(100vh-4rem)] items-center justify-center py-14 sm:py-16">
    <div class="w-full max-w-lg">
      <!-- Processing -->
      <div v-if="processing" class="y-panel px-6 py-12 text-center sm:px-10">
        <YLoader :size="44" class="!py-0" />
        <p class="y-eyebrow mt-6">{{ $t('uiPublic.payment.eyebrow') }}</p>
        <h1 class="y-h2 mt-3">{{ $t('uiPublic.payment.processingTitle') }}</h1>
        <p class="mt-2 text-[15px] leading-6 text-y-muted">{{ $t('uiPublic.payment.processingDesc') }}</p>
      </div>

      <!-- Success -->
      <div v-else-if="success" class="y-panel p-6 sm:p-8">
        <StatusDot tone="ok" :label="$t('uiPublic.payment.successStatus')" />
        <h1 class="y-h2 mt-4">{{ $t('uiPublic.payment.successTitle') }}</h1>
        <p class="mt-2 text-[15px] leading-6 text-y-muted">{{ $t('uiPublic.payment.successDesc') }}</p>

        <div class="y-mint-soft mt-6 rounded-lg border px-5 py-4">
          <p class="y-label">{{ $t('uiPublic.payment.domain') }}</p>
          <p class="y-domain mt-1.5 text-lg">{{ orderInfo?.full_domain }}</p>
        </div>

        <dl class="mt-4 grid grid-cols-1 divide-y divide-y-line rounded-lg border border-y-line sm:grid-cols-2 sm:divide-x sm:divide-y-0">
          <div class="px-5 py-4">
            <dt class="y-label">{{ $t('uiPublic.payment.orderNumber') }}</dt>
            <dd class="mt-1.5 break-all font-mono text-sm text-y-text">{{ orderInfo?.order_number }}</dd>
          </div>
          <div class="px-5 py-4">
            <dt class="y-label">{{ $t('uiPublic.payment.amountPaid') }}</dt>
            <dd class="mt-1.5 font-mono text-sm text-y-text">{{ formatPrice(orderInfo?.final_price || 0) }}</dd>
          </div>
        </dl>

        <div class="mt-8 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button type="button" class="btn btn-outline" @click="goToOrders">
            {{ $t('uiPublic.payment.viewOrders') }}
          </button>
          <button type="button" class="btn btn-primary gap-1.5" @click="goToDomains">
            {{ $t('uiPublic.payment.viewDomains') }}
            <ArrowRight :size="16" :stroke-width="1.5" />
          </button>
        </div>
      </div>

      <!-- Failure -->
      <div v-else class="y-panel p-6 sm:p-8">
        <StatusDot tone="danger" :label="$t('uiPublic.payment.failedStatus')" />
        <h1 class="y-h2 mt-4">{{ $t('uiPublic.payment.failedTitle') }}</h1>

        <div class="y-danger-soft mt-6 flex items-start gap-3 rounded-lg border px-5 py-4">
          <CircleX :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-danger-text" />
          <p class="text-sm leading-6 text-y-danger-text">{{ errorMessage || $t('uiPublic.payment.failedDesc') }}</p>
        </div>

        <div v-if="orderInfo" class="mt-4 rounded-lg border border-y-line px-5 py-4">
          <p class="y-label">{{ $t('uiPublic.payment.orderNumber') }}</p>
          <p class="mt-1.5 break-all font-mono text-sm text-y-text">{{ orderInfo.order_number }}</p>
          <p v-if="orderInfo.full_domain" class="y-domain mt-1 text-sm text-y-muted">{{ orderInfo.full_domain }}</p>
          <p class="mt-3 text-sm text-y-muted">{{ $t('uiPublic.payment.retryHint') }}</p>
        </div>

        <div class="mt-8 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button type="button" class="btn btn-ghost" @click="goToHome">
            {{ $t('uiPublic.backHome') }}
          </button>
          <button type="button" class="btn btn-primary gap-1.5" @click="goToOrders">
            {{ $t('uiPublic.payment.viewOrders') }}
            <ArrowRight :size="16" :stroke-width="1.5" />
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, CircleX } from 'lucide-vue-next'
import axios from '../utils/axios'
import { useCurrency } from '../composables/useCurrency'
import YLoader from '../components/brand/YLoader.vue'
import StatusDot from '../components/ui/StatusDot.vue'

const route = useRoute()
const router = useRouter()
const { formatPrice } = useCurrency()

const processing = ref(true)
const success = ref(false)
const errorMessage = ref('')
const orderInfo = ref(null)
const retryCount = ref(0)
const maxRetries = 5

onMounted(async () => {
  // 检查路由参数
  const orderId = route.query.order_id

  if (!orderId) {
    errorMessage.value = 'Invalid payment callback'
    processing.value = false
    return
  }

  // 延迟一下，给后端回调处理时间
  await new Promise(resolve => setTimeout(resolve, 2000))

  // 查询订单状态
  await checkOrderStatus(orderId)
})

const checkOrderStatus = async (orderId) => {
  try {
    const response = await axios.get(`/api/orders/${orderId}`)
    orderInfo.value = response.data

    if (orderInfo.value.status === 'paid') {
      success.value = true
    } else if (orderInfo.value.status === 'pending' && retryCount.value < maxRetries) {
      // 订单还在处理中，等待一下再查询
      retryCount.value++
      await new Promise(resolve => setTimeout(resolve, 2000))
      await checkOrderStatus(orderId)
      return
    } else if (orderInfo.value.status === 'pending' && retryCount.value >= maxRetries) {
      // 重试次数超限
      success.value = false
      errorMessage.value = 'Payment is still processing. Please check your order status later.'
    } else {
      success.value = false
      errorMessage.value = `Payment was ${orderInfo.value.status}`
    }
  } catch (error) {
    console.error('Failed to check order status:', error)
    success.value = false
    errorMessage.value = error.response?.data?.error || 'Failed to verify payment status'
  } finally {
    processing.value = false
  }
}

const goToDomains = () => {
  router.push('/domains')
}

const goToOrders = () => {
  router.push('/orders')
}

const goToHome = () => {
  router.push('/')
}
</script>
