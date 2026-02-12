<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">{{ $t('admin.couponManagement') }}</h1>
        <p class="text-lg opacity-70 mt-2">{{ $t('admin.couponManagementDesc') }}</p>
      </div>
      <button @click="showCreateModal = true" class="btn btn-primary">
        + {{ $t('admin.createCoupon') }}
      </button>
    </div>

    <!-- Coupons Table -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="coupons.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 opacity-40 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
        </svg>
        <h3 class="card-title text-2xl mb-2">{{ $t('admin.noCoupons') }}</h3>
        <p class="mb-4">{{ $t('admin.noCouponsDesc') }}</p>
        <button @click="showCreateModal = true" class="btn btn-primary">{{ $t('admin.createCoupon') }}</button>
      </div>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>{{ $t('admin.code') }}</th>
            <th>{{ $t('admin.type') }}</th>
            <th>{{ $t('admin.value') }}</th>
            <th>{{ $t('admin.uses') }}</th>
            <th>{{ $t('admin.validFrom') }}</th>
            <th>{{ $t('admin.validUntil') }}</th>
            <th>{{ $t('admin.status') }}</th>
            <th>{{ $t('admin.reusable') }}</th>
            <th>{{ $t('admin.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="coupon in coupons" :key="coupon.id">
            <td class="font-mono font-bold">{{ coupon.code }}</td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1 bg-blue-100 text-blue-800">
                <!-- Type icons -->
                <svg v-if="coupon.discount_type === 'quota_increase'" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 11l5-5m0 0l5 5m-5-5v12" />
                </svg>
                <svg v-else-if="coupon.discount_type === 'percentage'" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <svg v-else-if="coupon.discount_type === 'fixed'" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ formatType(coupon.discount_type) }}
              </div>
            </td>
            <td>{{ formatValue(coupon) }}</td>
            <td>{{ coupon.used_count }} / {{ coupon.max_uses || '∞' }}</td>
            <td>{{ formatDate(coupon.valid_from) }}</td>
            <td>{{ coupon.valid_until ? formatDate(coupon.valid_until) : $t('admin.noExpiry') }}</td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1" :class="coupon.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                <svg v-if="coupon.is_active" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ coupon.is_active ? $t('admin.active') : $t('admin.inactive') }}
              </div>
            </td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1" :class="coupon.is_reusable ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'">
                <svg v-if="coupon.is_reusable" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                {{ coupon.is_reusable ? $t('admin.yes') : $t('admin.no') }}
              </div>
            </td>
            <td>
              <div class="flex gap-2">
                <button @click="editCoupon(coupon)" class="btn btn-sm btn-ghost">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button @click="confirmDelete(coupon)" class="btn btn-sm btn-ghost text-error">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <dialog class="modal" :class="{ 'modal-open': showCreateModal || showEditModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="font-bold text-lg">
          {{ showEditModal ? t('admin.editCoupon') : t('admin.createCoupon') }}
        </h3>

        <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
          <div class="form-control">
            <label class="label"><span class="label-text">{{ t('admin.couponCode') }}</span></label>
            <input
              v-model="form.code"
              type="text"
              placeholder="SAVE20"
              class="input input-bordered uppercase"
              :disabled="showEditModal"
              required
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ t('admin.description') }}</span></label>
            <textarea
              v-model="form.description"
              class="textarea textarea-bordered"
              rows="2"
              placeholder="Optional description"
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">{{ t('admin.type') }}</span></label>
              <select v-model="form.discount_type" class="select select-bordered" :disabled="showEditModal" required>
                <option value="quota_increase">{{ t('admin.quotaIncrease') }}</option>
                <option value="percentage">{{ t('admin.percentage') }}</option>
                <option value="fixed">{{ t('admin.fixed') }}</option>
              </select>
            </div>

            <div v-if="form.discount_type === 'quota_increase'" class="form-control">
              <label class="label"><span class="label-text">{{ t('admin.quotaIncrease') }}</span></label>
              <input
                v-model.number="form.quota_increase"
                type="number"
                min="1"
                class="input input-bordered"
                required
              />
            </div>

            <div v-if="form.discount_type === 'percentage' || form.discount_type === 'fixed'" class="form-control">
              <label class="label"><span class="label-text">{{ t('admin.discountValue') }}</span></label>
              <input
                v-model.number="form.discount_value"
                type="number"
                :min="0"
                :max="form.discount_type === 'percentage' ? 100 : undefined"
                step="0.01"
                class="input input-bordered"
                required
              />
            </div>
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ t('admin.maxUses') }} {{ t('admin.maxUsesHint') }}</span></label>
            <input
              v-model.number="form.max_uses"
              type="number"
              min="0"
              class="input input-bordered"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">{{ t('admin.validFrom') }}</span></label>
              <input
                v-model="form.valid_from"
                type="datetime-local"
                class="input input-bordered"
              />
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text">{{ t('admin.validUntilOptional') }}</span></label>
              <input
                v-model="form.valid_until"
                type="datetime-local"
                class="input input-bordered"
              />
            </div>
          </div>

          <div v-if="showEditModal" class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">{{ t('admin.active') }}</span>
              <input v-model="form.is_active" type="checkbox" class="checkbox" />
            </label>
          </div>

          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">{{ t('admin.reusable') }} ({{ t('admin.reusableDesc') }})</span>
              <input v-model="form.is_reusable" type="checkbox" class="checkbox" />
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeModal" class="btn">{{ t('admin.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner"></span>
              <span v-else>{{ showEditModal ? t('admin.update') : t('admin.create') }}</span>
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const toast = useToast()
const coupons = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)
const editingCoupon = ref(null)

const form = ref({
  code: '',
  description: '',
  discount_type: 'quota_increase',
  discount_value: null,
  quota_increase: 1,
  max_uses: 0,
  valid_from: null,
  valid_until: null,
  is_active: true,
  is_reusable: false,
})

onMounted(async () => {
  await fetchCoupons()
})

const fetchCoupons = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/coupons')
    coupons.value = response.data.coupons
  } catch (error) {
    console.error('Failed to fetch coupons:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const payload = {
      ...form.value,
      code: form.value.code.toUpperCase(),
    }

    if (showEditModal.value) {
      await axios.put(`/api/admin/coupons/${editingCoupon.value.id}`, payload)
      toast.success(t('coupon.updateSuccess'))
    } else {
      await axios.post('/api/admin/coupons', payload)
      toast.success(t('coupon.createSuccess'))
    }
    closeModal()
    await fetchCoupons()
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.operationFailed'))
  } finally {
    submitting.value = false
  }
}

const editCoupon = (coupon) => {
  editingCoupon.value = coupon
  form.value = {
    code: coupon.code,
    description: coupon.description || '',
    discount_type: coupon.discount_type,
    discount_value: coupon.discount_value,
    quota_increase: coupon.quota_increase || 1,
    max_uses: coupon.max_uses || 0,
    valid_from: coupon.valid_from ? new Date(coupon.valid_from).toISOString().slice(0, 16) : null,
    valid_until: coupon.valid_until ? new Date(coupon.valid_until).toISOString().slice(0, 16) : null,
    is_active: coupon.is_active,
    is_reusable: coupon.is_reusable || false,
  }
  showEditModal.value = true
}

const confirmDelete = async (coupon) => {
  if (confirm(t('coupon.deleteConfirm'))) {
    try {
      await axios.delete(`/api/admin/coupons/${coupon.id}`)
      toast.success(t('coupon.deleteSuccess'))
      await fetchCoupons()
    } catch (error) {
      toast.error(error.response?.data?.error || t('admin.operationFailed'))
    }
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingCoupon.value = null
  form.value = {
    code: '',
    description: '',
    discount_type: 'quota_increase',
    discount_value: null,
    quota_increase: 1,
    max_uses: 0,
    valid_from: null,
    valid_until: null,
    is_active: true,
    is_reusable: false,
  }
}

const formatType = (type) => {
  const types = {
    quota_increase: 'admin.quotaIncrease',
    percentage: 'admin.percentage',
    fixed: 'admin.fixed',
  }
  return types[type] ? t(types[type]) : type
}

const formatValue = (coupon) => {
  switch (coupon.discount_type) {
    case 'quota_increase':
      return t('admin.valueDomains', { n: coupon.quota_increase })
    case 'percentage':
      return `${coupon.discount_value}%`
    case 'fixed':
      return `$${coupon.discount_value}`
    default:
      return '-'
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
