<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">{{ $t('admin.rootDomainManagement') }}</h1>
        <p class="text-lg opacity-70 mt-2">{{ $t('admin.rootDomainManagementDesc') }}</p>
      </div>
      <button @click="showCreateModal = true" class="btn btn-primary gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ $t('admin.addRootDomain') }}
      </button>
    </div>

    <!-- Root Domains List -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="rootDomains.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <div class="mb-6">
          <div class="w-24 h-24 rounded-full bg-base-300 flex items-center justify-center mx-auto">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-content opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
          </div>
        </div>
        <h3 class="card-title text-2xl mb-2">{{ $t('admin.noRootDomains') }}</h3>
        <p class="mb-4 opacity-70">{{ $t('admin.noRootDomainsDesc') }}</p>
        <button @click="showCreateModal = true" class="btn btn-primary gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ $t('admin.addRootDomain') }}
        </button>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="domain in rootDomains"
        :key="domain.id"
        class="card bg-base-100 shadow-xl border border-base-300 hover:shadow-2xl transition-all duration-300"
        :class="{ 'opacity-60': !domain.is_active }"
      >
        <div class="card-body">
          <div class="flex justify-between items-start mb-2">
            <div class="flex items-center gap-2 flex-wrap">
              <h2 class="card-title text-2xl font-mono">.{{ domain.domain }}</h2>
              <div class="flex gap-1.5">
                <span v-if="domain.is_hot" class="px-2 py-0.5 rounded-full text-xs font-medium inline-flex items-center gap-1 bg-warning/10 text-warning">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.879 16.121A3 3 0 1012.015 11L11 14H9c0 .768.293 1.536.879 2.121z" />
                  </svg>
                  {{ $t('admin.hot') }}
                </span>
                <span v-if="domain.is_new" class="px-2 py-0.5 rounded-full text-xs font-medium inline-flex items-center gap-1 bg-info/10 text-info">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                  </svg>
                  {{ $t('admin.new') }}
                </span>
              </div>
            </div>
            <div class="flex gap-2">
              <button @click="editDomain(domain)" class="btn btn-ghost btn-sm btn-square" title="Edit">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button @click="deleteDomain(domain)" class="btn btn-ghost btn-sm btn-square text-error" title="Delete">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

          <p v-if="domain.description" class="text-sm opacity-70 mb-4">{{ domain.description }}</p>

          <div class="divider my-2"></div>

          <div class="space-y-2">
            <div class="flex justify-between items-center">
              <span class="text-sm font-medium opacity-70">{{ $t('admin.status') }}:</span>
              <span class="px-3 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1.5" :class="domain.is_active ? 'bg-success/10 text-success' : 'bg-error/10 text-error'">
                <svg v-if="domain.is_active" xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ domain.is_active ? $t('admin.active') : $t('admin.inactive') }}
              </span>
            </div>

            <div class="flex justify-between items-center">
              <span class="text-sm font-medium opacity-70">Priority:</span>
              <span class="px-3 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1.5 bg-primary/10 text-primary">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 11l5-5m0 0l5 5m-5-5v12" />
                </svg>
                {{ domain.priority }}
              </span>
            </div>

            <div class="flex justify-between items-center">
              <span class="text-sm font-medium opacity-70">Registrations:</span>
              <router-link
                :to="`/admin/root-domains/${domain.id}/domains`"
                class="px-3 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1.5 bg-base-300 text-base-content hover:bg-primary/10 hover:text-primary transition-colors cursor-pointer"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                {{ domain.registration_count || 0 }}
              </router-link>
            </div>

            <div class="flex justify-between items-center">
              <span class="text-sm font-medium opacity-70">Pricing:</span>
              <span v-if="domain.is_free" class="px-3 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1.5 bg-success/10 text-success">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Free
              </span>
              <div v-else class="text-right">
                <div v-if="domain.price_per_year" class="text-xs font-medium text-primary">
                  {{ formatPrice(domain.price_per_year) }}/year
                </div>
                <div v-if="domain.lifetime_price" class="text-xs font-medium text-secondary">
                  {{ formatPrice(domain.lifetime_price) }} lifetime
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showCreateModal || showEditModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="font-bold text-2xl mb-4">
          {{ showEditModal ? 'Edit Root Domain' : 'Add New Root Domain' }}
        </h3>

        <form @submit.prevent="showEditModal ? updateDomain() : createDomain()" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">Domain</span>
            </label>
            <div class="relative">
              <input
                v-model="formData.domain"
                type="text"
                placeholder="example.com"
                class="input input-bordered w-full pl-8 font-mono"
                :disabled="showEditModal"
                required
              />
              <div class="absolute left-3 top-1/2 -translate-y-1/2 text-base-content font-mono font-bold text-lg">
                .
              </div>
            </div>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">Description (Optional)</span>
            </label>
            <textarea
              v-model="formData.description"
              class="textarea textarea-bordered h-20"
              placeholder="Brief description of this domain..."
            ></textarea>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text font-semibold">Priority</span>
            </label>
            <input
              v-model.number="formData.priority"
              type="number"
              class="input input-bordered"
              min="0"
              max="100"
            />
            <label class="label">
              <span class="label-text-alt opacity-60">Higher priority shows first (0-100)</span>
            </label>
          </div>

          <div class="divider">Pricing Configuration</div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-3">
              <input
                v-model="formData.is_free"
                type="checkbox"
                class="toggle toggle-success"
              />
              <span class="label-text font-semibold">{{ formData.is_free ? 'Free Domain' : 'Paid Domain' }}</span>
            </label>
            <label class="label">
              <span class="label-text-alt opacity-60">Free domains use user quota, paid domains require payment</span>
            </label>
          </div>

          <div v-if="!formData.is_free" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Price Per Year ({{ currencySymbol }})</span>
              </label>
              <input
                v-model.number="formData.price_per_year"
                type="number"
                step="0.01"
                min="0"
                class="input input-bordered"
                placeholder="0.00"
              />
              <label class="label">
                <span class="label-text-alt opacity-60">Annual registration price</span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Lifetime Price ({{ currencySymbol }})</span>
              </label>
              <input
                v-model.number="formData.lifetime_price"
                type="number"
                step="0.01"
                min="0"
                class="input input-bordered"
                placeholder="0.00"
              />
              <label class="label">
                <span class="label-text-alt opacity-60">One-time permanent registration price</span>
              </label>
            </div>
          </div>

          <div class="divider">Nameservers Configuration</div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-3">
              <input
                v-model="formData.use_default_nameservers"
                type="checkbox"
                class="toggle toggle-primary"
              />
              <span class="label-text font-semibold">{{ formData.use_default_nameservers ? 'Using Default Nameservers' : 'Using Custom Nameservers' }}</span>
            </label>
            <label class="label">
              <span class="label-text-alt opacity-60">Default nameservers are managed by the system</span>
            </label>
          </div>

          <div v-if="formData.use_default_nameservers" class="alert alert-info">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <div class="font-bold">Default Nameservers (Read-only)</div>
              <div class="text-sm">ns1.nodelook.com, ns2.nodelook.com</div>
            </div>
          </div>

          <div v-else class="space-y-2">
            <div v-for="(ns, index) in formData.nameservers" :key="index" class="flex gap-2">
              <input
                v-model="formData.nameservers[index]"
                type="text"
                class="input input-bordered flex-1"
                placeholder="ns1.example.com"
              />
              <button
                v-if="formData.nameservers.length > 1"
                type="button"
                class="btn btn-ghost btn-square"
                @click="formData.nameservers.splice(index, 1)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <button
              type="button"
              class="btn btn-sm btn-ghost gap-2"
              @click="formData.nameservers.push('')"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Add Nameserver
            </button>
          </div>

          <div class="divider">Display Options</div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Status</span>
              </label>
              <label class="label cursor-pointer justify-start gap-3">
                <input
                  v-model="formData.is_active"
                  type="checkbox"
                  class="toggle toggle-success"
                />
                <span class="label-text">{{ formData.is_active ? 'Active' : 'Inactive' }}</span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Mark as Hot</span>
              </label>
              <label class="label cursor-pointer justify-start gap-3">
                <input
                  v-model="formData.is_hot"
                  type="checkbox"
                  class="toggle toggle-warning"
                />
                <span class="label-text">{{ formData.is_hot ? 'Hot' : 'Normal' }}</span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Mark as New</span>
              </label>
              <label class="label cursor-pointer justify-start gap-3">
                <input
                  v-model="formData.is_new"
                  type="checkbox"
                  class="toggle toggle-info"
                />
                <span class="label-text">{{ formData.is_new ? 'New' : 'Regular' }}</span>
              </label>
            </div>
          </div>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModals">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ showEditModal ? 'Update' : 'Create' }}</span>
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeModals">close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'

const toast = useToast()
const { formatPrice, currencySymbol } = useCurrency()
const loading = ref(true)
const submitting = ref(false)
const rootDomains = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingDomain = ref(null)

const formData = ref({
  domain: '',
  description: '',
  priority: 0,
  is_active: true,
  is_hot: false,
  is_new: false,
  is_free: true,
  price_per_year: null,
  lifetime_price: null,
  use_default_nameservers: true,
  nameservers: [''],
})

onMounted(async () => {
  await fetchRootDomains()
})

const fetchRootDomains = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/root-domains')
    rootDomains.value = response.data.root_domains || []
  } catch (error) {
    console.error('Failed to fetch root domains:', error)
    toast.error(error.response?.data?.error || 'Failed to fetch root domains')
  } finally {
    loading.value = false
  }
}

const createDomain = async () => {
  submitting.value = true
  try {
    // Filter out empty nameservers
    const nameservers = formData.value.nameservers.filter(ns => ns.trim() !== '')

    const payload = {
      ...formData.value,
      nameservers: formData.value.use_default_nameservers ? [] : nameservers,
    }

    await axios.post('/api/admin/root-domains', payload)
    toast.success('Root domain created successfully!')
    closeModals()
    await fetchRootDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to create root domain')
  } finally {
    submitting.value = false
  }
}

const editDomain = (domain) => {
  editingDomain.value = domain

  // Parse nameservers from JSON string
  let nameservers = ['']
  try {
    if (domain.nameservers) {
      nameservers = JSON.parse(domain.nameservers)
      if (!Array.isArray(nameservers) || nameservers.length === 0) {
        nameservers = ['']
      }
    }
  } catch (e) {
    console.error('Failed to parse nameservers:', e)
    nameservers = ['']
  }

  formData.value = {
    domain: domain.domain,
    description: domain.description || '',
    priority: domain.priority,
    is_active: domain.is_active,
    is_hot: domain.is_hot,
    is_new: domain.is_new,
    is_free: domain.is_free ?? true,
    price_per_year: domain.price_per_year || null,
    lifetime_price: domain.lifetime_price || null,
    use_default_nameservers: domain.use_default_nameservers ?? true,
    nameservers: nameservers,
  }
  showEditModal.value = true
}

const updateDomain = async () => {
  submitting.value = true
  try {
    // Filter out empty nameservers
    const nameservers = formData.value.nameservers.filter(ns => ns.trim() !== '')

    await axios.put(`/api/admin/root-domains/${editingDomain.value.id}`, {
      description: formData.value.description || null,
      priority: formData.value.priority,
      is_active: formData.value.is_active,
      is_hot: formData.value.is_hot,
      is_new: formData.value.is_new,
      is_free: formData.value.is_free,
      price_per_year: formData.value.is_free ? null : (formData.value.price_per_year || null),
      lifetime_price: formData.value.is_free ? null : (formData.value.lifetime_price || null),
      use_default_nameservers: formData.value.use_default_nameservers,
      nameservers: formData.value.use_default_nameservers ? [] : nameservers,
    })
    toast.success('Root domain updated successfully!')
    closeModals()
    await fetchRootDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to update root domain')
  } finally {
    submitting.value = false
  }
}

const deleteDomain = async (domain) => {
  if (!confirm(`Are you sure you want to delete .${domain.domain}?\n\nThis cannot be undone.`)) {
    return
  }

  try {
    await axios.delete(`/api/admin/root-domains/${domain.id}`)
    toast.success('Root domain deleted successfully!')
    await fetchRootDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to delete root domain')
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingDomain.value = null
  formData.value = {
    domain: '',
    description: '',
    priority: 0,
    is_active: true,
    is_hot: false,
    is_new: false,
    is_free: true,
    price_per_year: null,
    lifetime_price: null,
    use_default_nameservers: true,
    nameservers: [''],
  }
}
</script>
