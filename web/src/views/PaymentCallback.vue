<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-12 max-w-2xl">
    <!-- 处理中状态 -->
    <div v-if="processing" class="card bg-base-100 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <span class="loading loading-spinner loading-lg text-primary mb-6"></span>
        <h2 class="card-title text-2xl mb-2">Processing Payment...</h2>
        <p class="opacity-70">Please wait while we verify your payment</p>
      </div>
    </div>

    <!-- 成功状态 -->
    <div v-else-if="success" class="card bg-base-100 shadow-xl border-2 border-success">
      <div class="card-body items-center text-center py-12">
        <div class="w-24 h-24 rounded-full bg-success/10 flex items-center justify-center mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>

        <h2 class="card-title text-3xl mb-4 text-success">Payment Successful!</h2>

        <p class="text-lg mb-2">Your domain has been registered successfully</p>
        <p class="text-xl font-mono font-bold mb-6">{{ orderInfo?.full_domain }}</p>

        <div class="stats shadow mb-6">
          <div class="stat">
            <div class="stat-title">Order Number</div>
            <div class="stat-value text-lg font-mono">{{ orderInfo?.order_number }}</div>
          </div>
          <div class="stat">
            <div class="stat-title">Amount Paid</div>
            <div class="stat-value text-lg text-primary">{{ formatPrice(orderInfo?.final_price || 0) }}</div>
          </div>
        </div>

        <div class="flex gap-4 w-full max-w-md">
          <button class="btn btn-ghost flex-1" @click="goToOrders">
            View Orders
          </button>
          <button class="btn btn-primary flex-1" @click="goToDomains">
            View My Domains
          </button>
        </div>
      </div>
    </div>

    <!-- 失败状态 -->
    <div v-else class="card bg-base-100 shadow-xl border-2 border-error">
      <div class="card-body items-center text-center py-12">
        <div class="w-24 h-24 rounded-full bg-error/10 flex items-center justify-center mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>

        <h2 class="card-title text-3xl mb-4 text-error">Payment Failed</h2>

        <p class="text-lg mb-6 opacity-70">{{ errorMessage || 'Your payment could not be processed' }}</p>

        <div v-if="orderInfo" class="alert alert-warning mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <div>
            <div class="font-bold">Order #{{ orderInfo.order_number }}</div>
            <div class="text-sm">You can retry payment from your orders page</div>
          </div>
        </div>

        <div class="flex gap-4 w-full max-w-md">
          <button class="btn btn-ghost flex-1" @click="goToHome">
            Back to Home
          </button>
          <button class="btn btn-primary flex-1" @click="goToOrders">
            View Orders
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../utils/axios'
import { useCurrency } from '../composables/useCurrency'

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
