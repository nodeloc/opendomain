<template>
  <div class="y-page">
    <PageHeader :eyebrow="$t('uiUser.coupons.eyebrow')" :title="$t('uiUser.coupons.title')" :description="$t('uiUser.coupons.description')" />

    <!-- Redeem -->
    <section class="y-panel">
      <div class="y-panel-header">
        <h2 class="y-h3">{{ $t('uiUser.coupons.redeemTitle') }}</h2>
      </div>
      <form class="space-y-4 px-5 py-5" @submit.prevent="applyCoupon">
        <div class="y-field">
          <label for="coupon-code" class="y-field-label">{{ $t('uiUser.coupons.codeLabel') }}</label>
          <div class="flex flex-col gap-2 sm:flex-row">
            <label class="input input-bordered flex min-w-0 items-center gap-3 sm:flex-1">
              <Ticket :size="18" :stroke-width="1.5" class="shrink-0 text-y-faint" />
              <input
                id="coupon-code"
                v-model="couponCode"
                type="text"
                :placeholder="$t('coupon.placeholder')"
                class="min-w-0 grow bg-transparent font-mono uppercase tracking-wide outline-none placeholder:normal-case placeholder:tracking-normal focus-visible:outline-none"
                autocomplete="off"
                autocapitalize="characters"
                spellcheck="false"
                :disabled="applying"
              />
            </label>
            <button type="submit" class="btn btn-primary gap-1.5 px-6" :disabled="!couponCode || applying">
              <span v-if="applying" class="loading loading-spinner loading-sm"></span>
              <Check v-else :size="16" :stroke-width="1.5" />
              {{ $t('uiUser.coupons.redeem') }}
            </button>
          </div>
          <span class="y-field-hint">{{ $t('uiUser.coupons.caseHint') }}</span>
        </div>

        <div class="y-inset flex items-start gap-3 px-4 py-3 text-sm text-y-muted">
          <Info :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-info" />
          <span>{{ $t('coupon.infoText') }}</span>
        </div>
      </form>
    </section>

    <!-- History -->
    <section class="y-panel">
      <div class="y-panel-header">
        <h2 class="y-h3">{{ $t('uiUser.coupons.historyTitle') }}</h2>
        <span v-if="usageHistory.length > 0" class="y-chip font-mono">{{ $t('uiUser.coupons.usedCount', usageHistory.length) }}</span>
      </div>

      <YLoader v-if="loadingHistory" :label="$t('uiUser.common.loading')" />

      <EmptyState v-else-if="usageHistory.length === 0" :title="$t('uiUser.coupons.emptyTitle')" :description="$t('uiUser.coupons.emptyDesc')" />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('coupon.code') }}</th>
              <th>{{ $t('coupon.type') }}</th>
              <th>{{ $t('coupon.benefit') }}</th>
              <th>{{ $t('coupon.usedAt') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="usage in usageHistory" :key="usage.id">
              <td class="whitespace-nowrap font-mono text-[13px] font-medium">{{ usage.coupon?.code }}</td>
              <td class="whitespace-nowrap">
                <span :class="getBadgeClass(usage.coupon?.discount_type)">{{ formatType(usage.coupon?.discount_type) }}</span>
              </td>
              <td class="text-y-text">{{ usage.benefit_applied }}</td>
              <td class="whitespace-nowrap font-mono text-[13px] text-y-muted">{{ formatDate(usage.used_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Ticket, Check, Info } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const couponCode = ref('')
const applying = ref(false)
const loadingHistory = ref(true)
const usageHistory = ref([])

onMounted(async () => {
  await fetchUsageHistory()
})

const applyCoupon = async () => {
  applying.value = true
  try {
    const response = await axios.post('/api/coupons/apply', {
      code: couponCode.value.trim().toUpperCase(),
    })
    toast.success(response.data.message + '\n' + response.data.benefit_applied, 4000)
    couponCode.value = ''
    await fetchUsageHistory()
  } catch (error) {
    const errorData = error.response?.data
    let errorMessage = errorData?.error || 'Failed to apply coupon'
    
    // 处理时间相关的错误，显示本地时间
    if (errorData?.reason === 'coupon_not_yet_valid' && errorData?.valid_from) {
      const validFrom = new Date(errorData.valid_from)
      errorMessage += `\nThis coupon will be available after ${validFrom.toLocaleString()}`
    } else if (errorData?.reason === 'coupon_expired' && errorData?.valid_until) {
      const validUntil = new Date(errorData.valid_until)
      errorMessage += `\nThis coupon expired on ${validUntil.toLocaleString()}`
    } else if (errorData?.reason) {
      errorMessage += `\n${errorData.reason}`
    }
    
    toast.error(errorMessage)
  } finally {
    applying.value = false
  }
}

const fetchUsageHistory = async () => {
  loadingHistory.value = true
  try {
    const response = await axios.get('/api/coupons/my-usage')
    usageHistory.value = response.data.usages
  } catch (error) {
    console.error('Failed to fetch usage history:', error)
  } finally {
    loadingHistory.value = false
  }
}

const formatType = (type) => {
  const types = {
    quota_increase: t('coupon.quotaIncrease'),
    percentage: t('coupon.percentageOff'),
    fixed: t('coupon.fixedDiscount'),
  }
  return types[type] || type
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const getBadgeClass = (type) => {
  const classes = {
    quota_increase: 'y-chip-mint',
    percentage: 'y-chip-info',
    fixed: 'y-chip-info',
  }
  return classes[type] || 'y-chip'
}
</script>
