<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('uiAdmin2.eyebrow')"
      :title="$t('uiAdmin2.users.title')"
      :description="$t('uiAdmin2.users.description')"
    >
      <template #actions>
        <span class="font-mono text-sm text-y-faint">{{ $t('uiAdmin2.users.totalCount', { count: pagination.total }) }}</span>
      </template>
    </PageHeader>

    <!-- Stats -->
    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.users.stats.total') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ pagination.total }}</p>
      </div>
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.users.stats.active') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ activeUsers }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.users.stats.admins') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight text-y-text">{{ adminUsers }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
      <div class="y-panel p-5">
        <p class="y-label">{{ $t('uiAdmin2.users.stats.banned') }}</p>
        <p class="mt-2 font-mono text-3xl font-medium tracking-tight" :class="bannedUsers ? 'text-y-danger-text' : 'text-y-text'">{{ bannedUsers }}</p>
        <p class="y-field-hint mt-1">{{ $t('uiAdmin2.onThisPage') }}</p>
      </div>
    </div>

    <!-- Users -->
    <section class="y-panel">
      <div class="y-panel-header">
        <label class="input input-bordered input-sm flex h-10 w-full items-center gap-2 sm:max-w-sm">
          <Search :size="16" :stroke-width="1.5" class="shrink-0 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('adminUsers.searchPlaceholder')"
            class="min-w-0 grow bg-transparent text-sm"
            @input="debouncedSearch"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="btn btn-ghost btn-xs btn-square"
            :aria-label="$t('uiAdmin2.clearSearch')"
            :title="$t('uiAdmin2.clearSearch')"
            @click="clearSearch"
          >
            <X :size="14" :stroke-width="1.5" />
          </button>
        </label>
        <p v-if="searchQuery" class="text-sm text-y-muted">
          {{ $t('adminUsers.searchResults', { count: users.length, query: searchQuery }) }}
        </p>
      </div>

      <YLoader v-if="loading && users.length === 0" :label="$t('uiAdmin2.loading')" />

      <EmptyState
        v-else-if="users.length === 0"
        :title="$t('uiAdmin2.users.emptyTitle')"
        :description="searchQuery ? $t('uiAdmin2.users.emptySearchDescription') : $t('uiAdmin2.users.emptyDescription')"
      >
        <button v-if="searchQuery" type="button" class="btn btn-outline btn-sm" @click="clearSearch">
          {{ $t('uiAdmin2.clearSearch') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>{{ $t('adminUsers.username') }}</th>
              <th>{{ $t('adminUsers.email') }}</th>
              <th>{{ $t('adminUsers.level') }}</th>
              <th>{{ $t('adminUsers.domainQuota') }}</th>
              <th>{{ $t('adminUsers.status') }}</th>
              <th>{{ $t('adminUsers.isAdmin') }}</th>
              <th>{{ $t('adminUsers.registeredAt') }}</th>
              <th class="text-right">{{ $t('adminUsers.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td class="font-mono text-xs text-y-faint">{{ user.id }}</td>
              <td>
                <div class="flex items-center gap-3">
                  <UserAvatar :src="user.avatar" :name="user.username" :size="32" />
                  <div class="min-w-0">
                    <p class="truncate font-medium text-y-text">{{ user.username }}</p>
                    <p class="truncate font-mono text-xs text-y-faint">{{ user.invite_code }}</p>
                  </div>
                </div>
              </td>
              <td class="text-y-muted">{{ user.email }}</td>
              <td>
                <span :class="user.user_level === 'vip' ? 'y-chip-mint' : 'y-chip'">{{ user.user_level }}</span>
              </td>
              <td class="font-mono">{{ user.domain_quota }}</td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="userStatusTone(user.status)" :label="userStatusLabel(user.status)" />
              </td>
              <td>
                <span v-if="user.is_admin" class="y-chip-mint">
                  <ShieldCheck :size="12" :stroke-width="2" />
                  {{ $t('uiAdmin2.users.admin') }}
                </span>
                <span v-else class="text-y-faint">—</span>
              </td>
              <td class="min-w-[7.5rem] text-xs leading-5 text-y-muted">{{ formatDate(user.created_at) }}</td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('adminUsers.editTooltip')"
                    :aria-label="$t('adminUsers.editTooltip')"
                    @click="openEditModal(user)"
                  >
                    <Pencil :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('adminUsers.detailsTooltip')"
                    :aria-label="$t('adminUsers.detailsTooltip')"
                    @click="viewUserDetails(user)"
                  >
                    <Eye :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text disabled:bg-transparent disabled:text-y-faint"
                    :title="$t('adminUsers.deleteTooltip')"
                    :aria-label="$t('adminUsers.deleteTooltip')"
                    :disabled="user.is_admin"
                    @click="confirmDelete(user)"
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

    <!-- Pagination -->
    <div v-if="pagination.total_pages > 1" class="flex flex-col items-center justify-between gap-3 sm:flex-row">
      <p class="font-mono text-xs text-y-faint">{{ $t('uiAdmin2.pageOf', { page: pagination.page, total: pagination.total_pages }) }}</p>
      <div class="join">
        <button
          type="button"
          class="join-item btn btn-outline btn-sm disabled:border-y-line disabled:bg-transparent disabled:text-y-faint"
          :disabled="pagination.page === 1"
          :aria-label="$t('uiAdmin2.previousPage')"
          @click="goToPage(pagination.page - 1)"
        >
          <ChevronLeft :size="16" :stroke-width="1.5" />
        </button>
        <template v-for="(page, index) in visiblePages" :key="index">
          <button
            v-if="typeof page === 'number'"
            type="button"
            class="join-item btn btn-outline btn-sm font-mono"
            :class="{ 'bg-y-surface-2 !text-y-mint-text': page === pagination.page }"
            :aria-current="page === pagination.page ? 'page' : undefined"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
          <button v-else type="button" class="join-item btn btn-outline btn-sm btn-disabled font-mono" tabindex="-1">…</button>
        </template>
        <button
          type="button"
          class="join-item btn btn-outline btn-sm disabled:border-y-line disabled:bg-transparent disabled:text-y-faint"
          :disabled="pagination.page === pagination.total_pages"
          :aria-label="$t('uiAdmin2.nextPage')"
          @click="goToPage(pagination.page + 1)"
        >
          <ChevronRight :size="16" :stroke-width="1.5" />
        </button>
      </div>
    </div>

    <!-- Edit User Modal -->
    <dialog class="modal" :class="{ 'modal-open': showEditModal }">
      <div class="modal-box relative">
        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square absolute right-3 top-3"
          :aria-label="$t('uiAdmin2.close')"
          @click="closeEditModal"
        >
          <X :size="16" :stroke-width="1.5" />
        </button>
        <h3 class="text-lg font-semibold text-y-text">{{ $t('uiAdmin2.users.editTitle') }}</h3>
        <p class="mt-1 text-sm text-y-muted">{{ $t('uiAdmin2.users.editDescription') }}</p>

        <form class="mt-6 space-y-4" @submit.prevent="handleUpdateUser">
          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.users.fields.username') }}</span>
            <input v-model="editForm.username" type="text" class="input input-bordered w-full" required minlength="3" />
          </label>

          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.users.fields.email') }}</span>
            <input v-model="editForm.email" type="email" class="input input-bordered w-full" required />
          </label>

          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.users.fields.newPassword') }}</span>
            <input
              v-model="editForm.password"
              type="password"
              class="input input-bordered w-full"
              :placeholder="$t('adminUsers.passwordPlaceholder')"
              minlength="6"
              autocomplete="new-password"
            />
            <span class="y-field-hint">{{ $t('uiAdmin2.users.fields.passwordHint') }}</span>
          </label>

          <label class="y-field">
            <span class="y-field-label">{{ $t('uiAdmin2.users.fields.status') }}</span>
            <select v-model="editForm.status" class="select select-bordered w-full">
              <option value="active">{{ $t('uiAdmin2.users.status.active') }}</option>
              <option value="frozen">{{ $t('uiAdmin2.users.status.frozen') }}</option>
              <option value="banned">{{ $t('uiAdmin2.users.status.banned') }}</option>
            </select>
          </label>

          <label class="flex cursor-pointer items-start gap-3 rounded-md border border-y-line px-4 py-3">
            <input v-model="editForm.is_admin" type="checkbox" class="checkbox checkbox-sm mt-0.5" />
            <span class="flex flex-col gap-0.5">
              <span class="y-field-label">{{ $t('uiAdmin2.users.fields.isAdmin') }}</span>
              <span class="y-field-hint">{{ $t('uiAdmin2.users.fields.isAdminHint') }}</span>
            </span>
          </label>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeEditModal">{{ $t('uiAdmin2.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              {{ $t('uiAdmin2.saveChanges') }}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeEditModal">
        <button>{{ $t('uiAdmin2.close') }}</button>
      </form>
    </dialog>

    <!-- User Details Modal -->
    <dialog class="modal" :class="{ 'modal-open': showDetailsModal }">
      <div class="modal-box relative max-w-2xl">
        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square absolute right-3 top-3"
          :aria-label="$t('uiAdmin2.close')"
          @click="showDetailsModal = false"
        >
          <X :size="16" :stroke-width="1.5" />
        </button>
        <h3 class="text-lg font-semibold text-y-text">{{ $t('uiAdmin2.users.detailsTitle') }}</h3>

        <div v-if="selectedUser" class="mt-6 space-y-5">
          <div class="flex items-center gap-4">
            <UserAvatar :src="selectedUser.avatar" :name="selectedUser.username" :size="48" />
            <div class="min-w-0">
              <p class="truncate text-base font-semibold text-y-text">{{ selectedUser.username }}</p>
              <p class="truncate text-sm text-y-muted">{{ selectedUser.email }}</p>
            </div>
          </div>

          <dl class="grid grid-cols-1 gap-x-6 gap-y-4 rounded-md border border-y-line p-4 sm:grid-cols-2">
            <div>
              <dt class="y-label">{{ $t('adminUsers.userId') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedUser.id }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.inviteCode') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedUser.invite_code }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.userLevel') }}</dt>
              <dd class="mt-1">
                <span :class="selectedUser.user_level === 'vip' ? 'y-chip-mint' : 'y-chip'">{{ selectedUser.user_level }}</span>
              </dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.domainQuota') }}</dt>
              <dd class="mt-1 font-mono text-sm text-y-text">{{ selectedUser.domain_quota }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.status') }}</dt>
              <dd class="mt-1">
                <StatusDot :tone="userStatusTone(selectedUser.status)" :label="userStatusLabel(selectedUser.status)" />
              </dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.isAdmin') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ selectedUser.is_admin ? $t('uiAdmin2.yes') : $t('uiAdmin2.no') }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.registeredAt') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ formatDate(selectedUser.created_at) }}</dd>
            </div>
            <div>
              <dt class="y-label">{{ $t('adminUsers.lastLogin') }}</dt>
              <dd class="mt-1 text-sm text-y-text">{{ selectedUser.last_login_at ? formatDate(selectedUser.last_login_at) : $t('adminUsers.never') }}</dd>
            </div>
          </dl>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDetailsModal = false">{{ $t('uiAdmin2.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showDetailsModal = false">
        <button>{{ $t('uiAdmin2.close') }}</button>
      </form>
    </dialog>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { Search, X, Pencil, Eye, Trash2, ShieldCheck, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import UserAvatar from '../components/ui/UserAvatar.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()

const users = ref([])
const loading = ref(true)
const searchQuery = ref('')
const showDetailsModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)
const selectedUser = ref(null)
const editingUser = ref(null)
const pagination = ref({
  page: 1,
  page_size: 20,
  total: 0,
  total_pages: 0
})
let searchTimeout = null

const editForm = ref({
  username: '',
  email: '',
  password: '',
  status: 'active',
  is_admin: false,
})

const activeUsers = computed(() => {
  return users.value.filter(u => u.status === 'active').length
})

const adminUsers = computed(() => {
  return users.value.filter(u => u.is_admin).length
})

const bannedUsers = computed(() => {
  return users.value.filter(u => u.status === 'banned').length
})

const fetchUsers = async (search = '') => {
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.page_size
    }
    if (search) {
      params.search = search
    }
    const response = await axios.get('/api/admin/users', { params })
    users.value = response.data.users || []
    if (response.data.pagination) {
      pagination.value = response.data.pagination
    }
  } catch (error) {
    toast.error(t('adminUsers.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const debouncedSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    fetchUsers(searchQuery.value)
  }, 500)
}

const clearSearch = () => {
  searchQuery.value = ''
  pagination.value.page = 1
  fetchUsers()
}

const goToPage = (page) => {
  if (page >= 1 && page <= pagination.value.total_pages) {
    pagination.value.page = page
    fetchUsers(searchQuery.value)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const visiblePages = computed(() => {
  const current = pagination.value.page
  const total = pagination.value.total_pages
  const pages = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages.filter(p => p !== '...' || pages.indexOf(p) === pages.lastIndexOf(p))
})

const viewUserDetails = (user) => {
  selectedUser.value = user
  showDetailsModal.value = true
}

const openEditModal = (user) => {
  editingUser.value = user
  editForm.value = {
    username: user.username,
    email: user.email,
    password: '',
    status: user.status,
    is_admin: user.is_admin,
  }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingUser.value = null
  editForm.value = { username: '', email: '', password: '', status: 'active', is_admin: false }
}

const handleUpdateUser = async () => {
  submitting.value = true
  try {
    const payload = {
      username: editForm.value.username,
      email: editForm.value.email,
      is_admin: editForm.value.is_admin,
      status: editForm.value.status,
    }
    if (editForm.value.password) {
      payload.password = editForm.value.password
    }
    await axios.put(`/api/admin/users/${editingUser.value.id}`, payload)
    toast.success(t('adminUsers.updateSuccess'))
    closeEditModal()
    await fetchUsers()
  } catch (error) {
    toast.error(error.response?.data?.error || t('adminUsers.updateFailed'))
  } finally {
    submitting.value = false
  }
}

const confirmDelete = (user) => {
  if (user.is_admin) {
    toast.error(t('adminUsers.cannotDeleteAdmin'))
    return
  }

  if (confirm(t('adminUsers.confirmDelete', { username: user.username }))) {
    deleteUser(user)
  }
}

const deleteUser = async (user) => {
  try {
    await axios.delete(`/api/admin/users/${user.id}`)
    toast.success(t('adminUsers.deleteSuccess'))
    await fetchUsers(searchQuery.value)
  } catch (error) {
    toast.error(t('adminUsers.deleteFailed') + ': ' + (error.response?.data?.error || error.message))
  }
}

const userStatusTone = (status) => ({ active: 'ok', frozen: 'warn', banned: 'danger' }[status] || 'muted')

const userStatusLabel = (status) =>
  ['active', 'frozen', 'banned'].includes(status) ? t(`uiAdmin2.users.status.${status}`) : status

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  fetchUsers()
})
</script>
