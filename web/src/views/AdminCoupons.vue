<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('uiAdmin2.eyebrow')"
      :title="$t('uiAdmin2.coupons.title')"
      :description="$t('uiAdmin2.coupons.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-primary gap-1.5" @click="showCreateModal = true">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('uiAdmin2.coupons.create') }}
        </button>
      </template>
    </PageHeader>

    <section class="y-panel">
      <YLoader v-if="loading" :label="$t('uiAdmin2.loading')" />

      <EmptyState
        v-else-if="coupons.length === 0"
        :title="$t('uiAdmin2.coupons.emptyTitle')"
        :description="$t('uiAdmin2.coupons.emptyDescription')"
      >
        <button type="button" class="btn btn-outline btn-sm gap-1.5" @click="showCreateModal = true">
          <Plus :size="14" :stroke-width="1.5" />
          {{ $t('uiAdmin2.coupons.create') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
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
              <th class="text-right">{{ $t('admin.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="coupon in coupons" :key="coupon.id">
              <td>
                <p class="whitespace-nowrap font-mono font-medium text-y-text">{{ coupon.code }}</p>
                <p v-if="coupon.description" class="max-w-[16rem] truncate text-xs text-y-faint">{{ coupon.description }}</p>
              </td>
              <td>
                <span class="y-chip whitespace-nowrap">
                  <component :is="couponTypeIcon(coupon.discount_type)" :size="12" :stroke-width="2" class="text-y-mint-text" />
                  {{ formatType(coupon.discount_type) }}
                </span>
              </td>
              <td class="whitespace-nowrap font-mono">{{ formatValue(coupon) }}</td>
              <td class="whitespace-nowrap font-mono text-y-muted">
                <span class="text-y-text">{{ coupon.used_count }}</span> / {{ coupon.max_uses || '∞' }}
              </td>
              <td class="min-w-[8.5rem] text-xs leading-5 text-y-muted">{{ formatDate(coupon.valid_from) }}</td>
              <td class="min-w-[8.5rem] text-xs leading-5 text-y-muted">
                {{ coupon.valid_until ? formatDate(coupon.valid_until) : $t('admin.noExpiry') }}
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="coupon.is_active ? 'ok' : 'muted'" :label="coupon.is_active ? $t('admin.active') : $t('admin.inactive')" />
              </td>
              <td class="whitespace-nowrap">
                <span v-if="coupon.is_reusable" class="y-chip-info">
                  <Repeat :size="12" :stroke-width="2" />
                  {{ $t('uiAdmin2.coupons.reusable') }}
                </span>
                <span v-else class="text-sm text-y-faint">{{ $t('uiAdmin2.coupons.singleUse') }}</span>
              </td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('uiAdmin2.edit')"
                    :aria-label="$t('uiAdmin2.edit')"
                    @click="editCoupon(coupon)"
                  >
                    <Pencil :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('uiAdmin2.delete')"
                    :aria-label="$t('uiAdmin2.delete')"
                    @click="confirmDelete(coupon)"
                  >
                    <Trash2 :size="16" :stroke-width="1.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Create/Edit Modal -->
    <dialog class="modal" :class="{ 'modal-open': showCreateModal || showEditModal }">
      <div class="modal-box relative w-11/12 max-w-2xl">
        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square absolute right-3 top-3"
          :aria-label="$t('uiAdmin2.close')"
          @click="closeModal"
        >
          <X :size="16" :stroke-width="1.5" />
        </button>
        <h3 class="text-lg font-semibold text-y-text">
          {{ showEditModal ? t('uiAdmin2.coupons.editTitle') : t('uiAdmin2.coupons.createTitle') }}
        </h3>
        <p class="mt-1 text-sm text-y-muted">{{ t('uiAdmin2.coupons.modalDescription') }}</p>

        <form class="mt-6 space-y-4" @submit.prevent="handleSubmit">
          <label class="y-field">
            <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.code') }}</span>
            <input
              v-model="form.code"
              type="text"
              placeholder="SAVE20"
              class="input input-bordered w-full font-mono uppercase"
              :disabled="showEditModal"
              autocomplete="off"
              spellcheck="false"
              required
            />
          </label>

          <label class="y-field">
            <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.description') }}</span>
            <textarea
              v-model="form.description"
              class="textarea textarea-bordered w-full"
              rows="2"
              :placeholder="t('uiAdmin2.coupons.fields.descriptionPlaceholder')"
            ></textarea>
          </label>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.type') }}</span>
              <select v-model="form.discount_type" class="select select-bordered w-full" :disabled="showEditModal" required>
                <option value="quota_increase">{{ t('admin.quotaIncrease') }}</option>
                <option value="percentage">{{ t('admin.percentage') }}</option>
                <option value="fixed">{{ t('admin.fixed') }}</option>
              </select>
            </label>

            <label v-if="form.discount_type === 'quota_increase'" class="y-field">
              <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.quotaIncrease') }}</span>
              <input
                v-model.number="form.quota_increase"
                type="number"
                min="1"
                class="input input-bordered w-full font-mono"
                required
              />
            </label>

            <label v-if="form.discount_type === 'percentage' || form.discount_type === 'fixed'" class="y-field">
              <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.discountValue') }}</span>
              <input
                v-model.number="form.discount_value"
                type="number"
                :min="0"
                :max="form.discount_type === 'percentage' ? 100 : undefined"
                step="0.01"
                class="input input-bordered w-full font-mono"
                required
              />
              <span class="y-field-hint">
                {{ form.discount_type === 'percentage' ? t('uiAdmin2.coupons.fields.discountPercentHint') : t('uiAdmin2.coupons.fields.discountFixedHint') }}
              </span>
            </label>
          </div>

          <label class="y-field">
            <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.maxUses') }}</span>
            <input
              v-model.number="form.max_uses"
              type="number"
              min="0"
              class="input input-bordered w-full font-mono"
            />
            <span class="y-field-hint">{{ t('uiAdmin2.coupons.fields.maxUsesHint') }}</span>
          </label>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.validFrom') }}</span>
              <input
                v-model="form.valid_from"
                type="datetime-local"
                class="input input-bordered w-full font-mono"
              />
            </label>

            <label class="y-field">
              <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.validUntil') }}</span>
              <input
                v-model="form.valid_until"
                type="datetime-local"
                class="input input-bordered w-full font-mono"
              />
              <span class="y-field-hint">{{ t('uiAdmin2.coupons.fields.validUntilHint') }}</span>
            </label>
          </div>

          <div class="space-y-2">
            <label v-if="showEditModal" class="flex cursor-pointer items-start gap-3 rounded-md border border-y-line px-4 py-3">
              <input v-model="form.is_active" type="checkbox" class="checkbox checkbox-sm mt-0.5" />
              <span class="flex flex-col gap-0.5">
                <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.active') }}</span>
                <span class="y-field-hint">{{ t('uiAdmin2.coupons.fields.activeHint') }}</span>
              </span>
            </label>

            <label class="flex cursor-pointer items-start gap-3 rounded-md border border-y-line px-4 py-3">
              <input v-model="form.is_reusable" type="checkbox" class="checkbox checkbox-sm mt-0.5" />
              <span class="flex flex-col gap-0.5">
                <span class="y-field-label">{{ t('uiAdmin2.coupons.fields.reusable') }}</span>
                <span class="y-field-hint">{{ t('uiAdmin2.coupons.fields.reusableHint') }}</span>
              </span>
            </label>
          </div>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModal">{{ t('uiAdmin2.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              {{ showEditModal ? t('uiAdmin2.saveChanges') : t('uiAdmin2.coupons.create') }}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>{{ t('uiAdmin2.close') }}</button>
      </form>
    </dialog>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Plus, Pencil, Trash2, X, Repeat, ArrowUp, Percent, Coins, Ticket } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()
const coupons = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)
const editingCoupon = ref(null)

// 时区转换辅助函数：将UTC时间字符串转换为本地时间字符串（用于datetime-local输入）
const toLocalDatetimeString = (utcDateString) => {
  if (!utcDateString) return null
  const date = new Date(utcDateString)
  // 获取本地时间的各个部分
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

// 时区转换辅助函数：将本地时间字符串（来自datetime-local输入）转换为UTC ISO字符串
const toUTCString = (localDatetimeString) => {
  if (!localDatetimeString) return null
  // datetime-local 返回的格式类似 "2024-01-01T12:00"
  // 创建Date对象时，它会被当作本地时间
  const date = new Date(localDatetimeString)
  return date.toISOString()
}

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
      // 将本地时间转换为UTC时间
      valid_from: toUTCString(form.value.valid_from),
      valid_until: toUTCString(form.value.valid_until),
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
    // 将UTC时间转换为本地时间填充表单
    valid_from: toLocalDatetimeString(coupon.valid_from),
    valid_until: toLocalDatetimeString(coupon.valid_until),
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

const couponTypeIcon = (type) => ({ quota_increase: ArrowUp, percentage: Percent, fixed: Coins }[type] || Ticket)

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
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}
</script>
