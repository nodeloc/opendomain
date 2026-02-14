<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">{{ $t('dnsManagement.title') }}</h1>
        <p class="text-lg opacity-70 mt-2">{{ domain?.full_domain }}</p>
      </div>
      <div class="flex gap-2">
        <button @click="syncFromPowerDNS" class="btn btn-outline gap-2" :disabled="syncing">
          <span v-if="syncing" class="loading loading-spinner loading-sm"></span>
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ $t('dnsManagement.syncFromPowerDNS') }}
        </button>
        <button @click="showAddModal = true" class="btn btn-primary gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ $t('dnsManagement.addRecord') }}
        </button>
      </div>
    </div>

    <!-- DNS Records Table -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="records.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <div class="mb-6">
          <div class="w-24 h-24 rounded-full bg-base-300 flex items-center justify-center mx-auto">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-content opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
            </svg>
          </div>
        </div>
        <h3 class="card-title text-2xl mb-2">{{ $t('dnsManagement.noRecords') }}</h3>
        <p class="mb-4 opacity-70">{{ $t('dnsManagement.noRecordsDesc') }}</p>
        <button @click="showAddModal = true" class="btn btn-primary gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ $t('dnsManagement.addRecord') }}
        </button>
      </div>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>{{ $t('dnsManagement.name') }}</th>
            <th>{{ $t('dnsManagement.type') }}</th>
            <th>{{ $t('dnsManagement.content') }}</th>
            <th>{{ $t('dnsManagement.ttl') }}</th>
            <th>{{ $t('dnsManagement.priority') }}</th>
            <th>{{ $t('dnsManagement.status') }}</th>
            <th>{{ $t('dnsManagement.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id">
            <td class="font-mono">{{ record.name }}</td>
            <td>
              <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-primary/10 text-primary">
                {{ record.type }}
              </span>
            </td>
            <td class="font-mono text-sm">{{ truncate(record.content, 40) }}</td>
            <td>{{ record.ttl }}s</td>
            <td>{{ record.priority || '-' }}</td>
            <td>
              <div class="flex items-center gap-2">
                <span
                  class="px-2 py-0.5 rounded-full text-xs font-medium inline-flex items-center gap-1"
                  :class="record.synced_to_powerdns ? 'bg-success/10 text-success' : 'bg-warning/10 text-warning'"
                >
                  <svg v-if="record.synced_to_powerdns" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {{ record.synced_to_powerdns ? $t('dnsManagement.synced') : $t('dnsManagement.pending') }}
                </span>
                <span v-if="!record.is_active" class="px-2 py-0.5 rounded-full text-xs font-medium bg-error/10 text-error">
                  {{ $t('dnsManagement.inactive') }}
                </span>
                <div v-if="record.sync_error" class="tooltip tooltip-error" :data-tip="$t('dnsManagement.clickToViewError')">
                  <button @click="showErrorDetails(record.sync_error)" class="btn btn-ghost btn-xs btn-square">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-error" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                  </button>
                </div>
              </div>
            </td>
            <td>
              <div class="flex gap-1">
                <button @click="editRecord(record)" class="btn btn-ghost btn-sm btn-square" :title="$t('dnsManagement.editTitle')">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button @click="confirmDelete(record)" class="btn btn-ghost btn-sm btn-square text-error" :title="$t('dnsManagement.deleteConfirm')">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add/Edit Modal -->
    <dialog ref="modalRef" class="modal" :class="{ 'modal-open': showAddModal || showEditModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="font-bold text-lg">
          {{ showEditModal ? $t('dnsManagement.editTitle') : $t('dnsManagement.addTitle') }}
        </h3>

        <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('dnsManagement.name') }}</span></label>
              <input
                v-model="form.name"
                type="text"
                :placeholder="$t('dnsManagement.namePlaceholder')"
                class="input input-bordered"
                required
              />
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('dnsManagement.type') }}</span></label>
              <select v-model="form.type" class="select select-bordered" required>
                <option value="A">A</option>
                <option value="AAAA">AAAA</option>
                <option value="CNAME">CNAME</option>
                <option value="MX">MX</option>
                <option value="TXT">TXT</option>
                <option value="NS">NS</option>
                <option value="SRV">SRV</option>
                <option value="CAA">CAA</option>
              </select>
            </div>
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('dnsManagement.content') }}</span></label>
            <textarea
              v-model="form.content"
              class="textarea textarea-bordered"
              rows="3"
              :placeholder="$t('dnsManagement.contentPlaceholder')"
              required
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">{{ $t('dnsManagement.ttlSeconds') }}</span></label>
              <input
                v-model.number="form.ttl"
                type="number"
                min="60"
                max="86400"
                class="input input-bordered"
                required
              />
            </div>

            <div v-if="form.type === 'MX' || form.type === 'SRV'" class="form-control">
              <label class="label"><span class="label-text">{{ $t('dnsManagement.priority') }}</span></label>
              <input
                v-model.number="form.priority"
                type="number"
                min="0"
                class="input input-bordered"
              />
            </div>
          </div>

          <div v-if="showEditModal" class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">{{ $t('dnsManagement.active') }}</span>
              <input v-model="form.is_active" type="checkbox" class="checkbox" />
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeModal" class="btn">{{ $t('common.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner"></span>
              <span v-else>{{ showEditModal ? $t('dnsManagement.update') : $t('dnsManagement.create') }}</span>
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>close</button>
      </form>
    </dialog>

    <!-- Error Details Modal -->
    <dialog class="modal" :class="{ 'modal-open': showErrorModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg text-error flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          {{ $t('dnsManagement.syncErrorTitle') }}
        </h3>
        
        <div class="mt-4">
          <div class="bg-base-200 p-4 rounded-lg">
            <pre class="text-sm whitespace-pre-wrap break-words font-mono">{{ currentError }}</pre>
          </div>
        </div>

        <div class="modal-action">
          <button @click="copyError" class="btn btn-outline gap-2">
            <svg v-if="!errorCopied" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            {{ errorCopied ? $t('dnsManagement.copied') : $t('dnsManagement.copyError') }}
          </button>
          <button @click="closeErrorModal" class="btn">{{ $t('common.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeErrorModal">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const route = useRoute()
const toast = useToast()
const domainId = route.params.domainId

const domain = ref(null)
const records = ref([])
const loading = ref(true)
const syncing = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const showErrorModal = ref(false)
const submitting = ref(false)
const editingRecord = ref(null)
const currentError = ref('')
const errorCopied = ref(false)

const form = ref({
  name: '@',
  type: 'A',
  content: '',
  ttl: 3600,
  priority: null,
  is_active: true,
})

onMounted(async () => {
  await fetchDomain()
  await fetchRecords()
})

const fetchDomain = async () => {
  try {
    const response = await axios.get(`/api/domains/${domainId}`)
    domain.value = response.data
  } catch (error) {
    console.error('Failed to fetch domain:', error)
  }
}

const fetchRecords = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/api/dns/${domainId}/records`)
    records.value = response.data.records
  } catch (error) {
    console.error('Failed to fetch DNS records:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (showEditModal.value) {
      await axios.put(`/api/dns/${domainId}/records/${editingRecord.value.id}`, form.value)
      toast.success(t('dnsManagement.recordUpdated'))
    } else {
      await axios.post(`/api/dns/${domainId}/records`, form.value)
      toast.success(t('dnsManagement.recordAdded'))
    }
    closeModal()
    await fetchRecords()
  } catch (error) {
    toast.error(error.response?.data?.error || t('dnsManagement.operationFailed'))
  } finally {
    submitting.value = false
  }
}

const editRecord = (record) => {
  editingRecord.value = record
  form.value = {
    name: record.name,
    type: record.type,
    content: record.content,
    ttl: record.ttl,
    priority: record.priority,
    is_active: record.is_active,
  }
  showEditModal.value = true
}

const confirmDelete = async (record) => {
  if (confirm(t('dnsManagement.deleteConfirm'))) {
    try {
      await axios.delete(`/api/dns/${domainId}/records/${record.id}`)
      toast.success(t('dnsManagement.recordDeleted'))
      await fetchRecords()
    } catch (error) {
      toast.error(error.response?.data?.error || t('dnsManagement.operationFailed'))
    }
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingRecord.value = null
  form.value = {
    name: '@',
    type: 'A',
    content: '',
    ttl: 3600,
    priority: null,
    is_active: true,
  }
}

const syncFromPowerDNS = async () => {
  if (!confirm(t('dnsManagement.syncConfirm'))) {
    return
  }

  syncing.value = true
  try {
    const response = await axios.post(`/api/dns/${domainId}/records/sync-from-powerdns`)
    const stats = response.data.stats
    toast.success(
      t('dnsManagement.syncSuccess', {
        created: stats.created,
        updated: stats.updated,
        skipped: stats.skipped
      })
    )
    await fetchRecords()
  } catch (error) {
    toast.error(error.response?.data?.error || t('dnsManagement.syncFailed'))
  } finally {
    syncing.value = false
  }
}

const truncate = (str, length) => {
  return str.length > length ? str.substring(0, length) + '...' : str
}

const showErrorDetails = (error) => {
  currentError.value = error
  errorCopied.value = false
  showErrorModal.value = true
}

const closeErrorModal = () => {
  showErrorModal.value = false
  currentError.value = ''
  errorCopied.value = false
}

const copyError = async () => {
  try {
    await navigator.clipboard.writeText(currentError.value)
    errorCopied.value = true
    setTimeout(() => {
      errorCopied.value = false
    }, 2000)
  } catch (error) {
    console.error('Failed to copy error:', error)
    toast.error(t('dnsManagement.copyFailed'))
  }
}
</script>
