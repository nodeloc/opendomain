<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('shell.adminArea')"
      :title="$t('partnerAdmin.title')"
      :description="$t('partnerAdmin.description')"
    >
      <template #actions>
        <router-link to="/sponsor#become-partner" class="btn btn-ghost gap-2">
          <ArrowUpRight :size="16" :stroke-width="1.5" />
          {{ $t('partnerAdmin.viewPage') }}
        </router-link>
        <button type="button" class="btn btn-outline gap-2" :disabled="loading" @click="fetchApplications">
          <span v-if="loading" class="loading loading-spinner loading-sm"></span>
          <RefreshCw v-else :size="16" :stroke-width="1.5" />
          {{ $t('partnerAdmin.refresh') }}
        </button>
      </template>
    </PageHeader>

    <!-- Status filter tabs with counts -->
    <div class="y-seg w-fit max-w-full overflow-x-auto" role="group" :aria-label="$t('partnerAdmin.fields.status')">
      <button
        v-for="tab in statusTabs"
        :key="tab.key"
        type="button"
        class="whitespace-nowrap"
        :aria-pressed="statusFilter === tab.key"
        @click="filterByStatus(tab.key)"
      >
        {{ tab.label }}
        <span class="ml-1.5 font-mono text-xs opacity-80">{{ tab.count }}</span>
      </button>
    </div>

    <section class="y-panel overflow-hidden">
      <div class="y-panel-header">
        <label class="relative block w-full sm:max-w-md">
          <span class="sr-only">{{ $t('common.search') }}</span>
          <Search :size="16" :stroke-width="1.5" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-y-faint" />
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('partnerAdmin.searchPlaceholder')"
            class="input input-bordered h-10 min-h-0 w-full pl-9 pr-10 text-sm"
            @input="debouncedSearch"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="btn btn-ghost btn-xs btn-square absolute right-2 top-1/2 -translate-y-1/2 text-y-faint"
            :aria-label="$t('uiAdmin1.clearSearch')"
            @click="clearSearch"
          >
            <X :size="14" :stroke-width="1.5" />
          </button>
        </label>
        <span class="font-mono text-xs text-y-faint">{{ $t('partnerAdmin.totalCount', { count: total }) }}</span>
      </div>

      <YLoader v-if="loading && applications.length === 0" :label="$t('common.loading')" />

      <EmptyState
        v-else-if="applications.length === 0"
        :title="searchQuery || statusFilter ? $t('partnerAdmin.noResults') : $t('partnerAdmin.noApplications')"
        :description="searchQuery || statusFilter ? '' : $t('partnerAdmin.noApplicationsDesc')"
      />

      <ul v-else class="m-0 list-none divide-y divide-y-line p-0" :class="{ 'opacity-60': loading }">
        <li v-for="app in applications" :key="app.id" class="flex flex-col gap-4 px-5 py-5">
          <!-- Header: who, what, when -->
          <div class="flex flex-wrap items-start justify-between gap-x-6 gap-y-3">
            <div class="flex min-w-0 flex-col gap-1">
              <div class="flex flex-wrap items-center gap-x-2.5 gap-y-1.5">
                <h3 class="m-0 break-words text-base font-semibold text-y-text">{{ app.name }}</h3>
                <span v-if="app.organization" class="break-words text-sm text-y-muted">· {{ app.organization }}</span>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <span :class="typeChipClass(app.partnership_type)">{{ typeLabel(app.partnership_type) }}</span>
                <span :class="statusChipClass(app.status)">{{ statusLabel(app.status) }}</span>
                <span class="font-mono text-xs text-y-faint">#{{ app.id }}</span>
              </div>
            </div>
            <time class="whitespace-nowrap font-mono text-xs text-y-faint" :datetime="app.created_at" :title="$t('partnerAdmin.fields.submitted')">
              {{ formatDateTime(app.created_at) }}
            </time>
          </div>

          <!-- Contact details -->
          <dl class="m-0 grid grid-cols-1 gap-x-6 gap-y-2 text-sm sm:grid-cols-2 xl:grid-cols-4">
            <div class="flex min-w-0 items-start gap-2">
              <dt class="shrink-0 text-y-faint"><Mail :size="15" :stroke-width="1.5" class="mt-0.5" /><span class="sr-only">{{ $t('partnerAdmin.fields.email') }}</span></dt>
              <dd class="m-0 min-w-0 break-all">
                <a :href="`mailto:${app.email}`" class="text-y-mint-text hover:text-y-mint-hover">{{ app.email }}</a>
              </dd>
            </div>
            <div v-if="app.contact" class="flex min-w-0 items-start gap-2">
              <dt class="shrink-0 text-y-faint"><MessageCircle :size="15" :stroke-width="1.5" class="mt-0.5" /><span class="sr-only">{{ $t('partnerAdmin.fields.contact') }}</span></dt>
              <dd class="m-0 min-w-0 break-all text-y-text">{{ app.contact }}</dd>
            </div>
            <div v-if="app.website" class="flex min-w-0 items-start gap-2">
              <dt class="shrink-0 text-y-faint"><Globe :size="15" :stroke-width="1.5" class="mt-0.5" /><span class="sr-only">{{ $t('partnerAdmin.fields.website') }}</span></dt>
              <dd class="m-0 min-w-0 break-all">
                <a
                  v-if="isHttpUrl(app.website)"
                  :href="app.website"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center gap-1 text-y-mint-text hover:text-y-mint-hover"
                >
                  <span class="break-all">{{ displayUrl(app.website) }}</span>
                  <ArrowUpRight :size="13" :stroke-width="1.75" class="shrink-0" />
                </a>
                <span v-else class="text-y-muted">{{ app.website }}</span>
              </dd>
            </div>
            <div v-if="app.ip || app.locale" class="flex min-w-0 items-start gap-2">
              <dt class="shrink-0 text-y-faint"><Info :size="15" :stroke-width="1.5" class="mt-0.5" /><span class="sr-only">{{ $t('partnerAdmin.fields.ip') }}</span></dt>
              <dd class="m-0 min-w-0 break-all font-mono text-xs leading-5 text-y-muted">
                {{ [app.ip, app.locale].filter(Boolean).join(' · ') }}
              </dd>
            </div>
          </dl>

          <!-- Message -->
          <div>
            <p class="y-label m-0">{{ $t('partnerAdmin.fields.message') }}</p>
            <div
              :id="`partner-msg-${app.id}`"
              class="y-inset mt-1.5 whitespace-pre-wrap break-words px-3.5 py-2.5 text-sm leading-6 text-y-text"
              :class="{ 'partner-msg-collapsed': isLong(app.message) && !expanded[app.id] }"
            >{{ app.message }}</div>
            <button
              v-if="isLong(app.message)"
              type="button"
              class="btn btn-ghost btn-xs mt-1.5 gap-1 text-y-muted"
              :aria-expanded="expanded[app.id] ? 'true' : 'false'"
              :aria-controls="`partner-msg-${app.id}`"
              @click="expanded[app.id] = !expanded[app.id]"
            >
              <ChevronDown :size="14" :stroke-width="1.5" class="transition-transform" :class="{ 'rotate-180': expanded[app.id] }" />
              {{ expanded[app.id] ? $t('partnerAdmin.showLess') : $t('partnerAdmin.showMore') }}
            </button>
          </div>

          <!-- Admin note -->
          <div class="y-field">
            <label class="y-field-label" :for="`partner-note-${app.id}`">{{ $t('partnerAdmin.fields.adminNote') }}</label>
            <textarea
              :id="`partner-note-${app.id}`"
              v-model="notes[app.id]"
              rows="2"
              maxlength="5000"
              class="textarea textarea-bordered w-full text-sm"
              :placeholder="$t('partnerAdmin.adminNotePlaceholder')"
            ></textarea>
          </div>

          <!-- Actions -->
          <div class="flex flex-wrap items-center gap-2">
            <label class="flex items-center gap-2">
              <span class="sr-only">{{ $t('partnerAdmin.fields.status') }}</span>
              <select
                class="select select-bordered select-sm"
                :value="app.status"
                :disabled="busy[app.id]"
                :aria-label="$t('partnerAdmin.fields.status')"
                @change="changeStatus(app, $event.target.value)"
              >
                <option v-for="key in statusKeys" :key="key" :value="key">{{ statusLabel(key) }}</option>
              </select>
            </label>
            <button
              type="button"
              class="btn btn-outline btn-sm"
              :disabled="busy[app.id] || !noteDirty(app)"
              @click="saveNote(app)"
            >
              <span v-if="busy[app.id] === 'note'" class="loading loading-spinner loading-xs"></span>
              <Save v-else :size="14" :stroke-width="1.5" />
              {{ $t('partnerAdmin.saveNote') }}
            </button>
            <a :href="replyMailto(app)" class="btn btn-ghost btn-sm gap-1.5">
              <Reply :size="14" :stroke-width="1.5" />
              {{ $t('partnerAdmin.reply') }}
            </a>
            <button
              type="button"
              class="btn btn-ghost btn-sm ml-auto gap-1.5 text-y-danger-text"
              :disabled="!!busy[app.id]"
              @click="confirmDelete(app)"
            >
              <Trash2 :size="14" :stroke-width="1.5" />
              {{ $t('partnerAdmin.delete') }}
            </button>
          </div>
        </li>
      </ul>

      <div
        v-if="!loading && totalPages > 1"
        class="flex flex-col items-center justify-between gap-3 border-t border-y-line px-5 py-3 sm:flex-row"
      >
        <p class="font-mono text-xs text-y-faint">{{ $t('uiAdmin1.pageOf', { page, total: totalPages }) }}</p>
        <div class="join">
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="page === 1"
            :aria-label="$t('common.previous')"
            @click="goToPage(page - 1)"
          >
            <ChevronLeft :size="16" :stroke-width="1.5" />
          </button>
          <span class="join-item flex h-[2.125rem] items-center px-3 font-mono text-xs text-y-muted">{{ page }}</span>
          <button
            type="button"
            class="join-item btn btn-ghost btn-sm btn-square"
            :disabled="page >= totalPages"
            :aria-label="$t('common.next')"
            @click="goToPage(page + 1)"
          >
            <ChevronRight :size="16" :stroke-width="1.5" />
          </button>
        </div>
      </div>
    </section>

    <!-- Delete confirmation modal -->
    <dialog class="modal" :class="{ 'modal-open': showDelete }">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('partnerAdmin.confirmDeleteTitle') }}</h3>
        <p class="mt-3 break-words text-sm leading-6 text-y-muted">
          {{ $t('partnerAdmin.confirmDeleteMessage', { name: deleteTarget ? deleteTarget.organization || deleteTarget.name : '' }) }}
        </p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="showDelete = false">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-error" :disabled="deleting" @click="handleDelete">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            <span v-else>{{ $t('common.delete') }}</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="showDelete = false">
        <button>{{ $t('common.close') }}</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { formatDateTime } from '../utils/time'
import {
  RefreshCw, Search, X, Trash2, ChevronLeft, ChevronRight, ChevronDown, ArrowUpRight,
  Mail, MessageCircle, Globe, Info, Save, Reply,
} from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

const { t } = useI18n()
const toast = useToast()

const statusKeys = ['new', 'contacted', 'accepted', 'declined']
const typeKeys = ['sponsorship', 'resources', 'business', 'community', 'other']
const LONG_MESSAGE_CHARS = 360
const LONG_MESSAGE_LINES = 6

const loading = ref(false)
const applications = ref([])
const counts = ref({ new: 0, contacted: 0, accepted: 0, declined: 0 })
const total = ref(0)
const page = ref(1)
const pageSize = 20
const statusFilter = ref('')
const searchQuery = ref('')

const notes = reactive({})
const savedNotes = {}
const expanded = reactive({})
const busy = reactive({})

const showDelete = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

const statusTabs = computed(() => {
  const allCount = statusKeys.reduce((sum, k) => sum + (counts.value[k] || 0), 0)
  return [
    { key: '', label: t('partnerAdmin.tabs.all'), count: allCount },
    ...statusKeys.map((k) => ({ key: k, label: t(`partnerAdmin.tabs.${k}`), count: counts.value[k] || 0 })),
  ]
})

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

const syncNote = (app) => {
  notes[app.id] = app.admin_note || ''
  savedNotes[app.id] = notes[app.id]
}

const fetchApplications = async () => {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize }
    if (statusFilter.value) params.status = statusFilter.value
    if (searchQuery.value.trim()) params.search = searchQuery.value.trim()
    const { data } = await axios.get('/api/admin/partner-applications', { params })
    applications.value = data.applications || []
    total.value = data.total || 0
    counts.value = { new: 0, contacted: 0, accepted: 0, declined: 0, ...(data.counts || {}) }
    // Keep a note the admin is still typing; refresh untouched ones from the server.
    for (const app of applications.value) {
      if (notes[app.id] === undefined || notes[app.id] === savedNotes[app.id]) syncNote(app)
    }
  } catch (error) {
    toast.error(error.response?.data?.error || t('partnerAdmin.toasts.fetchFailed'))
  } finally {
    loading.value = false
  }
}

let searchTimeout = null
const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchApplications()
  }, 400)
}

const clearSearch = () => {
  searchQuery.value = ''
  page.value = 1
  fetchApplications()
}

const filterByStatus = (key) => {
  statusFilter.value = key
  page.value = 1
  fetchApplications()
}

const goToPage = (p) => {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  fetchApplications()
}

const replaceApp = (updated) => {
  const idx = applications.value.findIndex((a) => a.id === updated.id)
  if (idx !== -1) applications.value[idx] = updated
}

const update = async (app, body, kind) => {
  busy[app.id] = kind
  try {
    const { data } = await axios.put(`/api/admin/partner-applications/${app.id}`, body)
    const updated = data.application
    replaceApp(updated)
    return updated
  } catch (error) {
    toast.error(error.response?.data?.error || t('partnerAdmin.toasts.saveFailed'))
    return null
  } finally {
    delete busy[app.id]
  }
}

const changeStatus = async (app, status) => {
  if (status === app.status) return
  const previous = app.status
  const updated = await update(app, { status }, 'status')
  if (!updated) {
    // Put the select back to the saved value.
    replaceApp({ ...app, status: previous })
    return
  }
  toast.success(t('partnerAdmin.toasts.statusUpdated', { status: statusLabel(updated.status) }))
  // Counts change and the row may no longer match the active tab.
  fetchApplications()
}

const noteDirty = (app) => (notes[app.id] ?? '') !== (app.admin_note || '')

const saveNote = async (app) => {
  const updated = await update(app, { admin_note: notes[app.id] ?? '' }, 'note')
  if (!updated) return
  syncNote(updated)
  toast.success(t('partnerAdmin.noteSaved'))
}

const confirmDelete = (app) => {
  deleteTarget.value = app
  showDelete.value = true
}

const handleDelete = async () => {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await axios.delete(`/api/admin/partner-applications/${deleteTarget.value.id}`)
    toast.success(t('partnerAdmin.toasts.deleteSuccess'))
    showDelete.value = false
    delete notes[deleteTarget.value.id]
    delete savedNotes[deleteTarget.value.id]
    deleteTarget.value = null
    if (applications.value.length === 1 && page.value > 1) page.value -= 1
    fetchApplications()
  } catch (error) {
    toast.error(error.response?.data?.error || t('partnerAdmin.toasts.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

const isHttpUrl = (value) => {
  try {
    const u = new URL(value)
    return u.protocol === 'http:' || u.protocol === 'https:'
  } catch {
    return false
  }
}

const displayUrl = (value) => value.replace(/^https?:\/\//i, '').replace(/\/$/, '')

const isLong = (text) => !!text && (text.length > LONG_MESSAGE_CHARS || text.split('\n').length > LONG_MESSAGE_LINES)

const replyMailto = (app) => `mailto:${app.email}?subject=${encodeURIComponent(t('partnerAdmin.replySubject'))}`

const statusLabel = (s) => (statusKeys.includes(s) ? t(`partnerAdmin.status.${s}`) : s)
const typeLabel = (type) => (typeKeys.includes(type) ? t(`partnerAdmin.types.${type}`) : type)

const statusChipClass = (s) => ({
  new: 'y-chip-info',
  contacted: 'y-chip-warn',
  accepted: 'y-chip-mint',
  declined: 'y-chip',
}[s] || 'y-chip')

const typeChipClass = () => 'y-tag'

onMounted(fetchApplications)
</script>

<style scoped>
.partner-msg-collapsed {
  max-height: 9rem;
  overflow: hidden;
  -webkit-mask-image: linear-gradient(to bottom, #000 60%, transparent);
  mask-image: linear-gradient(to bottom, #000 60%, transparent);
}
</style>
