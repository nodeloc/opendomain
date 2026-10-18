<template>
  <div class="y-page">
    <PageHeader
      :eyebrow="$t('sponsorAdmin.eyebrow')"
      :title="$t('sponsorAdmin.title')"
      :description="$t('sponsorAdmin.description')"
    >
      <template #actions>
        <router-link to="/sponsor" target="_blank" class="btn btn-ghost gap-1.5">
          {{ $t('sponsorAdmin.viewPage') }}
          <ArrowUpRight :size="16" :stroke-width="1.5" />
        </router-link>
        <button type="button" class="btn btn-primary gap-1.5" @click="openCreate">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('sponsorAdmin.list.add') }}
        </button>
      </template>
    </PageHeader>

    <!-- Payment methods -->
    <section class="y-panel">
      <div class="y-panel-header">
        <div>
          <h2 class="m-0 text-base font-semibold text-y-text">{{ $t('sponsorAdmin.payment.title') }}</h2>
          <p class="m-0 mt-0.5 text-sm text-y-muted">{{ $t('sponsorAdmin.payment.description') }}</p>
        </div>
      </div>
      <YLoader v-if="settingsLoading" compact />
      <form v-else class="flex flex-col gap-5 p-5" novalidate @submit.prevent="saveSettings">
        <div class="grid gap-5 lg:grid-cols-2">
          <label class="y-field min-w-0">
            <span class="y-field-label">{{ $t('sponsorAdmin.payment.usdtLabel') }}</span>
            <input
              v-model.trim="settings.usdt_trc20_address"
              type="text"
              class="input input-bordered w-full font-mono text-[14px]"
              :class="{ 'input-error': usdtInvalid || settingsErrorField === 'usdt_trc20_address' }"
              placeholder="T…"
              autocomplete="off"
              spellcheck="false"
              maxlength="64"
              @input="clearSettingsError('usdt_trc20_address')"
            />
            <span v-if="usdtInvalid" class="text-xs text-y-danger-text">{{ $t('sponsorAdmin.payment.usdtInvalid') }}</span>
            <span v-else class="y-field-hint">{{ $t('sponsorAdmin.payment.usdtHint') }}</span>
          </label>
          <label class="y-field min-w-0">
            <span class="y-field-label">{{ $t('sponsorAdmin.payment.githubLabel') }}</span>
            <input
              v-model.trim="settings.github_sponsors_url"
              type="url"
              class="input input-bordered w-full font-mono text-[14px]"
              :class="{ 'input-error': githubInvalid || settingsErrorField === 'github_sponsors_url' }"
              placeholder="https://github.com/sponsors/…"
              autocomplete="off"
              spellcheck="false"
              @input="clearSettingsError('github_sponsors_url')"
            />
            <span v-if="githubInvalid" class="text-xs text-y-danger-text">{{ $t('sponsorAdmin.payment.githubInvalid') }}</span>
            <span v-else class="y-field-hint">{{ $t('sponsorAdmin.payment.githubHint') }}</span>
          </label>
        </div>
        <div v-if="settingsError" role="alert" class="flex items-start gap-2 text-sm text-y-danger-text">
          <CircleAlert :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" />
          <span>{{ settingsError }}</span>
        </div>
        <div>
          <button type="submit" class="btn btn-primary btn-sm gap-1.5" :disabled="settingsSaving || usdtInvalid || githubInvalid">
            <span v-if="settingsSaving" class="loading loading-spinner loading-xs"></span>
            <Save v-else :size="14" :stroke-width="1.5" />
            {{ $t('sponsorAdmin.payment.save') }}
          </button>
        </div>
      </form>
    </section>

    <!-- Sponsors -->
    <section class="y-panel">
      <div class="y-panel-header">
        <div class="flex items-baseline gap-3">
          <h2 class="m-0 text-base font-semibold text-y-text">{{ $t('sponsorAdmin.list.title') }}</h2>
          <span v-if="!loading" class="font-mono text-xs text-y-faint">{{ $t('sponsorAdmin.list.count', { n: sponsors.length }) }}</span>
        </div>
      </div>

      <YLoader v-if="loading" />

      <EmptyState
        v-else-if="sponsors.length === 0"
        :title="$t('sponsorAdmin.list.emptyTitle')"
        :description="$t('sponsorAdmin.list.emptyDesc')"
      >
        <button type="button" class="btn btn-outline btn-sm gap-1.5" @click="openCreate">
          <Plus :size="14" :stroke-width="1.5" />
          {{ $t('sponsorAdmin.list.add') }}
        </button>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="y-table">
          <thead>
            <tr>
              <th class="w-[112px]">{{ $t('sponsorAdmin.columns.logo') }}</th>
              <th>{{ $t('sponsorAdmin.columns.name') }}</th>
              <th>{{ $t('sponsorAdmin.columns.tier') }}</th>
              <th class="text-right">{{ $t('sponsorAdmin.columns.order') }}</th>
              <th>{{ $t('sponsorAdmin.columns.active') }}</th>
              <th class="text-right">{{ $t('sponsorAdmin.columns.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in sponsors" :key="s.id" :class="{ 'opacity-60': !s.is_active }">
              <td>
                <div class="y-inset flex h-10 w-20 items-center justify-center overflow-hidden p-1.5">
                  <SponsorLogo :src="logoFor(s)" :name="s.name" small />
                </div>
              </td>
              <td class="max-w-[320px]">
                <p class="m-0 truncate font-medium text-y-text">{{ s.name }}</p>
                <a
                  v-if="s.website_url"
                  :href="s.website_url"
                  target="_blank"
                  rel="noopener"
                  class="block truncate font-mono text-xs text-y-faint hover:text-y-mint-text"
                >{{ hostOf(s.website_url) }}</a>
              </td>
              <td class="whitespace-nowrap"><span :class="tierChip(s.tier)">{{ $t(`sponsor.tierSingle.${s.tier}`) }}</span></td>
              <td class="text-right font-mono">{{ s.display_order }}</td>
              <td>
                <input
                  type="checkbox"
                  class="toggle toggle-sm"
                  :checked="s.is_active"
                  :disabled="togglingId === s.id"
                  :aria-label="$t('sponsorAdmin.toggleActive')"
                  :title="$t('sponsorAdmin.toggleActive')"
                  @change="toggleActive(s, $event)"
                />
              </td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square"
                    :title="$t('sponsorAdmin.edit')"
                    :aria-label="$t('sponsorAdmin.edit')"
                    @click="openEdit(s)"
                  >
                    <Pencil :size="16" :stroke-width="1.5" />
                  </button>
                  <button
                    type="button"
                    class="btn btn-ghost btn-sm btn-square text-y-danger-text"
                    :title="$t('sponsorAdmin.delete')"
                    :aria-label="$t('sponsorAdmin.delete')"
                    @click="askDelete(s)"
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

    <!-- Create / edit modal -->
    <dialog class="modal" :class="{ 'modal-open': showModal }">
      <div class="modal-box relative w-11/12 max-w-3xl">
        <button
          type="button"
          class="btn btn-ghost btn-sm btn-square absolute right-3 top-3"
          :aria-label="$t('common.close')"
          @click="closeModal"
        >
          <X :size="16" :stroke-width="1.5" />
        </button>
        <h3 class="pr-10 text-lg font-semibold text-y-text">
          {{ editingId ? $t('sponsorAdmin.modal.editTitle') : $t('sponsorAdmin.modal.createTitle') }}
        </h3>
        <p class="mt-1 text-sm text-y-muted">{{ $t('sponsorAdmin.modal.description') }}</p>

        <form class="mt-6 flex flex-col gap-4" novalidate @submit.prevent="submit">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.name') }} <span class="text-y-danger-text" aria-hidden="true">*</span></span>
              <input
                v-model="form.name"
                type="text"
                maxlength="100"
                class="input input-bordered w-full"
                :class="{ 'input-error': fieldError === 'name' }"
                :placeholder="$t('sponsorAdmin.modal.namePlaceholder')"
                @input="clearFieldError('name')"
              />
            </label>
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.website') }}</span>
              <input
                v-model.trim="form.website_url"
                type="url"
                class="input input-bordered w-full font-mono text-[14px]"
                :class="{ 'input-error': fieldError === 'website_url' || !isHttpUrl(form.website_url) }"
                placeholder="https://"
                @input="clearFieldError('website_url')"
              />
            </label>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.tier') }}</span>
              <select v-model="form.tier" class="select select-bordered w-full" :class="{ 'select-error': fieldError === 'tier' }" @change="clearFieldError('tier')">
                <option v-for="tier in TIERS" :key="tier" :value="tier">{{ $t(`sponsor.tierSingle.${tier}`) }}</option>
              </select>
            </label>
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.order') }}</span>
              <input
                v-model.number="form.display_order"
                type="number"
                step="1"
                class="input input-bordered w-full font-mono"
                :class="{ 'input-error': fieldError === 'display_order' }"
                @input="clearFieldError('display_order')"
              />
              <span class="y-field-hint">{{ $t('sponsorAdmin.modal.orderHint') }}</span>
            </label>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.logo') }}</span>
              <input
                v-model.trim="form.logo_url"
                type="url"
                class="input input-bordered w-full font-mono text-[14px]"
                :class="{ 'input-error': fieldError === 'logo_url' || !isHttpUrl(form.logo_url) }"
                placeholder="https://"
                @input="clearFieldError('logo_url')"
              />
              <span class="y-field-hint">{{ $t('sponsorAdmin.modal.logoHint') }}</span>
            </label>
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.logoDark') }}</span>
              <input
                v-model.trim="form.logo_dark_url"
                type="url"
                class="input input-bordered w-full font-mono text-[14px]"
                :class="{ 'input-error': fieldError === 'logo_dark_url' || !isHttpUrl(form.logo_dark_url) }"
                placeholder="https://"
                @input="clearFieldError('logo_dark_url')"
              />
              <span class="y-field-hint">{{ $t('sponsorAdmin.modal.logoDarkHint') }}</span>
            </label>
          </div>

          <!-- Live preview in both themes -->
          <div class="flex flex-col gap-2">
            <span class="y-label">{{ $t('sponsorAdmin.modal.preview') }}</span>
            <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div
                v-for="p in previews"
                :key="p.theme"
                :data-theme="p.theme"
                class="flex flex-col gap-2 rounded-md border border-y-line bg-y-bg p-3"
              >
                <span class="font-mono text-[11px] uppercase tracking-[0.08em] text-y-faint">{{ p.label }}</span>
                <div class="flex h-20 items-center rounded-md border border-y-line bg-y-surface px-4">
                  <div class="flex h-12 min-w-0 items-center">
                    <SponsorLogo :src="p.src" :name="form.name || '?'" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.descEn') }}</span>
              <textarea
                v-model="form.description_en"
                rows="3"
                maxlength="500"
                class="textarea textarea-bordered w-full text-sm leading-6"
                :class="{ 'textarea-error': fieldError === 'description_en' }"
                @input="clearFieldError('description_en')"
              ></textarea>
            </label>
            <label class="y-field">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.descZh') }}</span>
              <textarea
                v-model="form.description_zh"
                rows="3"
                maxlength="500"
                class="textarea textarea-bordered w-full text-sm leading-6"
                :class="{ 'textarea-error': fieldError === 'description_zh' }"
                @input="clearFieldError('description_zh')"
              ></textarea>
            </label>
          </div>

          <label class="flex cursor-pointer items-start gap-3 rounded-md border border-y-line px-4 py-3">
            <input v-model="form.is_active" type="checkbox" class="checkbox checkbox-sm mt-0.5" />
            <span class="flex flex-col gap-0.5">
              <span class="y-field-label">{{ $t('sponsorAdmin.modal.active') }}</span>
              <span class="y-field-hint">{{ $t('sponsorAdmin.modal.activeHint') }}</span>
            </span>
          </label>

          <div v-if="formError" role="alert" class="flex items-start gap-2 text-sm text-y-danger-text">
            <CircleAlert :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0" />
            <span>{{ formError }}</span>
          </div>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModal">{{ $t('sponsorAdmin.modal.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              {{ editingId ? $t('sponsorAdmin.modal.save') : $t('sponsorAdmin.modal.create') }}
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>{{ $t('common.close') }}</button>
      </form>
    </dialog>

    <!-- Delete confirmation -->
    <dialog class="modal" :class="{ 'modal-open': !!deleteTarget }">
      <div class="modal-box">
        <h3 class="text-lg font-semibold text-y-text">{{ $t('sponsorAdmin.deleteTitle') }}</h3>
        <p class="mt-3 text-sm leading-6 text-y-muted">{{ $t('sponsorAdmin.deleteMessage', { name: deleteTarget?.name }) }}</p>
        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="deleteTarget = null">{{ $t('common.cancel') }}</button>
          <button type="button" class="btn btn-error" :disabled="deleting" @click="handleDelete">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            <span v-else>{{ $t('common.delete') }}</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="deleteTarget = null">
        <button>{{ $t('common.close') }}</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, Pencil, Trash2, X, Save, CircleAlert, ArrowUpRight } from 'lucide-vue-next'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useTheme } from '../composables/useTheme'
import PageHeader from '../components/ui/PageHeader.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'
import SponsorLogo from '../components/sponsor/SponsorLogo.vue'

const { t } = useI18n()
const toast = useToast()
const { isDark } = useTheme()

const TIERS = ['partner', 'sponsor', 'backer']
// Mirrors the server-side checks so admins see problems before saving.
const TRC20_RE = /^T[1-9A-HJ-NP-Za-km-z]{33}$/

const emptyForm = () => ({
  name: '',
  website_url: '',
  logo_url: '',
  logo_dark_url: '',
  description_zh: '',
  description_en: '',
  tier: 'sponsor',
  display_order: 0,
  is_active: true,
})

const sponsors = ref([])
const loading = ref(true)
const settings = ref({ usdt_trc20_address: '', github_sponsors_url: '' })
const settingsLoading = ref(true)
const settingsSaving = ref(false)
const settingsError = ref('')
const settingsErrorField = ref('')

const showModal = ref(false)
const editingId = ref(null)
const form = ref(emptyForm())
const submitting = ref(false)
const formError = ref('')
const fieldError = ref('')
const togglingId = ref(null)
const deleteTarget = ref(null)
const deleting = ref(false)

const isHttpUrl = (value) => {
  if (!value) return true
  try {
    const u = new URL(value)
    return u.protocol === 'http:' || u.protocol === 'https:'
  } catch {
    return false
  }
}

const isGithubUrl = (value) => {
  try {
    const u = new URL(value)
    return u.protocol === 'https:' && u.hostname === 'github.com'
  } catch {
    return false
  }
}

const usdtInvalid = computed(() => !!settings.value.usdt_trc20_address && !TRC20_RE.test(settings.value.usdt_trc20_address))
const githubInvalid = computed(() => !!settings.value.github_sponsors_url && !isGithubUrl(settings.value.github_sponsors_url))

const previews = computed(() => [
  { theme: 'yunamo-light', label: t('sponsorAdmin.modal.previewLight'), src: form.value.logo_url },
  { theme: 'yunamo', label: t('sponsorAdmin.modal.previewDark'), src: form.value.logo_dark_url || form.value.logo_url },
])

const logoFor = (s) => (isDark.value && s.logo_dark_url ? s.logo_dark_url : s.logo_url) || ''

const hostOf = (url) => {
  try {
    return new URL(url).host
  } catch {
    return url
  }
}

const tierChip = (tier) => ({ partner: 'y-chip-mint', sponsor: 'y-chip-info' }[tier] || 'y-chip')

const toPayload = (s) => ({
  name: (s.name || '').trim(),
  website_url: (s.website_url || '').trim(),
  logo_url: (s.logo_url || '').trim(),
  logo_dark_url: (s.logo_dark_url || '').trim(),
  description_zh: (s.description_zh || '').trim(),
  description_en: (s.description_en || '').trim(),
  tier: s.tier,
  display_order: Number.isFinite(Number(s.display_order)) ? Math.trunc(Number(s.display_order)) : 0,
  is_active: !!s.is_active,
})

const fetchSponsors = async () => {
  loading.value = true
  try {
    const { data } = await axios.get('/api/admin/sponsors')
    sponsors.value = Array.isArray(data?.sponsors) ? data.sponsors : []
  } catch (err) {
    toast.error(err.response?.data?.error || t('sponsorAdmin.list.loadFailed'))
  } finally {
    loading.value = false
  }
}

const fetchSettings = async () => {
  settingsLoading.value = true
  try {
    const { data } = await axios.get('/api/admin/sponsor-settings')
    settings.value = {
      usdt_trc20_address: data?.usdt_trc20_address || '',
      github_sponsors_url: data?.github_sponsors_url || '',
    }
  } catch (err) {
    toast.error(err.response?.data?.error || t('sponsorAdmin.payment.loadFailed'))
  } finally {
    settingsLoading.value = false
  }
}

const clearSettingsError = (field) => {
  if (settingsErrorField.value === field) {
    settingsErrorField.value = ''
    settingsError.value = ''
  }
}

const saveSettings = async () => {
  if (usdtInvalid.value || githubInvalid.value) return
  settingsSaving.value = true
  settingsError.value = ''
  settingsErrorField.value = ''
  try {
    const { data } = await axios.put('/api/admin/sponsor-settings', {
      usdt_trc20_address: settings.value.usdt_trc20_address,
      github_sponsors_url: settings.value.github_sponsors_url,
    })
    settings.value = {
      usdt_trc20_address: data?.usdt_trc20_address ?? settings.value.usdt_trc20_address,
      github_sponsors_url: data?.github_sponsors_url ?? settings.value.github_sponsors_url,
    }
    toast.success(t('sponsorAdmin.payment.saved'))
  } catch (err) {
    const msg = err.response?.data?.error || t('sponsorAdmin.payment.saveFailed')
    settingsError.value = msg
    settingsErrorField.value = err.response?.data?.field || ''
    toast.error(msg)
  } finally {
    settingsSaving.value = false
  }
}

const openCreate = () => {
  editingId.value = null
  const maxOrder = sponsors.value.reduce((m, s) => Math.max(m, Number(s.display_order) || 0), 0)
  form.value = { ...emptyForm(), display_order: sponsors.value.length ? maxOrder + 10 : 0 }
  formError.value = ''
  fieldError.value = ''
  showModal.value = true
}

const openEdit = (s) => {
  editingId.value = s.id
  form.value = { ...emptyForm(), ...toPayload(s) }
  formError.value = ''
  fieldError.value = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingId.value = null
  formError.value = ''
  fieldError.value = ''
}

const clearFieldError = (field) => {
  if (fieldError.value === field) {
    fieldError.value = ''
    formError.value = ''
  }
}

const validateForm = () => {
  if (!form.value.name.trim()) return { field: 'name', message: t('sponsorAdmin.modal.nameRequired') }
  for (const field of ['website_url', 'logo_url', 'logo_dark_url']) {
    if (!isHttpUrl(form.value[field])) return { field, message: t('sponsorAdmin.modal.urlInvalid') }
  }
  return null
}

const submit = async () => {
  const problem = validateForm()
  if (problem) {
    fieldError.value = problem.field
    formError.value = problem.message
    return
  }
  submitting.value = true
  formError.value = ''
  fieldError.value = ''
  try {
    const payload = toPayload(form.value)
    if (editingId.value) {
      await axios.put(`/api/admin/sponsors/${editingId.value}`, payload)
      toast.success(t('sponsorAdmin.toasts.updated'))
    } else {
      await axios.post('/api/admin/sponsors', payload)
      toast.success(t('sponsorAdmin.toasts.created'))
    }
    closeModal()
    await fetchSponsors()
  } catch (err) {
    formError.value = err.response?.data?.error || t('sponsorAdmin.toasts.saveFailed')
    fieldError.value = err.response?.data?.field || ''
  } finally {
    submitting.value = false
  }
}

const toggleActive = async (s, event) => {
  const next = event.target.checked
  togglingId.value = s.id
  try {
    const { data } = await axios.put(`/api/admin/sponsors/${s.id}`, toPayload({ ...s, is_active: next }))
    Object.assign(s, data?.sponsor || { is_active: next })
    toast.success(t(next ? 'sponsorAdmin.toasts.activated' : 'sponsorAdmin.toasts.deactivated', { name: s.name }))
  } catch (err) {
    event.target.checked = s.is_active
    toast.error(err.response?.data?.error || t('sponsorAdmin.toasts.saveFailed'))
  } finally {
    togglingId.value = null
  }
}

const askDelete = (s) => {
  deleteTarget.value = s
}

const handleDelete = async () => {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await axios.delete(`/api/admin/sponsors/${deleteTarget.value.id}`)
    toast.success(t('sponsorAdmin.toasts.deleted'))
    deleteTarget.value = null
    await fetchSponsors()
  } catch (err) {
    toast.error(err.response?.data?.error || t('sponsorAdmin.toasts.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchSponsors()
  fetchSettings()
})
</script>

<style scoped>
/* The global focus style paints inputs mint, which would hide a field the server rejected. */
.input-error,
.input-error:focus,
.select-error,
.select-error:focus,
.textarea-error,
.textarea-error:focus {
  border-color: rgb(var(--y-danger) / 0.8);
  outline-color: rgb(var(--y-danger) / 0.4);
}
</style>
