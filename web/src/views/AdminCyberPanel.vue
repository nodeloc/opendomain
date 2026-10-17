<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('admin.cyberPanel')"
      :description="$t('uiAdmin3.cp.description')"
    >
      <template #actions>
        <button type="button" class="btn btn-primary btn-sm gap-1.5" @click="openCreateServer">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('admin.cpAddServer') }}
        </button>
      </template>
    </PageHeader>

    <!-- Tabs -->
    <div role="tablist" class="tabs tabs-boxed w-full sm:w-fit">
      <button
        type="button"
        role="tab"
        class="tab flex-1 gap-2 sm:flex-none"
        :class="{ 'tab-active': activeTab === 'servers' }"
        :aria-selected="activeTab === 'servers'"
        @click="activeTab = 'servers'"
      >
        <Server :size="16" :stroke-width="1.5" />
        {{ $t('admin.cpServers') }}
        <span class="font-mono text-xs text-y-faint">{{ servers.length }}</span>
      </button>
      <button
        type="button"
        role="tab"
        class="tab flex-1 gap-2 sm:flex-none"
        :class="{ 'tab-active': activeTab === 'accounts' }"
        :aria-selected="activeTab === 'accounts'"
        @click="activeTab = 'accounts'; fetchAccounts()"
      >
        <Users :size="16" :stroke-width="1.5" />
        {{ $t('admin.cpAccounts') }}
        <span class="font-mono text-xs text-y-faint">{{ accountTotal }}</span>
      </button>
    </div>

    <!-- Servers -->
    <section v-if="activeTab === 'servers'" class="y-panel">
      <YLoader v-if="serversLoading" :label="$t('common.loading')" />

      <EmptyState
        v-else-if="servers.length === 0"
        :title="$t('admin.cpNoServers')"
        :description="$t('admin.cpNoServersDesc')"
      >
        <button type="button" class="btn btn-primary btn-sm gap-1.5" @click="openCreateServer">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('admin.cpAddServer') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th>{{ $t('admin.cpServerName') }}</th>
              <th>{{ $t('admin.cpServerUrl') }}</th>
              <th>{{ $t('admin.cpAdminUser') }}</th>
              <th>{{ $t('admin.cpPackage') }}</th>
              <th>{{ $t('admin.cpAccUsage') }}</th>
              <th>{{ $t('admin.status') }}</th>
              <th class="text-right">{{ $t('admin.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="server in servers" :key="server.id">
              <td>
                <div class="flex flex-wrap items-center gap-2">
                  <span class="whitespace-nowrap font-medium text-y-text">{{ server.name }}</span>
                  <span v-if="server.is_default" class="y-chip-mint">{{ $t('uiAdmin3.cp.defaultBadge') }}</span>
                </div>
              </td>
              <td class="whitespace-nowrap font-mono text-sm text-y-muted">{{ server.url }}</td>
              <td class="whitespace-nowrap font-mono text-sm">{{ server.admin_user }}</td>
              <td class="whitespace-nowrap font-mono text-sm">{{ server.package_name }}</td>
              <td>
                <div class="flex min-w-[7rem] flex-col gap-1.5">
                  <span class="whitespace-nowrap font-mono text-sm">
                    {{ server.current_accounts }}<span class="text-y-faint"> / {{ server.max_accounts === 0 ? '∞' : server.max_accounts }}</span>
                  </span>
                  <div
                    v-if="server.max_accounts > 0"
                    class="h-1 w-20 overflow-hidden rounded-full bg-y-surface-2"
                    role="progressbar"
                    :aria-valuenow="server.current_accounts"
                    aria-valuemin="0"
                    :aria-valuemax="server.max_accounts"
                  >
                    <div class="h-full rounded-full" :class="usageBarClass(server)" :style="{ width: usagePercent(server) + '%' }"></div>
                  </div>
                </div>
              </td>
              <td class="whitespace-nowrap">
                <StatusDot :tone="server.is_active ? 'ok' : 'muted'" :label="server.is_active ? $t('admin.active') : $t('admin.inactive')" />
              </td>
              <td class="whitespace-nowrap">
                <div class="flex justify-end gap-1">
                  <button type="button" class="btn btn-ghost btn-xs gap-1" :disabled="testingId === server.id" @click="testServer(server)">
                    <span v-if="testingId === server.id" class="loading loading-spinner loading-xs"></span>
                    <PlugZap v-else :size="14" :stroke-width="1.5" />
                    {{ $t('admin.cpTest') }}
                  </button>
                  <button type="button" class="btn btn-outline btn-xs gap-1" :disabled="openingPanelId === server.id" @click="openServerPanel(server)">
                    <span v-if="openingPanelId === server.id" class="loading loading-spinner loading-xs"></span>
                    <ExternalLink v-else :size="14" :stroke-width="1.5" />
                    {{ $t('admin.cpOpenPanel') }}
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-xs btn-square"
                    :title="$t('admin.edit')"
                    :aria-label="$t('admin.edit')"
                    @click="openEditServer(server)"
                  >
                    <Pencil :size="14" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-xs btn-square text-y-danger-text"
                    :title="$t('admin.delete')"
                    :aria-label="$t('admin.delete')"
                    @click="confirmDeleteServer(server)"
                  >
                    <Trash2 :size="14" :stroke-width="1.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Accounts -->
    <section v-if="activeTab === 'accounts'" class="y-panel">
      <div class="y-panel-header">
        <h2 class="text-base font-semibold text-y-text">{{ $t('admin.cpAccounts') }}</h2>
        <div class="flex w-full flex-col gap-2 sm:w-auto sm:flex-row">
          <select
            v-model="accountFilter.status"
            class="select select-bordered select-sm"
            :aria-label="$t('admin.status')"
            @change="fetchAccounts"
          >
            <option value="">{{ $t('admin.cpAllStatuses') }}</option>
            <option v-for="status in accountStatuses" :key="status" :value="status">{{ accountStatusLabel(status) }}</option>
          </select>
          <select
            v-model="accountFilter.server_id"
            class="select select-bordered select-sm"
            :aria-label="$t('admin.cpServer')"
            @change="fetchAccounts"
          >
            <option value="">{{ $t('admin.cpAllServers') }}</option>
            <option v-for="s in servers" :key="s.id" :value="s.id">{{ s.name }}</option>
          </select>
        </div>
      </div>

      <YLoader v-if="accountsLoading" :label="$t('common.loading')" />

      <EmptyState
        v-else-if="accounts.length === 0"
        :title="$t('admin.cpNoAccounts')"
        :description="$t('uiAdmin3.cp.noAccountsDesc')"
      />

      <div v-else>
        <div class="overflow-x-auto">
          <table class="y-table">
            <thead>
              <tr>
                <th>{{ $t('admin.cpDomain') }}</th>
                <th>{{ $t('admin.cpCpUser') }}</th>
                <th>{{ $t('admin.cpServer') }}</th>
                <th>{{ $t('admin.cpOwner') }}</th>
                <th>{{ $t('admin.status') }}</th>
                <th>{{ $t('admin.cpCreatedAt') }}</th>
                <th class="text-right">{{ $t('admin.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="acc in accounts" :key="acc.id">
                <td><span class="y-domain whitespace-nowrap">{{ acc.domain?.full_domain || acc.domain_id }}</span></td>
                <td class="whitespace-nowrap font-mono text-sm">{{ acc.cp_username }}</td>
                <td class="whitespace-nowrap">{{ acc.server?.name || acc.server_id }}</td>
                <td class="whitespace-nowrap text-y-muted">{{ acc.user?.username || acc.user_id }}</td>
                <td class="whitespace-nowrap">
                  <div class="flex items-center gap-2">
                    <StatusDot :tone="statusTone(acc.status)" :label="accountStatusLabel(acc.status)" />
                    <span v-if="acc.error_msg" class="tooltip tooltip-left" :data-tip="acc.error_msg">
                      <AlertTriangle :size="14" :stroke-width="1.5" class="text-y-danger-text" :aria-label="$t('uiAdmin3.cp.lastError')" />
                    </span>
                  </div>
                </td>
                <td class="whitespace-nowrap font-mono text-xs text-y-muted">{{ formatDate(acc.created_at) }}</td>
                <td class="whitespace-nowrap">
                  <div class="flex justify-end gap-1">
                    <button
                      v-if="acc.status === 'active'"
                      type="button"
                      class="btn btn-ghost btn-xs gap-1"
                      @click="adminSuspend(acc)"
                    >
                      <PauseCircle :size="14" :stroke-width="1.5" />
                      {{ $t('admin.suspend') }}
                    </button>
                    <button
                      v-if="acc.status === 'suspended'"
                      type="button"
                      class="btn btn-ghost btn-xs gap-1"
                      @click="adminUnsuspend(acc)"
                    >
                      <PlayCircle :size="14" :stroke-width="1.5" />
                      {{ $t('admin.unsuspend') }}
                    </button>
                    <button
                      v-if="acc.status !== 'terminated'"
                      type="button"
                      class="btn btn-ghost btn-xs gap-1 text-y-danger-text"
                      @click="confirmTerminateAccount(acc)"
                    >
                      <Ban :size="14" :stroke-width="1.5" />
                      {{ $t('admin.cpTerminate') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="accountTotal > accountPageSize" class="flex items-center justify-between gap-3 border-t border-y-line px-5 py-3">
          <span class="font-mono text-xs text-y-faint">
            {{ $t('uiAdmin3.pageOf', { page: accountPage, total: Math.ceil(accountTotal / accountPageSize) }) }}
          </span>
          <div class="flex items-center gap-1">
            <button
              type="button"
              class="btn btn-ghost btn-sm btn-square"
              :disabled="accountPage === 1"
              :aria-label="$t('uiAdmin3.prevPage')"
              @click="accountPage--; fetchAccounts()"
            >
              <ChevronLeft :size="16" :stroke-width="1.5" />
            </button>
            <button
              type="button"
              class="btn btn-ghost btn-sm btn-square"
              :disabled="accountPage >= Math.ceil(accountTotal / accountPageSize)"
              :aria-label="$t('uiAdmin3.nextPage')"
              @click="accountPage++; fetchAccounts()"
            >
              <ChevronRight :size="16" :stroke-width="1.5" />
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Server modal -->
    <dialog v-if="showServerModal" class="modal modal-open">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="text-lg font-semibold text-y-text">
          {{ editingServer ? $t('admin.cpEditServer') : $t('admin.cpAddServer') }}
        </h3>
        <p v-if="editingServer" class="mt-1 text-sm text-y-muted">{{ editingServer.name }}</p>

        <div class="mt-5 grid grid-cols-1 gap-4 md:grid-cols-2">
          <label class="y-field">
            <span class="y-field-label">{{ $t('admin.cpServerName') }} <span class="text-y-faint" :title="$t('uiAdmin3.cp.required')">*</span></span>
            <input v-model="serverForm.name" type="text" class="input input-bordered w-full" :placeholder="$t('admin.cpServerNamePlaceholder')" />
          </label>
          <label class="y-field">
            <span class="y-field-label">{{ $t('admin.cpServerUrl') }} <span class="text-y-faint" :title="$t('uiAdmin3.cp.required')">*</span></span>
            <input v-model="serverForm.url" type="text" class="input input-bordered w-full font-mono" placeholder="http://1.2.3.4:8090" autocomplete="off" spellcheck="false" />
          </label>
          <label class="y-field">
            <span class="y-field-label">{{ $t('admin.cpAdminUser') }} <span class="text-y-faint" :title="$t('uiAdmin3.cp.required')">*</span></span>
            <input v-model="serverForm.admin_user" type="text" class="input input-bordered w-full font-mono" placeholder="admin" autocomplete="off" spellcheck="false" />
          </label>
          <label class="y-field">
            <span class="y-field-label">{{ $t('admin.cpAdminPass') }} <span class="text-y-faint" :title="$t('uiAdmin3.cp.required')">*</span></span>
            <input v-model="serverForm.admin_pass" type="password" class="input input-bordered w-full font-mono" placeholder="••••••••" autocomplete="new-password" />
            <span v-if="editingServer" class="y-field-hint">{{ $t('admin.cpPassHint') }}</span>
          </label>
          <label class="y-field">
            <span class="y-field-label">{{ $t('admin.cpPackage') }} <span class="text-y-faint" :title="$t('uiAdmin3.cp.required')">*</span></span>
            <input v-model="serverForm.package_name" type="text" class="input input-bordered w-full font-mono" placeholder="Default" spellcheck="false" />
          </label>
          <label class="y-field">
            <span class="y-field-label">{{ $t('admin.cpMaxAccounts') }}</span>
            <input v-model.number="serverForm.max_accounts" type="number" min="0" class="input input-bordered w-full font-mono" :placeholder="$t('uiAdmin3.cp.unlimitedPlaceholder')" />
          </label>
          <label class="y-field md:col-span-2">
            <span class="y-field-label">{{ $t('admin.description') }}</span>
            <input v-model="serverForm.description" type="text" class="input input-bordered w-full" />
          </label>
          <label class="flex cursor-pointer items-center justify-between gap-4 rounded-md border border-y-line bg-y-surface-2 px-4 py-3">
            <span class="text-sm font-medium text-y-text">{{ $t('admin.active') }}</span>
            <input v-model="serverForm.is_active" type="checkbox" class="toggle toggle-primary" />
          </label>
          <label class="flex cursor-pointer items-center justify-between gap-4 rounded-md border border-y-line bg-y-surface-2 px-4 py-3">
            <span class="text-sm font-medium text-y-text">{{ $t('admin.cpDefault') }}</span>
            <input v-model="serverForm.is_default" type="checkbox" class="toggle toggle-primary" />
          </label>
        </div>

        <div v-if="serverModalError" role="alert" class="alert alert-error mt-4 text-sm">
          <AlertTriangle :size="16" :stroke-width="1.5" />
          <span>{{ serverModalError }}</span>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showServerModal = false">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-primary" :disabled="serverSaving" @click="saveServer">
            <span v-if="serverSaving" class="loading loading-spinner loading-sm"></span>
            {{ $t('common.save') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showServerModal = false"></form>
    </dialog>

    <!-- Confirm delete server -->
    <dialog v-if="showDeleteServerConfirm" class="modal modal-open">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('admin.cpDeleteServerTitle') }}</h3>
        <p class="mt-2 text-sm leading-6 text-y-muted">{{ $t('admin.cpDeleteServerConfirm', { name: deletingServer?.name }) }}</p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDeleteServerConfirm = false">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-error" :disabled="serverSaving" @click="deleteServer">
            <span v-if="serverSaving" class="loading loading-spinner loading-sm"></span>
            {{ $t('admin.delete') }}
          </button>
        </div>
      </div>
    </dialog>

    <!-- Confirm terminate account -->
    <dialog v-if="showTerminateConfirm" class="modal modal-open">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('admin.cpTerminateTitle') }}</h3>
        <p class="mt-2 text-sm leading-6 text-y-muted">{{ $t('admin.cpTerminateConfirm', { domain: terminatingAccount?.domain?.full_domain }) }}</p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showTerminateConfirm = false">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-error" :disabled="accountActionLoading" @click="terminateAccount">
            <span v-if="accountActionLoading" class="loading loading-spinner loading-sm"></span>
            {{ $t('admin.cpTerminate') }}
          </button>
        </div>
      </div>
    </dialog>

    <!-- Toast -->
    <div v-if="toast.show" class="toast toast-top toast-end z-50">
      <div role="status" class="alert text-sm" :class="toast.type === 'success' ? 'alert-success' : 'alert-error'">
        <CheckCircle2 v-if="toast.type === 'success'" :size="16" :stroke-width="1.5" class="text-y-mint-text" />
        <AlertTriangle v-else :size="16" :stroke-width="1.5" />
        <span>{{ toast.message }}</span>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import axios from '../utils/axios'
import { useI18n } from 'vue-i18n'
import {
  Plus, Server, Users, PlugZap, ExternalLink, Pencil, Trash2, AlertTriangle,
  PauseCircle, PlayCircle, Ban, ChevronLeft, ChevronRight, CheckCircle2,
} from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()

// ── State ────────────────────────────────────────────────────────────
const activeTab = ref('servers')

// Servers
const servers = ref([])
const serversLoading = ref(false)
const showServerModal = ref(false)
const editingServer = ref(null)
const serverSaving = ref(false)
const serverModalError = ref('')
const testingId = ref(null)
const openingPanelId = ref(null)
const showDeleteServerConfirm = ref(false)
const deletingServer = ref(null)

const defaultServerForm = () => ({
  name: '',
  url: '',
  admin_user: '',
  admin_pass: '',
  package_name: 'Default',
  is_active: true,
  is_default: false,
  max_accounts: 0,
  description: '',
})
const serverForm = ref(defaultServerForm())

// Accounts
const accounts = ref([])
const accountsLoading = ref(false)
const accountTotal = ref(0)
const accountPage = ref(1)
const accountPageSize = 20
const accountFilter = ref({ status: '', server_id: '' })
const showTerminateConfirm = ref(false)
const terminatingAccount = ref(null)
const accountActionLoading = ref(false)

// Toast
const toast = ref({ show: false, type: 'success', message: '' })

// ── Toast ────────────────────────────────────────────────────────────
const showToast = (message, type = 'success') => {
  toast.value = { show: true, type, message }
  setTimeout(() => (toast.value.show = false), 3000)
}

// ── Servers ──────────────────────────────────────────────────────────
const fetchServers = async () => {
  serversLoading.value = true
  try {
    const res = await axios.get('/api/admin/cyberpanel/servers')
    servers.value = res.data.servers || []
  } catch (e) {
    showToast(e.response?.data?.error || 'Failed to load servers', 'error')
  } finally {
    serversLoading.value = false
  }
}

const openCreateServer = () => {
  editingServer.value = null
  serverForm.value = defaultServerForm()
  serverModalError.value = ''
  showServerModal.value = true
}

const openEditServer = (server) => {
  editingServer.value = server
  serverForm.value = {
    name: server.name,
    url: server.url,
    admin_user: server.admin_user,
    admin_pass: '',
    package_name: server.package_name,
    is_active: server.is_active,
    is_default: server.is_default,
    max_accounts: server.max_accounts,
    description: server.description || '',
  }
  serverModalError.value = ''
  showServerModal.value = true
}

const saveServer = async () => {
  serverModalError.value = ''
  if (!serverForm.value.name || !serverForm.value.url || !serverForm.value.admin_user || !serverForm.value.package_name) {
    serverModalError.value = t('admin.cpRequiredFields')
    return
  }
  if (!editingServer.value && !serverForm.value.admin_pass) {
    serverModalError.value = t('admin.cpPassRequired')
    return
  }

  serverSaving.value = true
  try {
    const payload = { ...serverForm.value }
    if (!payload.admin_pass) delete payload.admin_pass

    if (editingServer.value) {
      await axios.put(`/api/admin/cyberpanel/servers/${editingServer.value.id}`, payload)
      showToast(t('admin.cpServerUpdated'))
    } else {
      await axios.post('/api/admin/cyberpanel/servers', payload)
      showToast(t('admin.cpServerCreated'))
    }
    showServerModal.value = false
    await fetchServers()
  } catch (e) {
    serverModalError.value = e.response?.data?.error || 'Operation failed'
  } finally {
    serverSaving.value = false
  }
}

const confirmDeleteServer = (server) => {
  deletingServer.value = server
  showDeleteServerConfirm.value = true
}

const deleteServer = async () => {
  if (!deletingServer.value) return
  serverSaving.value = true
  try {
    await axios.delete(`/api/admin/cyberpanel/servers/${deletingServer.value.id}`)
    showToast(t('admin.cpServerDeleted'))
    showDeleteServerConfirm.value = false
    await fetchServers()
  } catch (e) {
    showToast(e.response?.data?.error || 'Delete failed', 'error')
  } finally {
    serverSaving.value = false
  }
}

const testServer = async (server) => {
  testingId.value = server.id
  try {
    const res = await axios.post(`/api/admin/cyberpanel/servers/${server.id}/test`)
    if (res.data.success) {
      showToast(t('admin.cpTestSuccess'))
    } else {
      showToast(res.data.error || 'Connection failed', 'error')
    }
  } catch (e) {
    showToast(e.response?.data?.error || 'Test failed', 'error')
  } finally {
    testingId.value = null
  }
}

const cpFormLogin = (username, password, loginUrl) => {
  // 打开命名窗口，先 POST 到 /api/loginAPI（免 CSRF）建立 session cookie，
  // 再导航到 /dashboard/，浏览器会携带已设置的 cookie
  const winName = 'cpanel_' + Date.now()
  const cpWindow = window.open('about:blank', winName)
  const form = document.createElement('form')
  form.method = 'POST'
  form.action = `${loginUrl}/api/loginAPI`
  form.target = winName
  ;[['username', username], ['password', password]].forEach(([name, value]) => {
    const input = document.createElement('input')
    input.type = 'hidden'
    input.name = name
    input.value = value
    form.appendChild(input)
  })
  document.body.appendChild(form)
  form.submit()
  document.body.removeChild(form)
  // 等待 loginAPI 响应后（cookie 已写入浏览器），跳转到面板
  setTimeout(() => {
    if (cpWindow && !cpWindow.closed) {
      cpWindow.location.href = `${loginUrl}/dashboard/`
    }
  }, 1500)
}

const openServerPanel = async (server) => {
  openingPanelId.value = server.id
  try {
    const res = await axios.get(`/api/admin/cyberpanel/servers/${server.id}/autologin`)
    const { username, password, login_url } = res.data
    cpFormLogin(username, password, login_url)
  } catch (e) {
    showToast(e.response?.data?.error || 'Login failed', 'error')
  } finally {
    openingPanelId.value = null
  }
}

// ── Accounts ─────────────────────────────────────────────────────────
const fetchAccounts = async () => {
  accountsLoading.value = true
  try {
    const params = {
      page: accountPage.value,
      page_size: accountPageSize,
    }
    if (accountFilter.value.status) params.status = accountFilter.value.status
    if (accountFilter.value.server_id) params.server_id = accountFilter.value.server_id

    const res = await axios.get('/api/admin/cyberpanel/accounts', { params })
    accounts.value = res.data.accounts || []
    accountTotal.value = res.data.total || 0
  } catch (e) {
    showToast(e.response?.data?.error || 'Failed to load accounts', 'error')
  } finally {
    accountsLoading.value = false
  }
}

const adminSuspend = async (acc) => {
  accountActionLoading.value = true
  try {
    await axios.post(`/api/admin/cyberpanel/accounts/${acc.id}/suspend`)
    showToast(t('admin.cpSuspended'))
    await fetchAccounts()
  } catch (e) {
    showToast(e.response?.data?.error || 'Failed', 'error')
  } finally {
    accountActionLoading.value = false
  }
}

const adminUnsuspend = async (acc) => {
  accountActionLoading.value = true
  try {
    await axios.post(`/api/admin/cyberpanel/accounts/${acc.id}/unsuspend`)
    showToast(t('admin.cpUnsuspended'))
    await fetchAccounts()
  } catch (e) {
    showToast(e.response?.data?.error || 'Failed', 'error')
  } finally {
    accountActionLoading.value = false
  }
}

const confirmTerminateAccount = (acc) => {
  terminatingAccount.value = acc
  showTerminateConfirm.value = true
}

const terminateAccount = async () => {
  if (!terminatingAccount.value) return
  accountActionLoading.value = true
  try {
    await axios.delete(`/api/admin/cyberpanel/accounts/${terminatingAccount.value.id}`)
    showToast(t('admin.cpTerminated'))
    showTerminateConfirm.value = false
    await fetchAccounts()
  } catch (e) {
    showToast(e.response?.data?.error || 'Failed', 'error')
  } finally {
    accountActionLoading.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────────────────
const accountStatuses = ['pending', 'active', 'suspended', 'terminated']

const accountStatusLabel = (status) => ({
  pending: t('domains.cpHostingPending'),
  active: t('domains.cpHostingActive'),
  suspended: t('domains.cpHostingSuspended'),
  terminated: t('domains.cpHostingTerminated'),
}[status] || status)

const statusTone = (status) => ({
  pending: 'info',
  active: 'ok',
  suspended: 'warn',
  terminated: 'danger',
}[status] || 'muted')

const usagePercent = (server) => {
  if (!server.max_accounts) return 0
  return Math.min(100, Math.round((server.current_accounts / server.max_accounts) * 100))
}

const usageBarClass = (server) => {
  const pct = usagePercent(server)
  if (pct >= 100) return 'bg-y-danger'
  if (pct >= 85) return 'bg-y-warn'
  return 'bg-y-mint'
}
const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleDateString()
}

onMounted(async () => {
  await fetchServers()
})
</script>
