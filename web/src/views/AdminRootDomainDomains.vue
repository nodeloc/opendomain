<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <div class="flex items-center gap-3 mb-2">
          <router-link to="/admin/root-domains" class="btn btn-ghost btn-sm gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            {{ $t('admin.backToRootDomains') }}
          </router-link>
        </div>
        <h1 class="text-3xl font-bold">{{ rootDomain ? $t('admin.domainsUnder', { domain: rootDomain.domain }) : '...' }}</h1>
      </div>
    </div>

    <!-- Search -->
    <div class="flex gap-4">
      <div class="form-control flex-1 max-w-md">
        <div class="relative">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 absolute left-3 top-1/2 -translate-y-1/2 opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('admin.searchDomains')"
            class="input input-bordered w-full pl-10"
            @input="debouncedSearch"
          />
        </div>
      </div>
      <div class="text-sm opacity-70 flex items-center">
        {{ domains.length }} {{ $t('admin.totalDomains') }}
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Empty State -->
    <div v-else-if="domains.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <div class="mb-6">
          <div class="w-24 h-24 rounded-full bg-base-300 flex items-center justify-center mx-auto">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-base-content opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
          </div>
        </div>
        <h3 class="card-title text-2xl mb-2">{{ $t('admin.noDomains') }}</h3>
        <p class="mb-4 opacity-70">{{ $t('admin.noDomainsDesc') }}</p>
      </div>
    </div>

    <!-- Domains Table -->
    <div v-else class="overflow-x-auto bg-base-100 rounded-lg shadow">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>{{ $t('admin.totalDomains') }}</th>
            <th>{{ $t('admin.owner') }}</th>
            <th>{{ $t('admin.status') }}</th>
            <th>{{ $t('admin.dnsSynced') }}</th>
            <th>{{ $t('admin.registeredAt') }}</th>
            <th>{{ $t('admin.expiresAt') }}</th>
            <th>{{ $t('admin.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="domain in domains" :key="domain.id">
            <td>
              <div class="font-mono font-bold">{{ domain.full_domain }}</div>
              <div class="text-xs opacity-50">ID: {{ domain.id }}</div>
            </td>
            <td>
              <div v-if="domain.user">
                <div class="font-bold">{{ domain.user.username }}</div>
                <div class="text-xs opacity-50">{{ domain.user.email }}</div>
              </div>
              <span v-else class="opacity-40">-</span>
            </td>
            <td>
              <span
                class="px-2 py-0.5 rounded-full text-xs font-medium inline-flex items-center gap-1"
                :class="getStatusClass(domain.status)"
              >
                <svg v-if="domain.status === 'active'" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else-if="domain.status === 'suspended'" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ getStatusLabel(domain.status) }}
              </span>
            </td>
            <td>
              <span
                class="px-2 py-0.5 rounded-full text-xs font-medium inline-flex items-center gap-1"
                :class="domain.dns_synced ? 'bg-success/10 text-success' : 'bg-warning/10 text-warning'"
              >
                <svg v-if="domain.dns_synced" xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ domain.dns_synced ? $t('admin.yes') : $t('admin.no') }}
              </span>
            </td>
            <td class="text-sm">{{ formatDate(domain.registered_at) }}</td>
            <td class="text-sm">{{ formatDate(domain.expires_at) }}</td>
            <td>
              <div class="flex gap-1">
                <button
                  v-if="domain.user"
                  @click="viewUser(domain.user)"
                  class="btn btn-ghost btn-sm btn-square"
                  :title="$t('admin.viewUser')"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                </button>
                <button
                  v-if="domain.status === 'active'"
                  @click="updateDomainStatus(domain, 'suspended')"
                  class="btn btn-ghost btn-sm btn-square text-warning"
                  :title="$t('admin.suspend')"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                  </svg>
                </button>
                <button
                  v-else-if="domain.status === 'suspended'"
                  @click="updateDomainStatus(domain, 'active')"
                  class="btn btn-ghost btn-sm btn-square text-success"
                  :title="$t('admin.activate')"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
                <button
                  @click="deleteDomain(domain)"
                  class="btn btn-ghost btn-sm btn-square text-error"
                  :title="$t('admin.delete')"
                >
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

    <!-- User Detail Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showUserModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="font-bold text-lg mb-4">{{ $t('admin.userInfo') }}</h3>
        <div v-if="selectedUser" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="label"><span class="label-text font-semibold">ID</span></label>
              <div>{{ selectedUser.id }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.users') }}</span></label>
              <div class="font-bold">{{ selectedUser.username }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">Email</span></label>
              <div>{{ selectedUser.email }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.status') }}</span></label>
              <div>
                <span
                  class="px-2 py-0.5 rounded-full text-xs font-medium inline-flex items-center gap-1"
                  :class="getUserStatusClass(selectedUser.status)"
                >
                  {{ selectedUser.status }}
                </span>
              </div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('admin.registeredAt') }}</span></label>
              <div>{{ formatDate(selectedUser.created_at) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">Domain Quota</span></label>
              <div>{{ selectedUser.domain_quota }}</div>
            </div>
          </div>

          <div class="divider"></div>

          <div class="flex gap-2">
            <button
              v-if="selectedUser.status === 'active'"
              @click="updateUserStatus(selectedUser, 'banned')"
              class="btn btn-warning btn-sm gap-1"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
              </svg>
              {{ $t('admin.suspend') }}
            </button>
            <button
              v-else-if="selectedUser.status === 'banned' || selectedUser.status === 'frozen'"
              @click="updateUserStatus(selectedUser, 'active')"
              class="btn btn-success btn-sm gap-1"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ $t('admin.activate') }}
            </button>
          </div>
        </div>
        <div class="modal-action">
          <button @click="showUserModal = false" class="btn">{{ $t('admin.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="showUserModal = false">close</button>
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
const rootDomainId = route.params.id

const rootDomain = ref(null)
const domains = ref([])
const loading = ref(true)
const searchQuery = ref('')
const showUserModal = ref(false)
const selectedUser = ref(null)

let searchTimeout = null

onMounted(async () => {
  await fetchDomains()
})

const fetchDomains = async () => {
  loading.value = true
  try {
    const params = {}
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    const response = await axios.get(`/api/admin/root-domains/${rootDomainId}/domains`, { params })
    domains.value = response.data.domains || []
    rootDomain.value = response.data.root_domain
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.operationFailed'))
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchDomains()
  }, 300)
}

const updateDomainStatus = async (domain, status) => {
  const msg = status === 'suspended'
    ? t('admin.confirmSuspendDomain')
    : t('admin.confirmActivateDomain')

  if (!confirm(msg)) return

  try {
    await axios.put(`/api/admin/domains/${domain.id}/status`, { status })
    toast.success(status === 'suspended' ? t('admin.domainSuspended') : t('admin.domainActivated'))
    await fetchDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.operationFailed'))
  }
}

const deleteDomain = async (domain) => {
  if (!confirm(t('admin.confirmDeleteDomain'))) return

  try {
    await axios.delete(`/api/admin/domains/${domain.id}`)
    toast.success(t('admin.domainDeleted'))
    await fetchDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.operationFailed'))
  }
}

const viewUser = (user) => {
  selectedUser.value = { ...user }
  showUserModal.value = true
}

const updateUserStatus = async (user, status) => {
  const msg = status === 'banned'
    ? t('admin.confirmSuspendUser')
    : t('admin.confirmActivateUser')

  if (!confirm(msg)) return

  try {
    const response = await axios.put(`/api/admin/users/${user.id}/status`, { status })
    toast.success(status === 'banned' ? t('admin.userSuspended') : t('admin.userActivated'))
    selectedUser.value = response.data.user
    // Update user status in domain list too
    domains.value.forEach(d => {
      if (d.user && d.user.id === user.id) {
        d.user.status = status
      }
    })
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.operationFailed'))
  }
}

const getStatusClass = (status) => {
  const classes = {
    active: 'bg-success/10 text-success',
    suspended: 'bg-warning/10 text-warning',
    expired: 'bg-error/10 text-error',
  }
  return classes[status] || 'bg-base-300 text-base-content'
}

const getStatusLabel = (status) => {
  const labels = {
    active: t('admin.active'),
    suspended: t('admin.suspended'),
    expired: t('admin.expired'),
  }
  return labels[status] || status
}

const getUserStatusClass = (status) => {
  const classes = {
    active: 'bg-success/10 text-success',
    frozen: 'bg-warning/10 text-warning',
    banned: 'bg-error/10 text-error',
  }
  return classes[status] || 'bg-base-300 text-base-content'
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}
</script>
