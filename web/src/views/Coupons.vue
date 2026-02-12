<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-8">
    <!-- Header Section -->
    <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary/10 via-secondary/5 to-accent/10 p-8 backdrop-blur-sm">
      <div class="relative z-10">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-12 h-12 rounded-xl bg-primary/20 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-7 w-7 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h1 class="text-4xl font-bold">{{ $t('coupon.title') }}</h1>
        </div>
        <p class="text-lg opacity-70 max-w-2xl">{{ $t('coupon.applyCode') }}</p>
      </div>
      <div class="absolute top-0 right-0 w-64 h-64 bg-primary/5 rounded-full blur-3xl"></div>
      <div class="absolute bottom-0 left-0 w-48 h-48 bg-secondary/5 rounded-full blur-3xl"></div>
    </div>

    <!-- Apply Coupon Card -->
    <div class="card bg-base-100 shadow-2xl border border-base-300 hover:shadow-primary/20 transition-all duration-300">
      <div class="card-body">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
            </svg>
          </div>
          <h2 class="card-title text-2xl">{{ $t('coupon.applyCode') }}</h2>
        </div>

        <form @submit.prevent="applyCoupon" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold text-base">{{ $t('coupon.enterCode') }}</span>
              <span class="label-text-alt opacity-60">{{ $t('coupon.caseInsensitive') }}</span>
            </label>
            <div class="flex gap-3">
              <div class="relative flex-1">
                <input
                  v-model="couponCode"
                  type="text"
                  :placeholder="$t('coupon.placeholder')"
                  class="input input-bordered input-lg w-full pl-12 font-mono uppercase focus:input-primary transition-all"
                  :disabled="applying"
                />
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-base-content/40">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                  </svg>
                </div>
              </div>
              <button
                type="submit"
                class="btn btn-primary btn-lg px-8 gap-2 hover:scale-105 transition-transform"
                :disabled="!couponCode || applying"
              >
                <span v-if="applying" class="loading loading-spinner loading-sm"></span>
                <span v-else>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </span>
                <span>{{ $t('coupon.apply') }}</span>
              </button>
            </div>
          </div>

          <div class="alert alert-info shadow-lg">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span>{{ $t('coupon.infoText') }}</span>
          </div>
        </form>
      </div>
    </div>

    <!-- Usage History -->
    <div class="card bg-base-100 shadow-2xl border border-base-300">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-secondary/10 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <h2 class="card-title text-2xl">{{ $t('coupon.myHistory') }}</h2>
          </div>
          <div v-if="usageHistory.length > 0" class="px-2 py-1 rounded text-sm font-medium inline-flex items-center gap-2 bg-purple-100 text-purple-800">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ usageHistory.length }} {{ $t('coupon.used') }}
          </div>
        </div>

        <div v-if="loadingHistory" class="flex justify-center py-16">
          <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>

        <div v-else-if="usageHistory.length === 0" class="text-center py-16">
          <div class="flex justify-center mb-6">
            <div class="relative">
              <div class="w-32 h-32 rounded-full bg-base-200 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-base-content/20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                </svg>
              </div>
              <div class="absolute -top-2 -right-2 w-8 h-8 rounded-full bg-primary/20 animate-pulse"></div>
              <div class="absolute -bottom-2 -left-2 w-6 h-6 rounded-full bg-secondary/20 animate-pulse" style="animation-delay: 0.5s"></div>
            </div>
          </div>
          <h3 class="text-xl font-semibold mb-2">{{ $t('coupon.noHistory') }}</h3>
          <p class="text-base-content/60 max-w-md mx-auto">{{ $t('coupon.noHistoryDesc') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
              <tr class="border-b-2 border-base-300">
                <th class="bg-base-200">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
                    </svg>
                    {{ $t('coupon.code') }}
                  </div>
                </th>
                <th class="bg-base-200">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                    </svg>
                    {{ $t('coupon.type') }}
                  </div>
                </th>
                <th class="bg-base-200">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ $t('coupon.benefit') }}
                  </div>
                </th>
                <th class="bg-base-200">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    {{ $t('coupon.usedAt') }}
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="usage in usageHistory" :key="usage.id" class="hover:bg-base-200/50 transition-colors">
                <td>
                  <div class="flex items-center gap-2">
                    <div class="px-2 py-1 rounded border text-sm font-mono font-bold border-base-300">
                      {{ usage.coupon?.code }}
                    </div>
                  </div>
                </td>
                <td>
                  <div :class="getBadgeClass(usage.coupon?.discount_type)">
                    {{ formatType(usage.coupon?.discount_type) }}
                  </div>
                </td>
                <td>
                  <span class="font-semibold text-success">{{ usage.benefit_applied }}</span>
                </td>
                <td>
                  <div class="flex items-center gap-2 text-sm opacity-70">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ formatDate(usage.used_at) }}
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

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
    toast.error(error.response?.data?.error || 'Failed to apply coupon')
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
    quota_increase: 'px-2 py-1 rounded text-sm font-medium inline-flex items-center gap-2 bg-blue-100 text-blue-800',
    percentage: 'px-2 py-1 rounded text-sm font-medium inline-flex items-center gap-2 bg-purple-100 text-purple-800',
    fixed: 'px-2 py-1 rounded text-sm font-medium inline-flex items-center gap-2 bg-pink-100 text-pink-800',
  }
  return classes[type] || 'px-2 py-1 rounded text-sm font-medium inline-flex items-center gap-2 bg-gray-100 text-gray-800'
}
</script>
