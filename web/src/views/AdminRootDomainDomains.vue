<template>
  <div class="y-page">
    <div class="space-y-4">
      <router-link to="/admin/root-domains" class="btn btn-ghost btn-sm -ml-3 gap-1.5 text-y-muted">
        <ArrowLeft :size="16" :stroke-width="1.5" />
        {{ $t('admin.backToRootDomains') }}
      </router-link>
      <PageHeader :eyebrow="$t('shell.adminArea')" :description="$t('uiAdmin1.rootDomainDomains.description')">
        <template #title>
          <template v-if="rootDomain">{{ $t('admin.domainsUnder', { domain: rootDomain.domain }) }}</template>
          <span v-else class="text-y-faint">…</span>
        </template>
        <template #actions>
          <span class="y-chip font-mono">{{ $t('uiAdmin1.domainCount', { count: pagination.total }) }}</span>
        </template>
      </PageHeader>
    </div>

    <section class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <label class="relative block w-full sm:max-w-sm">
          <span class="sr-only">{{ $t('common.search') }}</span>
          <Search :size="16" :stroke-width="1.5" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('admin.searchDomains')"
            class="input input-bordered h-10 min-h-0 w-full pl-9 font-mono text-sm"
            @input="debouncedSearch"
          />
        </label>
      </div>

      <YLoader v-if="loading" :label="$t('common.loading')" />

      <EmptyState v-else-if="domains.length === 0" :title="$t('admin.noDomains')" :description="$t('admin.noDomainsDesc')" />

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('adminDomains.domain') }}</th>
              <th>{{ $t('admin.owner') }}</th>
              <th>{{ $t('admin.status') }}</th>
              <th>{{ $t('admin.dnsSynced') }}</th>
              <th>{{ $t('admin.registeredAt') }}</th>
              <th>{{ $t('admin.expiresAt') }}</th>
              <th class="text-right">{{ $t('admin.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="domain in domains" :key="domain.id">
              <td>
                <span class="y-domain block">{{ domain.full_domain }}</span>
                <span class="font-mono text-xs text-y-faint">ID {{ domain.id }}</span>
              </td>
              <td>
                <div v-if="domain.user" class="min-w-0">
                  <span class="block text-sm font-medium text-y-text">{{ domain.user.username }}</span>
                  <span class="block text-xs text-y-faint">{{ domain.user.email }}</span>
                </div>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td>
                <StatusDot :tone="statusTone(domain.status)" :label="getStatusLabel(domain.status)" />
              </td>
              <td>
                <StatusDot :tone="domain.dns_synced ? 'ok' : 'warn'" :label="domain.dns_synced ? $t('admin.yes') : $t('admin.no')" />
              </td>
              <td class="whitespace-nowrap text-sm text-y-muted">{{ formatDate(domain.registered_at) }}</td>
              <td class="whitespace-nowrap text-sm text-y-muted">{{ formatDate(domain.expires_at) }}</td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    v-if="domain.user"
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('admin.viewUser')"
                    :aria-label="$t('admin.viewUser')"
                    @click="viewUser(domain.user)"
                  >
                    <UserRound :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    v-if="domain.status === 'active'"
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-warn"
                    :title="$t('admin.suspend')"
                    :aria-label="$t('admin.suspend')"
                    @click="updateDomainStatus(domain, 'suspended')"
                  >
                    <Ban :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    v-else-if="domain.status === 'suspended'"
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-mint-text"
                    :title="$t('admin.activate')"
                    :aria-label="$t('admin.activate')"
                    @click="updateDomainStatus(domain, 'active')"
                  >
                    <CircleCheck :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('admin.delete')"
                    :aria-label="$t('admin.delete')"
                    @click="deleteDomain(domain)"
                  >
                    <Trash2 :size="16" :stroke-width="1.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div
        v-if="!loading && pagination.total > 0"
        class="flex flex-col items-center justify-between gap-3 border-t border-y-line px-5 py-3 sm:flex-row"
      >
        <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
          <p class="font-mono text-xs text-y-faint">{{ $t('pagination.showing', pageRange) }}</p>
          <label class="flex items-center gap-2 text-xs text-y-muted">
            <span>{{ $t('pagination.perPage') }}</span>
            <select
              :value="pagination.page_size"
              class="select select-bordered select-sm h-8 min-h-0 font-mono text-xs"
              @change="changePageSize(Number($event.target.value))"
            >
              <option v-for="n in pageSizeOptions" :key="n" :value="n">{{ $t('pagination.perPageOption', { n }) }}</option>
            </select>
          </label>
        </div>
        <div v-if="pagination.total_pages > 1" class="join">
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="pagination.page === 1"
            :aria-label="$t('common.previous')"
            @click="goToPage(pagination.page - 1)"
          >
            <ChevronLeft :size="16" :stroke-width="1.5" />
          </button>
          <template v-for="(page, index) in visiblePages" :key="index">
            <button
              v-if="typeof page === 'number'"
              type="button"
              class="join-item btn btn-ghost btn-sm min-w-[2.125rem] px-2 font-mono font-medium"
              :class="page === pagination.page ? 'bg-y-surface-2 text-y-mint-text' : 'text-y-muted'"
              :aria-current="page === pagination.page ? 'page' : undefined"
              @click="goToPage(page)"
            >
              {{ page }}
            </button>
            <span v-else class="join-item flex h-[2.125rem] items-center px-2 font-mono text-xs text-y-faint">…</span>
          </template>
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="pagination.page === pagination.total_pages"
            :aria-label="$t('common.next')"
            @click="goToPage(pagination.page + 1)"
          >
            <ChevronRight :size="16" :stroke-width="1.5" />
          </button>
        </div>
      </div>
    </section>

    <!-- User Detail Modal -->    <dialog :class="{ 'modal': true, 'modal-open': showUserModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('admin.userInfo') }}</h3>
        <div v-if="selectedUser" class="mt-5 space-y-5">
          <dl class="grid grid-cols-1 gap-x-6 gap-y-4 sm:grid-cols-2">
            <div>
              <dt class="y-label">{{ $t('uiAdmin1.id') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedUser.id }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('uiAdmin1.username') }}</dt>
              <dd class="mt-1 text-sm font-medium text-y-text">{{ selectedUser.username }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('uiAdmin1.email') }}</dt>
              <dd class="mt-1 break-all text-sm text-y-text">{{ selectedUser.email }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.status') }}</dt>
              <dd class="mt-1">
                <span :class="getUserStatusClass(selectedUser.status)">{{ getUserStatusLabel(selectedUser.status) }}</span>
              </dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('admin.registeredAt') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ formatDate(selectedUser.created_at) }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('uiAdmin1.domainQuota') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedUser.domain_quota }}</dd>
            </div>
          </dl>

          <div
            v-if="selectedUser.status === 'active' || selectedUser.status === 'banned' || selectedUser.status === 'frozen'"
            class="y-divider pt-5"
          >
            <button
              v-if="selectedUser.status === 'active'"
              type="button"
              class="btn btn-outline btn-sm gap-1.5 text-y-danger-text"
              @click="updateUserStatus(selectedUser, 'banned')"
            >
              <Ban :size="16" :stroke-width="1.5" />
              {{ $t('admin.suspend') }}
            </button>
            <button
              v-else
              type="button"
              class="btn btn-primary btn-sm gap-1.5"
              @click="updateUserStatus(selectedUser, 'active')"
            >
              <CircleCheck :size="16" :stroke-width="1.5" />
              {{ $t('admin.activate') }}
            </button>
          </div>
        </div>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showUserModal = false">{{ $t('admin.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="showUserModal = false">{{ $t('common.close') }}</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { ArrowLeft, Search, UserRound, Ban, CircleCheck, Trash2, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

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
const pageSizeOptions = [20, 50, 100, 200]
const pageSizeStorageKey = 'admin.rootDomainDomains.pageSize'
const storedPageSize = Number(localStorage.getItem(pageSizeStorageKey))
const pagination = ref({
  page: 1,
  page_size: pageSizeOptions.includes(storedPageSize) ? storedPageSize : 50,
  total: 0,
  total_pages: 0,
})

const pageRange = computed(() => {
  const { page, page_size: size, total } = pagination.value
  const from = total === 0 ? 0 : (page - 1) * size + 1
  return { from, to: Math.min(page * size, total), total }
})

const changePageSize = (size) => {
  if (!pageSizeOptions.includes(size) || size === pagination.value.page_size) return
  localStorage.setItem(pageSizeStorageKey, String(size))
  pagination.value.page_size = size
  pagination.value.page = 1
  fetchDomains()
}

let searchTimeout = null

onMounted(async () => {
  await fetchDomains()
})

const fetchDomains = async () => {
  loading.value = true
  try {
    const params = { page: pagination.value.page, page_size: pagination.value.page_size }
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    const response = await axios.get(`/api/admin/root-domains/${rootDomainId}/domains`, { params })
    domains.value = response.data.domains || []
    rootDomain.value = response.data.root_domain
    if (response.data.pagination) {
      pagination.value = { ...pagination.value, ...response.data.pagination }
    }
  } catch (error) {
    toast.error(error.response?.data?.error || t('admin.operationFailed'))
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    fetchDomains()
  }, 300)
}

const goToPage = (page) => {
  if (page >= 1 && page <= pagination.value.total_pages && page !== pagination.value.page) {
    pagination.value.page = page
    fetchDomains()
  }
}

const visiblePages = computed(() => {
  const current = pagination.value.page
  const total = pagination.value.total_pages
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  if (current <= 4) return [1, 2, 3, 4, 5, '…', total]
  if (current >= total - 3) return [1, '…', total - 4, total - 3, total - 2, total - 1, total]
  return [1, '…', current - 1, current, current + 1, '…', total]
})

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

const statusTone = (status) => {
  const tones = {
    active: 'ok',
    suspended: 'warn',
    expired: 'danger',
  }
  return tones[status] || 'muted'
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
    active: 'y-chip-mint',
    frozen: 'y-chip-warn',
    banned: 'y-chip-danger',
  }
  return classes[status] || 'y-chip'
}

const getUserStatusLabel = (status) => {
  const labels = {
    active: t('uiAdmin1.userStatus.active'),
    frozen: t('uiAdmin1.userStatus.frozen'),
    banned: t('uiAdmin1.userStatus.banned'),
  }
  return labels[status] || status
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString()
}
</script>
