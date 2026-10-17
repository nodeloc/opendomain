<template>
  <div class="y-page">
    <YLoader v-if="loading" :label="$t('uiUser.common.loading')" />

    <template v-else>
      <div class="space-y-4">
        <button type="button" class="inline-flex items-center gap-1.5 text-sm text-y-muted transition-colors hover:text-y-text" @click="goBack">
          <ArrowLeft :size="16" :stroke-width="1.5" />
          {{ $t('uiUser.checkout.back') }}
        </button>
        <PageHeader :eyebrow="$t('uiUser.checkout.eyebrow')" :title="$t('uiUser.checkout.title')" :description="$t('uiUser.checkout.description')" />
      </div>

      <div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_360px] lg:items-start">
        <!-- Configuration -->
        <div class="space-y-6">
          <section class="y-panel flex items-center gap-4 px-5 py-5">
            <span class="y-icon-tile"><Globe :size="20" :stroke-width="1.5" /></span>
            <div class="min-w-0">
              <p class="y-domain text-xl">{{ fullDomain }}</p>
              <p class="mt-1 text-sm text-y-muted">{{ rootDomain?.description || $t('uiUser.checkout.premium') }}</p>
            </div>
          </section>

          <section class="y-panel">
            <div class="y-panel-header">
              <h2 class="y-h3">{{ $t('uiUser.checkout.termTitle') }}</h2>
              <span class="font-mono text-sm text-y-muted">
                {{ isLifetime ? $t('uiUser.checkout.lifetime') : $t('uiUser.checkout.yearsCount', years) }}
              </span>
            </div>
            <div class="space-y-5 px-5 py-5">
              <div class="y-field">
                <span class="y-field-label">{{ $t('uiUser.checkout.yearsLabel') }}</span>
                <div class="grid grid-cols-5 gap-2 sm:grid-cols-10" role="radiogroup" :aria-label="$t('uiUser.checkout.yearsLabel')">
                  <button
                    v-for="n in 10"
                    :key="n"
                    type="button"
                    role="radio"
                    :aria-checked="!isLifetime && years === n"
                    class="h-10 rounded-sm border font-mono text-sm transition-colors disabled:cursor-not-allowed disabled:opacity-40"
                    :class="!isLifetime && years === n
                      ? 'y-mint-soft text-y-mint-text'
                      : 'border-y-line bg-y-bg text-y-muted hover:border-y-line-strong hover:text-y-text'"
                    :disabled="isLifetime || calculating"
                    @click="years = n; calculatePrice()"
                  >
                    {{ n }}
                  </button>
                </div>
                <span class="y-field-hint">{{ $t('uiUser.checkout.yearsHint') }}</span>
              </div>

              <label v-if="rootDomain?.lifetime_price" class="y-inset flex cursor-pointer items-start gap-3 px-4 py-3">
                <input v-model="isLifetime" type="checkbox" class="checkbox checkbox-sm mt-0.5" @change="calculatePrice" />
                <span>
                  <span class="block text-sm font-medium text-y-text">{{ $t('uiUser.checkout.lifetimeTitle') }}</span>
                  <span class="block text-xs text-y-muted">{{ $t('uiUser.checkout.lifetimeHint') }}</span>
                </span>
              </label>
            </div>
          </section>

          <section class="y-panel">
            <div class="y-panel-header">
              <h2 class="y-h3">{{ $t('uiUser.checkout.couponTitle') }}</h2>
              <span class="text-xs text-y-faint">{{ $t('uiUser.checkout.optional') }}</span>
            </div>
            <div class="space-y-3 px-5 py-5">
              <div class="flex flex-col gap-2 sm:flex-row">
                <input
                  v-model="couponCode"
                  type="text"
                  :placeholder="$t('uiUser.checkout.couponPlaceholder')"
                  class="input input-bordered w-full min-w-0 font-mono sm:flex-1 uppercase placeholder:normal-case"
                  autocomplete="off"
                  autocapitalize="characters"
                  spellcheck="false"
                  :aria-label="$t('uiUser.checkout.couponTitle')"
                  @keyup.enter="applyCoupon"
                />
                <button type="button" class="btn btn-outline px-6" :disabled="calculating" @click="applyCoupon">
                  <span v-if="calculating" class="loading loading-spinner loading-sm"></span>
                  <span v-else>{{ $t('uiUser.checkout.apply') }}</span>
                </button>
              </div>

              <div v-if="priceInfo?.coupon_applied" class="alert alert-success text-sm">
                <CheckCircle2 :size="18" :stroke-width="1.5" class="text-y-mint-text" />
                <span>{{ $t('uiUser.checkout.couponApplied', { amount: formatPrice(priceInfo.discount_amount) }) }}</span>
              </div>

              <div v-if="couponError" class="alert alert-error text-sm">
                <XCircle :size="18" :stroke-width="1.5" />
                <span>{{ couponError }}</span>
              </div>
            </div>
          </section>
        </div>

        <!-- Summary -->
        <aside class="y-panel lg:sticky lg:top-20">
          <div class="y-panel-header">
            <h2 class="y-h3">{{ $t('uiUser.checkout.summary') }}</h2>
            <span v-if="calculating" class="loading loading-spinner loading-xs"></span>
          </div>
          <div class="space-y-3 px-5 py-5 text-sm">
            <div class="flex justify-between gap-4">
              <span class="text-y-muted">{{ $t('uiUser.checkout.name') }}</span>
              <span class="min-w-0 truncate text-right font-mono text-y-text" :title="fullDomain">{{ fullDomain }}</span>
            </div>
            <div class="flex justify-between gap-4">
              <span class="text-y-muted">{{ $t('uiUser.checkout.term') }}</span>
              <span class="font-mono text-y-text">{{ isLifetime ? $t('uiUser.checkout.lifetime') : $t('uiUser.checkout.yearsCount', years) }}</span>
            </div>
            <div class="flex justify-between gap-4">
              <span class="text-y-muted">{{ $t('uiUser.checkout.basePrice') }}</span>
              <span class="font-mono text-y-text">{{ formatPrice(priceInfo?.base_price || 0) }}</span>
            </div>
            <div v-if="priceInfo?.discount_amount > 0" class="flex justify-between gap-4">
              <span class="text-y-muted">{{ $t('uiUser.checkout.discount') }}</span>
              <span class="font-mono text-y-mint-text">−{{ formatPrice(priceInfo.discount_amount) }}</span>
            </div>

            <div class="y-divider !my-4"></div>

            <div class="flex items-baseline justify-between gap-4">
              <span class="font-medium text-y-text">{{ $t('uiUser.checkout.total') }}</span>
              <span class="font-mono text-3xl font-medium tracking-tight text-y-text">{{ formatPrice(priceInfo?.final_price || 0) }}</span>
            </div>
            <p class="text-right text-xs text-y-faint">
              <template v-if="isLifetime">{{ $t('uiUser.checkout.oneTime') }}</template>
              <template v-else>{{ $t('uiUser.checkout.perYear', { price: formatPrice((priceInfo?.final_price || 0) / years) }) }}</template>
            </p>
          </div>
          <div class="space-y-2 border-t border-y-line px-5 py-5">
            <button type="button" class="btn btn-primary w-full" :disabled="creating || !priceInfo" @click="createOrderAndPay">
              <span v-if="creating" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ $t('uiUser.checkout.pay') }}</span>
            </button>
            <button type="button" class="btn btn-ghost w-full" @click="goBack">{{ $t('common.cancel') }}</button>
          </div>
        </aside>
      </div>
    </template>
  </div>
</template>
<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { ArrowLeft, Globe, CheckCircle2, XCircle } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import YLoader from '../components/brand/YLoader.vue'

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

// 辅助函数：将错误消息中的 RFC3339 时间转换为本地时间
const convertTimeInErrorMessage = (message) => {
  if (!message) return message
  
  // 匹配 RFC3339 格式的时间字符串，如 "2026-02-17T00:00:00Z" 或 "2026-02-17T00:00:00+08:00"
  const rfc3339Pattern = /(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2}))/g
  
  return message.replace(rfc3339Pattern, (match) => {
    try {
      const date = new Date(match)
      return date.toLocaleString()
    } catch (e) {
      return match
    }
  })
}

const calculatePrice = async () => {
  calculating.value = true
  couponError.value = ''

  try {
    const payload = {
      root_domain_id: rootDomainId.value,
      years: isLifetime.value ? 0 : years.value,
      is_lifetime: isLifetime.value,
    }

    if (couponCode.value && couponCode.value.trim() !== '') {
      payload.coupon_code = couponCode.value.trim()
    }

    const response = await axios.post('/api/orders/calculate', payload)
    priceInfo.value = response.data

    // 显示优惠券错误信息（如果有）
    if (response.data.coupon_error) {
      couponError.value = convertTimeInErrorMessage(response.data.coupon_error)
    } else if (couponCode.value && !response.data.coupon_applied) {
      couponError.value = 'Coupon could not be applied'
    }
  } catch (error) {
    console.error('Failed to calculate price:', error)
    let errorMsg = error.response?.data?.error || 'Failed to calculate price'
    
    // 尝试从错误信息中提取时间并转换为本地时间
    errorMsg = convertTimeInErrorMessage(errorMsg)
    couponError.value = errorMsg

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
      years: isLifetime.value ? 0 : years.value,
      is_lifetime: isLifetime.value,
    }

    if (couponCode.value && priceInfo.value?.coupon_applied) {
      orderPayload.coupon_code = couponCode.value.trim()
    }

    const orderResponse = await axios.post('/api/orders', orderPayload)
    const order = orderResponse.data.order

    // 如果总价接近 0（< 0.01），直接标记为免费订单
    // 使用 < 0.01 而不是 === 0 来处理浮点精度问题和小额折扣
    if (order.final_price < 0.01) {
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
    const errorData = error.response?.data
    if (error.response?.status === 403 && errorData?.required_level) {
      toast.error(`${errorData.error} (required: ${errorData.required_level}, yours: ${errorData.your_level})`)
    } else {
      toast.error(errorData?.error || 'Failed to create order')
    }
  } finally {
    creating.value = false
  }
}

const goBack = () => {
  router.push('/')
}
</script>
