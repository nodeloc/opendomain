<template>
  <div class="y-page">
    <div class="space-y-4">
      <router-link to="/domains" class="inline-flex items-center gap-1.5 text-sm text-y-muted transition-colors hover:text-y-text">
        <ArrowLeft :size="16" :stroke-width="1.5" />
        {{ $t('uiUser.dns.back') }}
      </router-link>
      <PageHeader :eyebrow="$t('uiUser.dns.eyebrow')" :description="$t('uiUser.dns.description')">
        <template #title>
          <span class="y-domain font-semibold tracking-[-0.02em]">{{ domain?.full_domain || '—' }}</span>
        </template>
        <template #actions>
          <button type="button" class="btn btn-outline gap-1.5" :disabled="syncing" @click="syncFromPowerDNS">
            <span v-if="syncing" class="loading loading-spinner loading-sm"></span>
            <RefreshCw v-else :size="16" :stroke-width="1.5" />
            {{ $t('uiUser.dns.sync') }}
          </button>
          <button type="button" class="btn btn-primary gap-1.5" @click="showAddModal = true">
            <Plus :size="16" :stroke-width="1.5" />
            {{ $t('uiUser.dns.addRecord') }}
          </button>
        </template>
      </PageHeader>
    </div>

    <!-- Records -->
    <section class="y-panel">
      <div class="y-panel-header">
        <h2 class="y-h3">{{ $t('uiUser.dns.records') }}</h2>
        <span v-if="!loading && records.length" class="y-chip font-mono">{{ $t('uiUser.dns.recordCount', records.length) }}</span>
      </div>

      <YLoader v-if="loading" :label="$t('uiUser.common.loading')" />

      <EmptyState v-else-if="records.length === 0" :title="$t('uiUser.dns.emptyTitle')" :description="$t('uiUser.dns.emptyDesc')">
        <button type="button" class="btn btn-outline gap-1.5" @click="showAddModal = true">
          <Plus :size="16" :stroke-width="1.5" />
          {{ $t('uiUser.dns.addRecord') }}
        </button>
      </EmptyState>

      <template v-else>
        <!-- Desktop table -->
        <div class="hidden overflow-x-auto md:block">
          <table class="y-table">
            <thead>
              <tr>
                <th>{{ $t('dnsManagement.name') }}</th>
                <th>{{ $t('dnsManagement.type') }}</th>
                <th>{{ $t('dnsManagement.content') }}</th>
                <th>{{ $t('dnsManagement.ttl') }}</th>
                <th>{{ $t('dnsManagement.priority') }}</th>
                <th>{{ $t('dnsManagement.status') }}</th>
                <th class="text-right"><span class="sr-only">{{ $t('dnsManagement.actions') }}</span></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in records" :key="record.id">
                <td class="font-mono text-[13px]">{{ record.name }}</td>
                <td><span class="y-tag">{{ record.type }}</span></td>
                <td class="max-w-[22rem] font-mono text-[13px] text-y-muted">
                  <span class="block truncate" :title="record.content">{{ truncate(record.content, 40) }}</span>
                </td>
                <td class="whitespace-nowrap font-mono text-[13px] text-y-muted">{{ record.ttl }}s</td>
                <td class="font-mono text-[13px] text-y-muted">{{ record.priority || '—' }}</td>
                <td>
                  <div class="flex flex-wrap items-center gap-2 whitespace-nowrap">
                    <StatusDot
                      :tone="record.synced_to_powerdns ? 'ok' : 'warn'"
                      :label="record.synced_to_powerdns ? $t('dnsManagement.synced') : $t('dnsManagement.pending')"
                    />
                    <span v-if="!record.is_active" class="y-chip">{{ $t('dnsManagement.inactive') }}</span>
                    <button
                      v-if="record.sync_error"
                      type="button"
                      class="btn btn-ghost btn-xs btn-square text-y-danger-text"
                      :title="$t('uiUser.dns.viewError')"
                      :aria-label="$t('uiUser.dns.viewError')"
                      @click="showErrorDetails(record.sync_error)"
                    >
                      <AlertTriangle :size="16" :stroke-width="1.5" />
                    </button>
                  </div>
                </td>
                <td>
                  <div class="flex justify-end gap-1">
                    <button type="button" class="btn btn-ghost btn-sm btn-square" :title="$t('uiUser.dns.edit')" :aria-label="$t('uiUser.dns.edit')" @click="editRecord(record)">
                      <Pencil :size="16" :stroke-width="1.5" />
                    </button>
                    <button type="button" class="btn btn-ghost btn-sm btn-square text-y-danger-text" :title="$t('uiUser.dns.delete')" :aria-label="$t('uiUser.dns.delete')" @click="confirmDelete(record)">
                      <Trash2 :size="16" :stroke-width="1.5" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Mobile list -->
        <ul class="divide-y divide-y-line md:hidden">
          <li v-for="record in records" :key="record.id" class="space-y-2 px-5 py-4">
            <div class="flex items-start justify-between gap-3">
              <div class="flex min-w-0 items-center gap-2">
                <span class="y-tag">{{ record.type }}</span>
                <span class="truncate font-mono text-[13px] text-y-text">{{ record.name }}</span>
              </div>
              <div class="-mr-2 -mt-1 flex shrink-0 gap-1">
                <button type="button" class="btn btn-ghost btn-sm btn-square" :aria-label="$t('uiUser.dns.edit')" @click="editRecord(record)">
                  <Pencil :size="16" :stroke-width="1.5" />
                </button>
                <button type="button" class="btn btn-ghost btn-sm btn-square text-y-danger-text" :aria-label="$t('uiUser.dns.delete')" @click="confirmDelete(record)">
                  <Trash2 :size="16" :stroke-width="1.5" />
                </button>
              </div>
            </div>
            <p class="break-all font-mono text-[13px] text-y-muted">{{ record.content }}</p>
            <div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-xs text-y-faint">
              <span class="font-mono">TTL {{ record.ttl }}s</span>
              <span v-if="record.priority" class="font-mono">{{ $t('dnsManagement.priority') }} {{ record.priority }}</span>
              <StatusDot
                :tone="record.synced_to_powerdns ? 'ok' : 'warn'"
                :label="record.synced_to_powerdns ? $t('dnsManagement.synced') : $t('dnsManagement.pending')"
              />
              <span v-if="!record.is_active" class="y-chip">{{ $t('dnsManagement.inactive') }}</span>
              <button v-if="record.sync_error" type="button" class="inline-flex items-center gap-1 text-y-danger-text" @click="showErrorDetails(record.sync_error)">
                <AlertTriangle :size="14" :stroke-width="1.5" />
                {{ $t('uiUser.dns.viewError') }}
              </button>
            </div>
          </li>
        </ul>
      </template>
    </section>

    <!-- DNSSEC -->
    <section class="y-panel">
      <div class="y-panel-header">
        <div class="flex min-w-0 items-center gap-3">
          <span class="y-icon-tile" :class="{ '!text-y-faint': !dnssec.enabled }">
            <ShieldCheck v-if="dnssec.enabled" :size="20" :stroke-width="1.5" />
            <Shield v-else :size="20" :stroke-width="1.5" />
          </span>
          <div class="min-w-0">
            <h2 class="y-h3">DNSSEC</h2>
            <p class="mt-0.5 text-sm text-y-muted">{{ $t('uiUser.dns.dnssecDesc') }}</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <StatusDot
            :tone="dnssec.enabled ? 'ok' : 'muted'"
            :label="dnssec.enabled ? $t('dnsManagement.dnssecEnabled') : $t('dnsManagement.dnssecDisabled')"
          />
          <button
            v-if="usesDefaultNS"
            type="button"
            class="btn btn-sm"
            :class="dnssec.enabled ? 'btn-ghost text-y-danger-text' : 'btn-outline'"
            :disabled="dnssec.loading"
            @click="toggleDNSSEC"
          >
            <span v-if="dnssec.loading" class="loading loading-spinner loading-xs"></span>
            <span v-else>{{ dnssec.enabled ? $t('uiUser.dns.dnssecDisable') : $t('uiUser.dns.dnssecEnable') }}</span>
          </button>
        </div>
      </div>

      <div v-if="!usesDefaultNS || dnssec.enabled" class="space-y-4 px-5 py-5">
        <!-- Custom nameservers: DNSSEC lives at the provider -->
        <div v-if="!usesDefaultNS" class="flex items-start gap-3 text-sm text-y-muted">
          <Info :size="18" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-info" />
          <span>{{ $t('dnsManagement.dnssecCustomNS') }}</span>
        </div>

        <template v-else>
          <!-- Chain of trust status -->
          <div class="flex flex-wrap gap-2">
            <span :class="dnssec.parentSigned ? 'y-chip-mint' : 'y-chip-warn'">
              {{ dnssec.parentSigned ? $t('dnsManagement.dnssecParentSigned', { zone: dnssec.parentZone }) : $t('uiUser.dns.parentUnsigned', { zone: dnssec.parentZone }) }}
            </span>
            <span :class="dnssec.dsPublished ? 'y-chip-mint' : 'y-chip-warn'">
              {{ dnssec.dsPublished ? $t('dnsManagement.dnssecDSPublished') : $t('dnsManagement.dnssecDSMissing') }}
            </span>
            <span :class="dnssec.chainValid ? 'y-chip-mint' : 'y-chip'">
              {{ dnssec.chainValid ? $t('dnsManagement.dnssecChainValid') : $t('dnsManagement.dnssecChainIncomplete') }}
            </span>
          </div>

          <!-- Keys -->
          <div v-if="dnssec.keys.length > 0" class="space-y-3">
            <div v-for="key in dnssec.keys" :key="key.id" class="y-inset space-y-4 p-4">
              <div class="flex flex-wrap items-center gap-2">
                <span class="y-tag">{{ key.keytype.toUpperCase() }}</span>
                <span class="font-mono text-xs text-y-faint">{{ $t('uiUser.dns.algorithmBits', { algorithm: key.algorithm, bits: key.bits }) }}</span>
              </div>

              <div class="space-y-1.5">
                <p class="y-label">DNSKEY</p>
                <div class="flex items-start gap-2">
                  <code class="min-w-0 flex-1 break-all rounded-sm border border-y-line bg-y-surface px-3 py-2 text-xs text-y-text">{{ key.dnskey }}</code>
                  <button type="button" class="btn btn-ghost btn-sm btn-square shrink-0" :title="$t('uiUser.common.copy')" :aria-label="$t('uiUser.common.copy')" @click="copyText(key.dnskey)">
                    <Copy :size="16" :stroke-width="1.5" />
                  </button>
                </div>
              </div>

              <div v-if="key.ds && key.ds.length > 0" class="space-y-1.5">
                <p class="y-label">{{ $t('uiUser.dns.dsRecords') }}</p>
                <div v-for="(ds, idx) in key.ds" :key="idx" class="flex items-start gap-2">
                  <code class="min-w-0 flex-1 break-all rounded-sm border border-y-line bg-y-surface px-3 py-2 text-xs text-y-text">{{ ds }}</code>
                  <button type="button" class="btn btn-ghost btn-sm btn-square shrink-0" :title="$t('uiUser.common.copy')" :aria-label="$t('uiUser.common.copy')" @click="copyText(ds)">
                    <Copy :size="16" :stroke-width="1.5" />
                  </button>
                </div>
                <button type="button" class="btn btn-outline btn-sm mt-2 gap-1.5" :disabled="dnssec.publishingDS" @click="publishDS">
                  <span v-if="dnssec.publishingDS" class="loading loading-spinner loading-xs"></span>
                  <Upload v-else :size="14" :stroke-width="1.5" />
                  {{ dnssec.dsPublished ? $t('uiUser.dns.republishDS') : $t('uiUser.dns.publishDS') }}
                </button>
              </div>
            </div>
          </div>

          <p v-else-if="!dnssec.loading" class="text-sm text-y-muted">{{ $t('dnsManagement.dnssecNoKeys') }}</p>
        </template>
      </div>
    </section>

    <!-- Add/Edit Modal -->
    <dialog ref="modalRef" class="modal" :class="{ 'modal-open': showAddModal || showEditModal }">
      <div class="modal-box w-11/12 max-w-2xl">
        <h3 class="text-lg font-semibold text-y-text">
          {{ showEditModal ? $t('uiUser.dns.editTitle') : $t('uiUser.dns.addTitle') }}
        </h3>
        <p class="mt-1 text-sm text-y-muted">{{ $t('uiUser.dns.formDesc') }}</p>

        <form class="mt-6 space-y-5" @submit.prevent="handleSubmit">
          <div class="grid gap-5 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('dnsManagement.name') }}</span>
              <input
                v-model="form.name"
                type="text"
                :placeholder="$t('dnsManagement.namePlaceholder')"
                class="input input-bordered w-full font-mono"
                autocomplete="off"
                autocapitalize="off"
                spellcheck="false"
                required
              />
              <span class="y-field-hint">{{ $t('uiUser.dns.nameHint') }}</span>
            </label>

            <label class="y-field">
              <span class="y-field-label">{{ $t('dnsManagement.type') }}</span>
              <select v-model="form.type" class="select select-bordered w-full font-mono" required>
                <option value="A">A</option>
                <option value="AAAA">AAAA</option>
                <option value="CNAME">CNAME</option>
                <option value="ALIAS">ALIAS</option>
                <option value="MX">MX</option>
                <option value="TXT">TXT</option>
                <option value="NS">NS</option>
                <option value="SRV">SRV</option>
                <option value="CAA">CAA</option>
              </select>
            </label>
          </div>

          <div v-if="form.type === 'ALIAS'" class="y-inset flex items-start gap-3 px-4 py-3 text-sm">
            <Info :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-info" />
            <p class="text-y-muted"><span class="font-medium text-y-text">{{ $t('dnsManagement.aliasHintTitle') }}</span> {{ $t('dnsManagement.aliasHint') }}</p>
          </div>
          <div v-if="form.type === 'CNAME' && (form.name === '@' || form.name === '')" class="y-mint-soft flex items-start gap-3 rounded-lg border px-4 py-3 text-sm">
            <Info :size="16" :stroke-width="1.5" class="mt-0.5 shrink-0 text-y-mint-text" />
            <p class="text-y-muted"><span class="font-medium text-y-text">{{ $t('dnsManagement.cnameFlatteningTitle') }}</span> {{ $t('dnsManagement.cnameFlatteningHint') }}</p>
          </div>

          <label class="y-field">
            <span class="y-field-label">{{ $t('dnsManagement.content') }}</span>
            <textarea
              v-model="form.content"
              class="textarea textarea-bordered w-full font-mono text-sm"
              rows="3"
              :placeholder="$t('dnsManagement.contentPlaceholder')"
              spellcheck="false"
              required
            ></textarea>
          </label>

          <div class="grid gap-5 sm:grid-cols-2">
            <label class="y-field">
              <span class="y-field-label">{{ $t('uiUser.dns.ttlLabel') }}</span>
              <input v-model.number="form.ttl" type="number" min="60" max="86400" class="input input-bordered w-full font-mono" required />
              <span class="y-field-hint">{{ $t('uiUser.dns.ttlHint') }}</span>
            </label>

            <label v-if="form.type === 'MX' || form.type === 'SRV'" class="y-field">
              <span class="y-field-label">{{ $t('dnsManagement.priority') }}</span>
              <input v-model.number="form.priority" type="number" min="0" class="input input-bordered w-full font-mono" />
            </label>
          </div>

          <label v-if="showEditModal" class="y-inset flex cursor-pointer items-center justify-between gap-4 px-4 py-3">
            <span>
              <span class="block text-sm font-medium text-y-text">{{ $t('uiUser.dns.active') }}</span>
              <span class="block text-xs text-y-faint">{{ $t('uiUser.dns.activeHint') }}</span>
            </span>
            <input v-model="form.is_active" type="checkbox" class="toggle toggle-sm" />
          </label>

          <div class="modal-action">
            <button type="button" class="btn btn-ghost" @click="closeModal">{{ $t('common.cancel') }}</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              <span v-else>{{ showEditModal ? $t('uiUser.dns.save') : $t('uiUser.dns.addRecord') }}</span>
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
        <h3 class="flex items-center gap-2 text-lg font-semibold text-y-text">
          <AlertTriangle :size="18" :stroke-width="1.5" class="text-y-danger-text" />
          {{ $t('uiUser.dns.syncErrorTitle') }}
        </h3>
        <p class="mt-1 text-sm text-y-muted">{{ $t('uiUser.dns.syncErrorDesc') }}</p>

        <div class="y-inset mt-5 max-h-72 overflow-auto p-4">
          <pre class="whitespace-pre-wrap break-words font-mono text-xs text-y-text">{{ currentError }}</pre>
        </div>

        <div class="modal-action">
          <button type="button" class="btn btn-ghost" @click="closeErrorModal">{{ $t('common.close') }}</button>
          <button type="button" class="btn btn-outline gap-1.5" @click="copyError">
            <Check v-if="errorCopied" :size="16" :stroke-width="1.5" class="text-y-mint-text" />
            <Copy v-else :size="16" :stroke-width="1.5" />
            {{ errorCopied ? $t('uiUser.common.copied') : $t('uiUser.dns.copyError') }}
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeErrorModal">
        <button>close</button>
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
import { ArrowLeft, RefreshCw, Plus, Pencil, Trash2, AlertTriangle, Shield, ShieldCheck, Info, Copy, Check, Upload } from 'lucide-vue-next'
import PageHeader from '../components/ui/PageHeader.vue'
import StatusDot from '../components/ui/StatusDot.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import YLoader from '../components/brand/YLoader.vue'

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

const dnssec = ref({
  enabled: false,
  keys: [],
  parentZone: '',
  parentSigned: false,
  dsPublished: false,
  chainValid: false,
  loading: false,
  publishingDS: false,
})

const usesDefaultNS = computed(() => domain.value?.use_default_nameservers !== false)

const applyDNSSECStatus = (data) => {
  dnssec.value.enabled = !!data.enabled
  dnssec.value.keys = data.keys || []
  dnssec.value.parentZone = data.parent_zone || dnssec.value.parentZone
  dnssec.value.parentSigned = !!data.parent_signed
  dnssec.value.dsPublished = !!data.ds_published
  dnssec.value.chainValid = !!data.chain_valid
}

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
  await Promise.all([fetchRecords(), fetchDNSSEC()])
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

const fetchDNSSEC = async () => {
  try {
    const response = await axios.get(`/api/domains/${domainId}/dnssec`)
    applyDNSSECStatus(response.data)
  } catch (error) {
    console.error('Failed to fetch DNSSEC status:', error)
  }
}

const toggleDNSSEC = async () => {
  dnssec.value.loading = true
  try {
    if (dnssec.value.enabled) {
      if (!confirm(t('dnsManagement.disableDNSSECConfirm'))) {
        dnssec.value.loading = false
        return
      }
      const response = await axios.post(`/api/domains/${domainId}/dnssec/disable`)
      applyDNSSECStatus(response.data)
      toast.success(t('dnsManagement.dnssecDisabledSuccess'))
    } else {
      const response = await axios.post(`/api/domains/${domainId}/dnssec/enable`)
      applyDNSSECStatus(response.data)
      toast.success(t('dnsManagement.dnssecEnabledSuccess'))
      if (response.data.warning) {
        toast.warning(response.data.warning, 6000)
      }
    }
  } catch (error) {
    toast.error(error.response?.data?.error || t('dnsManagement.dnssecOperationFailed'))
  } finally {
    dnssec.value.loading = false
  }
}

const copyText = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    toast.success(t('dnsManagement.copied'))
  } catch {
    toast.error(t('dnsManagement.copyFailed'))
  }
}

const publishDS = async () => {
  dnssec.value.publishingDS = true
  try {
    const response = await axios.post(`/api/domains/${domainId}/dnssec/publish-ds`)
    applyDNSSECStatus(response.data)
    toast.success(response.data.message || t('dnsManagement.publishDSSuccess'))
  } catch (error) {
    toast.error(error.response?.data?.error || t('dnsManagement.publishDSFailed'))
  } finally {
    dnssec.value.publishingDS = false
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
