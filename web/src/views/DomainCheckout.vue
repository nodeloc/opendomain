<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-4xl">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- 主内容 -->
    <div v-else class="space-y-6">
      <!-- 标题 -->
      <div>
        <h1 class="text-3xl font-bold">{{ $t('checkout.title') }}</h1>
        <p class="text-lg opacity-70 mt-2">{{ $t('checkout.subtitle') }}</p>
      </div>

      <!-- 域名信息卡片 -->
      <div class="card bg-base-100 shadow-xl border border-base-300">
        <div class="card-body">
          <h2 class="card-title text-2xl mb-4">{{ $t('checkout.domainInfo') }}</h2>

          <div class="flex items-center gap-3 p-4 bg-base-200 rounded-lg">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
            <div>
              <div class="text-2xl font-mono font-bold">{{ fullDomain }}</div>
              <div class="text-sm opacity-70 mt-1">{{ rootDomain?.description || 'Premium Domain' }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 注册期限选择 -->
      <div class="card bg-base-100 shadow-xl border border-base-300">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">{{ $t('checkout.registrationPeriod') }}</h2>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">{{ $t('checkout.duration') }}</span>
            </label>
            <div class="flex items-center gap-4">
              <input
                v-model.number="years"
                type="range"
                min="1"
                max="10"
                class="range range-primary flex-1"
                :disabled="isLifetime"
                @input="calculatePrice"
              />
              <div class="badge badge-primary badge-lg font-mono">{{ years }} {{ years === 1 ? $t('checkout.year') : $t('checkout.years') }}</div>
            </div>
          </div>

          <div v-if="rootDomain?.lifetime_price" class="form-control mt-4">
            <label class="label cursor-pointer justify-start gap-3">
              <input
                v-model="isLifetime"
                type="checkbox"
                class="checkbox checkbox-secondary"
                @change="calculatePrice"
              />
              <span class="label-text">
                <span class="font-semibold">{{ $t('checkout.lifetimeRegistration') }}</span>
                <span class="ml-2 text-sm opacity-70">{{ $t('checkout.lifetimeDesc') }}</span>
              </span>
            </label>
          </div>
        </div>
      </div>

      <!-- 优惠券 -->
      <div class="card bg-base-100 shadow-xl border border-base-300">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">{{ $t('checkout.couponCode') }}</h2>

          <div class="join w-full">
            <input
              v-model="couponCode"
              type="text"
              :placeholder="$t('checkout.enterCoupon')"
              class="input input-bordered join-item flex-1"
              @keyup.enter="applyCoupon"
            />
            <button
              class="btn btn-primary join-item"
              @click="applyCoupon"
              :disabled="calculating"
            >
              <span v-if="calculating" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ $t('coupon.apply') }}</span>
            </button>
          </div>

          <div v-if="priceInfo?.coupon_applied" class="alert alert-success mt-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Coupon applied! You saved {{ formatPrice(priceInfo.discount_amount) }}</span>
          </div>

          <div v-if="couponError" class="alert alert-error mt-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ couponError }}</span>
          </div>
        </div>
      </div>

      <!-- 价格明细 -->
      <div class="card bg-base-100 shadow-xl border border-base-300">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">Price Summary</h2>

          <div class="space-y-3">
            <div class="flex justify-between items-center pb-2">
              <span class="opacity-70">Base Price:</span>
              <span class="font-mono">{{ formatPrice(priceInfo?.base_price || 0) }}</span>
            </div>

            <div v-if="priceInfo?.discount_amount > 0" class="flex justify-between items-center pb-2 text-success">
              <span>Discount:</span>
              <span class="font-mono">-{{ formatPrice(priceInfo.discount_amount) }}</span>
            </div>

            <div class="divider my-2"></div>

            <div class="flex justify-between items-center text-2xl font-bold">
              <span>Total:</span>
              <span class="font-mono text-primary">{{ formatPrice(priceInfo?.final_price || 0) }}</span>
            </div>

            <div v-if="isLifetime" class="text-sm opacity-70 text-center">
              One-time payment for permanent ownership
            </div>
            <div v-else class="text-sm opacity-70 text-center">
              {{ formatPrice((priceInfo?.final_price || 0) / years) }}/year
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex gap-4">
        <button class="btn btn-ghost flex-1" @click="goBack">
          Cancel
        </button>
        <button
          class="btn btn-primary flex-1"
          @click="createOrderAndPay"
          :disabled="creating || !priceInfo"
        >
          <span v-if="creating" class="loading loading-spinner loading-sm"></span>
          <span v-else>Create Order & Pay</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const { formatPrice } = useCurrency()

const loading = ref(true)
const calculating = ref(false)
const creating = ref(false)

const subdomain = ref('')
const rootDomainId = ref(null)
const rootDomain = ref(null)
const years = ref(1)
const isLifetime = ref(false)
const couponCode = ref('')
const priceInfo = ref(null)
const couponError = ref('')

const fullDomain = computed(() => {
  return rootDomain.value ? `${subdomain.value}.${rootDomain.value.domain}` : ''
})

onMounted(async () => {
  // 从路由获取参数
  subdomain.value = route.query.subdomain || ''
  rootDomainId.value = parseInt(route.query.root_domain_id || '0')

  if (!subdomain.value || !rootDomainId.value) {
    toast.error('Invalid domain information')
    router.push('/')
    return
  }

  await fetchRootDomain()
  await calculatePrice()
  loading.value = false
})

const fetchRootDomain = async () => {
  try {
    const response = await axios.get('/api/public/root-domains')
    const domains = response.data.root_domains || []
    rootDomain.value = domains.find(d => d.id === rootDomainId.value)

    if (!rootDomain.value) {
      throw new Error('Root domain not found')
    }
  } catch (error) {
    console.error('Failed to fetch root domain:', error)
    toast.error('Failed to load domain information')
    router.push('/')
  }
}

const calculatePrice = async () => {
  calculating.value = true
  couponError.value = ''

  try {
    const payload = {
      root_domain_id: rootDomainId.value,
      years: years.value,
      is_lifetime: isLifetime.value,
    }

    if (couponCode.value && couponCode.value.trim() !== '') {
      payload.coupon_code = couponCode.value.trim()
    }

    const response = await axios.post('/api/orders/calculate', payload)
    priceInfo.value = response.data

    // 清除优惠券输入如果没有应用成功
    if (couponCode.value && !response.data.coupon_applied) {
      couponError.value = 'Coupon could not be applied'
    }
  } catch (error) {
    console.error('Failed to calculate price:', error)
    couponError.value = error.response?.data?.error || 'Failed to calculate price'

    // 如果有优惠券错误，重新计算不带优惠券的价格
    if (couponCode.value) {
      couponCode.value = ''
      await calculatePrice()
    }
  } finally {
    calculating.value = false
  }
}

const applyCoupon = async () => {
  if (!couponCode.value || couponCode.value.trim() === '') {
    couponError.value = 'Please enter a coupon code'
    return
  }
  await calculatePrice()
}

const createOrderAndPay = async () => {
  creating.value = true
  try {
    // 创建订单
    const orderPayload = {
      subdomain: subdomain.value,
      root_domain_id: rootDomainId.value,
      years: years.value,
      is_lifetime: isLifetime.value,
    }

    if (couponCode.value && priceInfo.value?.coupon_applied) {
      orderPayload.coupon_code = couponCode.value.trim()
    }

    const orderResponse = await axios.post('/api/orders', orderPayload)
    const order = orderResponse.data.order

    // 如果总价为 0，直接跳转到成功页面（不需要支付）
    if (order.final_price === 0) {
      // 直接标记订单为已支付
      await axios.post(`/api/payments/${order.id}/complete-free`)
      // 跳转到支付成功页面
      window.location.href = `/payment/success?order_id=${order.id}`
      return
    }

    // 发起支付
    const paymentResponse = await axios.post(`/api/payments/${order.id}/initiate`)
    const redirectURL = paymentResponse.data.redirect_url

    // 跳转到支付页面
    window.location.href = redirectURL
  } catch (error) {
    console.error('Failed to create order:', error)
    toast.error(error.response?.data?.error || 'Failed to create order')
  } finally {
    creating.value = false
  }
}

const goBack = () => {
  router.push('/')
}
</script>
