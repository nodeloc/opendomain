<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">{{ $t('adminUsers.title') }}</h1>
      <div class="text-sm opacity-70">{{ $t('adminUsers.totalCount', { count: pagination.total }) }}</div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminUsers.totalUsers') }}</div>
        <div class="stat-value text-primary">{{ pagination.total }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminUsers.activeUsers') }}</div>
        <div class="stat-value text-success">{{ activeUsers }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminUsers.admins') }}</div>
        <div class="stat-value text-warning">{{ adminUsers }}</div>
      </div>
      <div class="stat bg-base-100 rounded-lg shadow">
        <div class="stat-title">{{ $t('adminUsers.bannedUsers') }}</div>
        <div class="stat-value text-error">{{ bannedUsers }}</div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <div class="relative">
        <div class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="$t('adminUsers.searchPlaceholder')"
          class="input input-bordered w-full pl-12 pr-10"
          @input="debouncedSearch"
        />
        <button
          v-if="searchQuery"
          @click="clearSearch"
          class="absolute inset-y-0 right-0 flex items-center pr-4 text-gray-400 hover:text-gray-600"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <p v-if="searchQuery" class="text-sm text-gray-500 mt-2">
        {{ $t('adminUsers.searchResults', { count: users.length, query: searchQuery }) }}
      </p>
    </div>

    <!-- Users Table -->
    <div class="overflow-x-auto bg-base-100 rounded-lg shadow">
      <table class="table w-full">
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
            <th>{{ $t('adminUsers.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>
              <div class="font-bold">{{ user.username }}</div>
              <div class="text-xs opacity-50">{{ $t('adminUsers.inviteCode') }}: {{ user.invite_code }}</div>
            </td>
            <td>{{ user.email }}</td>
            <td>
              <span :class="['px-2 py-0.5 rounded text-xs font-medium',
                user.user_level === 'vip' ? 'bg-yellow-100 text-yellow-800' : 'bg-gray-100 text-gray-800']">
                {{ user.user_level }}
              </span>
            </td>
            <td>{{ user.domain_quota }}</td>
            <td>
              <span :class="['px-2 py-0.5 rounded text-xs font-medium',
                user.status === 'active' ? 'bg-green-100 text-green-800' :
                user.status === 'banned' ? 'bg-red-100 text-red-800' : 'bg-gray-100 text-gray-800']">
                {{ user.status }}
              </span>
            </td>
            <td>
              <span v-if="user.is_admin" class="px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                Admin
              </span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td>
              <div class="flex gap-1">
                <button @click="openEditModal(user)" class="btn btn-sm btn-ghost" :title="$t('adminUsers.editTooltip')">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>
                <button @click="viewUserDetails(user)" class="btn btn-sm btn-ghost" :title="$t('adminUsers.detailsTooltip')">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                    <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                  </svg>
                </button>
                <button @click="confirmDelete(user)" class="btn btn-sm btn-ghost text-error" :title="$t('adminUsers.deleteTooltip')" :disabled="user.is_admin">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="users.length === 0" class="text-center py-8 text-gray-500">
        {{ $t('adminUsers.noData') }}
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="pagination.total_pages > 1" class="flex justify-center mt-6">
      <div class="join">
        <button
          class="join-item btn"
          :disabled="pagination.page === 1"
          @click="goToPage(pagination.page - 1)"
        >
          «
        </button>
        <template v-for="(page, index) in visiblePages" :key="index">
          <button
            v-if="typeof page === 'number'"
            class="join-item btn"
            :class="{ 'btn-active': page === pagination.page }"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
          <button v-else class="join-item btn btn-disabled">
            ...
          </button>
        </template>
        <button
          class="join-item btn"
          :disabled="pagination.page === pagination.total_pages"
          @click="goToPage(pagination.page + 1)"
        >
          »
        </button>
      </div>
    </div>

    <!-- Edit User Modal -->
    <dialog class="modal" :class="{ 'modal-open': showEditModal }">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">{{ $t('adminUsers.editUser') }}</h3>
        <form @submit.prevent="handleUpdateUser" class="space-y-4">
          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('adminUsers.username') }}</span></label>
            <input v-model="editForm.username" type="text" class="input input-bordered" required minlength="3" />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('adminUsers.email') }}</span></label>
            <input v-model="editForm.email" type="email" class="input input-bordered" required />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('adminUsers.newPassword') }}</span></label>
            <input v-model="editForm.password" type="password" class="input input-bordered" :placeholder="$t('adminUsers.passwordPlaceholder')" minlength="6" />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">{{ $t('adminUsers.status') }}</span></label>
            <select v-model="editForm.status" class="select select-bordered">
              <option value="active">Active</option>
              <option value="frozen">Frozen</option>
              <option value="banned">Banned</option>
            </select>
          </div>

          <div class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">{{ $t('adminUsers.isAdmin') }}</span>
              <input v-model="editForm.is_admin" type="checkbox" class="checkbox" />
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeEditModal" class="btn">{{ $t('common.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ $t('common.save') }}</span>
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeEditModal">
        <button>close</button>
      </form>
    </dialog>

    <!-- User Details Modal -->
    <dialog class="modal" :class="{ 'modal-open': showDetailsModal }">
      <div class="modal-box max-w-2xl">
        <h3 class="font-bold text-lg mb-4">{{ $t('adminUsers.userDetails') }}</h3>
        <div v-if="selectedUser" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.userId') }}</span></label>
              <div class="text-lg">{{ selectedUser.id }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.username') }}</span></label>
              <div class="text-lg">{{ selectedUser.username }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.email') }}</span></label>
              <div class="text-lg">{{ selectedUser.email }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.inviteCode') }}</span></label>
              <div class="text-lg font-mono">{{ selectedUser.invite_code }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.userLevel') }}</span></label>
              <div class="text-lg">{{ selectedUser.user_level }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.domainQuota') }}</span></label>
              <div class="text-lg">{{ selectedUser.domain_quota }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.status') }}</span></label>
              <div class="text-lg">{{ selectedUser.status }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.isAdmin') }}</span></label>
              <div class="text-lg">{{ selectedUser.is_admin ? $t('common.yes') : $t('common.no') }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.registeredAt') }}</span></label>
              <div class="text-lg">{{ formatDate(selectedUser.created_at) }}</div>
            </div>
            <div>
              <label class="label"><span class="label-text font-semibold">{{ $t('adminUsers.lastLogin') }}</span></label>
              <div class="text-lg">{{ selectedUser.last_login_at ? formatDate(selectedUser.last_login_at) : $t('adminUsers.never') }}</div>
            </div>
          </div>
        </div>
        <div class="modal-action">
          <button @click="showDetailsModal = false" class="btn">{{ $t('common.close') }}</button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showDetailsModal = false">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const toast = useToast()

const users = ref([])
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

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  fetchUsers()
})
</script>
