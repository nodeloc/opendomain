<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-8">
    <div class="flex justify-between items-center">
      <h1 class="text-4xl font-bold">{{ $t('domains.myDomains') }}</h1>
      <div class="flex gap-2">
        <button v-if="fossbillingEnabled" @click="openSyncModal" class="btn btn-outline btn-secondary">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Sync from FOSSBilling
        </button>
        <router-link to="/" class="btn btn-primary">
          {{ $t('domains.registerNew') }}
        </router-link>
      </div>
    </div>

    <!-- Domains List -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="domains.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 opacity-40 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
        </svg>
        <h3 class="card-title text-2xl mb-2">{{ $t('domains.noDomains') }}</h3>
        <p class="mb-4">{{ $t('domains.noDomainsDesc') }}</p>
        <div class="flex gap-3">
          <router-link to="/" class="btn btn-primary">{{ $t('domains.registerDomain') }}</router-link>
          <button v-if="fossbillingEnabled" @click="openSyncModal" class="btn btn-outline btn-secondary">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Sync from FOSSBilling
          </button>
        </div>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="domain in domains" :key="domain.id" class="card bg-base-200 shadow-xl hover:shadow-2xl transition-all" :class="{ 'opacity-60': domain.status === 'suspended' }">
        <div class="card-body">
          <h2 class="card-title text-lg break-all">
            {{ domain.full_domain }}
            <div class="px-2 py-0.5 rounded text-xs font-medium" :class="getStatusClass(domain.status)">
              {{ domain.status }}
            </div>
          </h2>

          <!-- Suspended Notice -->
          <div v-if="domain.status === 'suspended'" class="alert alert-warning py-2 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
            <span>{{ $t('domains.suspendedNotice') }}</span>
          </div>

          <div class="space-y-2 text-sm">
            <p>
              <span class="font-semibold">{{ $t('domains.registered') }}:</span>
              {{ formatDate(domain.registered_at) }}
            </p>
            <p>
              <span class="font-semibold">{{ $t('domains.expires') }}:</span>
              {{ formatDate(domain.expires_at) }}
              <span v-if="domain.status === 'active' && getDaysUntilExpiry(domain.expires_at) <= 30" class="text-warning ml-2">
                ({{ getDaysUntilExpiry(domain.expires_at) }} {{ $t('domains.daysLeft') }})
              </span>
            </p>
            <p>
              <span class="font-semibold">{{ $t('domains.dnsSynced') }}:</span>
              <span :class="domain.dns_synced ? 'text-success' : 'text-warning'">
                <svg v-if="domain.dns_synced" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                {{ domain.dns_synced ? $t('domains.yes') : $t('domains.no') }}
              </span>
            </p>
          </div>

          <div v-if="domain.status !== 'suspended'" class="card-actions justify-end mt-4 flex-wrap gap-2">
            <div class="tooltip" :data-tip="isUsingDefaultNS(domain) ? $t('domains.manageDNS') : $t('domains.dnsTooltip')">
              <button
                class="btn btn-sm btn-primary"
                @click="manageDNS(domain)"
                :disabled="!isUsingDefaultNS(domain)"
                :class="{ 'btn-disabled': !isUsingDefaultNS(domain) }"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                DNS
              </button>
            </div>

            <!-- More Actions Dropdown -->
            <div class="dropdown dropdown-end">
              <label tabindex="0" class="btn btn-sm btn-ghost">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                </svg>
              </label>
              <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                <li><a @click="openModifyNS(domain)">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
                  </svg>
                  Modify Nameservers
                </a></li>
                <li><a @click="openRenew(domain)">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  Renew Domain
                </a></li>
                <li><a @click="openTransfer(domain)">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                  </svg>
                  Transfer Domain
                </a></li>
                <li><a @click="confirmDelete(domain)" class="text-error">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Delete
                </a></li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modify Nameservers Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showNSModal }">
      <div class="modal-box" v-if="selectedDomain">
        <h3 class="font-bold text-2xl mb-4">Modify Nameservers</h3>
        <p class="text-sm opacity-70 mb-4">Domain: <span class="font-mono font-bold">{{ selectedDomain.full_domain }}</span></p>

        <div class="form-control mb-4">
          <label class="label cursor-pointer justify-start gap-3">
            <input
              v-model="useCustomNS"
              type="checkbox"
              class="toggle toggle-primary"
            />
            <span class="label-text font-semibold">Using Custom Nameservers</span>
          </label>
          <label class="label">
            <span class="label-text-alt opacity-60">Toggle off to use system default nameservers</span>
          </label>
        </div>

        <!-- Default Nameservers (Read-only) -->
        <div v-if="!useCustomNS" class="form-control space-y-3">
          <div v-for="(ns, index) in ['ns1.nodelook.com', 'ns2.nodelook.com']" :key="index" class="flex gap-2">
            <input
              :value="ns"
              type="text"
              class="input input-bordered flex-1 bg-base-200"
              readonly
              disabled
            />
            <button class="btn btn-square btn-disabled opacity-50">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Custom Nameservers (Editable) -->
        <div v-else class="form-control space-y-3">
          <div v-for="(ns, index) in nameservers" :key="index" class="flex gap-2">
            <input
              v-model="nameservers[index]"
              type="text"
              :placeholder="`Nameserver ${index + 1}`"
              class="input input-bordered flex-1"
            />
            <button
              v-if="nameservers.length > 1"
              class="btn btn-error btn-square"
              @click="removeNameserver(index)"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <button class="btn btn-sm btn-ghost" @click="addNameserver">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Nameserver
          </button>
        </div>

        <div class="modal-action">
          <button class="btn btn-ghost" @click="closeNSModal">Cancel</button>
          <button
            class="btn btn-primary"
            @click="saveNameservers"
            :disabled="submitting"
          >
            <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
            <span v-else>Save</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeNSModal">close</button>
      </form>
    </dialog>

    <!-- Renew Domain Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showRenewModal }">
      <div class="modal-box" v-if="selectedDomain">
        <h3 class="font-bold text-2xl mb-4">Renew Domain</h3>
        <p class="text-sm opacity-70 mb-4">Domain: <span class="font-mono font-bold">{{ selectedDomain.full_domain }}</span></p>

        <div class="alert alert-info mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <div class="font-bold">Current Expiration</div>
            <div class="text-sm">{{ formatDate(selectedDomain.expires_at) }}</div>
          </div>
        </div>

        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text font-semibold">Renewal Period</span>
          </label>
          <input
            v-model.number="renewYears"
            type="range"
            min="1"
            max="10"
            class="range range-primary"
          />
          <div class="flex justify-between text-xs px-2 mt-2">
            <span>1 Year</span>
            <span>{{ renewYears }} {{ renewYears === 1 ? 'Year' : 'Years' }}</span>
            <span>10 Years</span>
          </div>
        </div>

        <div class="alert alert-success">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <div class="font-bold">New Expiration</div>
            <div class="text-sm">{{ calculateNewExpiry(selectedDomain.expires_at, renewYears) }}</div>
          </div>
        </div>

        <div v-if="selectedDomain.root_domain && !selectedDomain.root_domain.is_free" class="mt-4">
          <div class="text-lg font-bold text-center">
            Price: {{ formatPrice(calculateRenewPrice()) }}
          </div>
        </div>

        <div class="modal-action">
          <button class="btn btn-ghost" @click="closeRenewModal">Cancel</button>
          <button class="btn btn-primary" @click="renewDomain" :disabled="submitting">
            <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
            <span v-else>Renew</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeRenewModal">close</button>
      </form>
    </dialog>

    <!-- Transfer Domain Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showTransferModal }">
      <div class="modal-box" v-if="selectedDomain">
        <h3 class="font-bold text-2xl mb-4">Transfer Domain</h3>
        <p class="text-sm opacity-70 mb-4">Domain: <span class="font-mono font-bold">{{ selectedDomain.full_domain }}</span></p>

        <div class="alert alert-warning mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <div>
            <div class="font-bold">Warning</div>
            <div class="text-sm">This action will transfer domain ownership. You will lose access to this domain.</div>
          </div>
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text font-semibold">Recipient Email or Username</span>
          </label>
          <input
            v-model="transferTarget"
            type="text"
            placeholder="Enter email or username"
            class="input input-bordered"
          />
          <label class="label">
            <span class="label-text-alt opacity-60">The user must have an account on this platform</span>
          </label>
        </div>

        <div class="modal-action">
          <button class="btn btn-ghost" @click="closeTransferModal">Cancel</button>
          <button class="btn btn-error" @click="transferDomain" :disabled="submitting || !transferTarget">
            <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
            <span v-else>Transfer</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeTransferModal">close</button>
      </form>
    </dialog>

    <!-- FOSSBilling Sync Modal -->
    <dialog :class="{ 'modal': true, 'modal-open': showSyncModal }">
      <div class="modal-box max-w-2xl">
        <h3 class="font-bold text-2xl mb-4">Sync Domains from FOSSBilling</h3>
        <p class="text-sm opacity-70 mb-6">Enter your FOSSBilling credentials to import your domains</p>

        <div v-if="!siteConfig.fossbilling?.url" class="form-control mb-4">
          <label class="label">
            <span class="label-text font-semibold">FOSSBilling Server URL</span>
          </label>
          <input
            v-model="syncForm.fossbilling_url"
            type="text"
            placeholder="https://your-fossbilling-server.com"
            class="input input-bordered"
            :disabled="syncInProgress"
          />
          <label class="label">
            <span class="label-text-alt opacity-60">Full URL including https://</span>
          </label>
        </div>

        <div v-if="siteConfig.fossbilling?.url" class="alert alert-info mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-sm">Syncing from: <strong>{{ siteConfig.fossbilling.url }}</strong></span>
        </div>

        <!-- API Key Input -->
        <div class="form-control mb-6">
          <label class="label">
            <span class="label-text font-semibold">FOSSBilling API Key</span>
          </label>
          <input
            v-model="syncForm.fossbilling_api_key"
            type="text"
            placeholder="Enter your FOSSBilling API Key"
            class="input input-bordered font-mono"
            :disabled="syncInProgress"
          />
          <label class="label">
            <span class="label-text-alt opacity-60">Get your API key from FOSSBilling admin panel</span>
          </label>
        </div>

        <!-- Sync Results -->
        <div v-if="syncResult" class="space-y-4 mb-4">
          <div class="alert" :class="syncResult.error_count > 0 ? 'alert-warning' : 'alert-success'">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div>
              <div class="font-bold">{{ syncResult.message }}</div>
              <div class="text-sm">
                Synced: {{ syncResult.synced_count }} |
                Existing: {{ syncResult.existing_count }} |
                Skipped: {{ syncResult.skipped_count }} |
                Errors: {{ syncResult.error_count }}
              </div>
            </div>
          </div>

          <div v-if="syncResult.details && syncResult.details.length > 0" class="bg-base-200 p-4 rounded-lg max-h-60 overflow-y-auto">
            <div class="text-sm font-semibold mb-2">Details:</div>
            <ul class="text-xs space-y-1">
              <li v-for="(detail, index) in syncResult.details" :key="index" class="opacity-80">
                {{ detail }}
              </li>
            </ul>
          </div>
        </div>

        <div class="modal-action">
          <button class="btn btn-ghost" @click="closeSyncModal" :disabled="syncInProgress">Cancel</button>
          <button
            class="btn btn-primary"
            @click="syncFromFOSSBilling"
            :disabled="syncInProgress || (!syncForm.fossbilling_url && !siteConfig.fossbilling?.url) || !syncForm.fossbilling_api_key"
          >
            <span v-if="syncInProgress" class="loading loading-spinner loading-sm"></span>
            <span v-else>Sync Domains</span>
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button type="button" @click="closeSyncModal" :disabled="syncInProgress">close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'
import { useSiteConfigStore } from '../stores/siteConfig'

const router = useRouter()
const toast = useToast()
const { formatPrice } = useCurrency()
const siteConfig = useSiteConfigStore()

// 确保加载站点配置
if (!siteConfig.loaded) {
  siteConfig.fetch()
}

// FOSSBilling 是否启用
const fossbillingEnabled = computed(() => siteConfig.fossbilling?.enabled || false)

const domains = ref([])
const loading = ref(true)
const submitting = ref(false)

// Modal states
const showNSModal = ref(false)
const showRenewModal = ref(false)
const showTransferModal = ref(false)
const showSyncModal = ref(false)
const selectedDomain = ref(null)

// NS modification
const nameservers = ref(['', ''])
const useCustomNS = ref(false)

// Renew
const renewYears = ref(1)

// Transfer
const transferTarget = ref('')

// FOSSBilling Sync
const syncInProgress = ref(false)
const syncForm = ref({
  fossbilling_url: '',
  fossbilling_api_key: ''
})
const syncResult = ref(null)

onMounted(async () => {
  await fetchDomains()
})

const fetchDomains = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/domains')
    domains.value = response.data.domains
  } catch (error) {
    console.error('Failed to fetch domains:', error)
  } finally {
    loading.value = false
  }
}

const getStatusClass = (status) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'expired':
      return 'bg-red-100 text-red-800'
    case 'suspended':
      return 'bg-amber-100 text-amber-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

const getDaysUntilExpiry = (expiryDate) => {
  const now = new Date()
  const expiry = new Date(expiryDate)
  const diff = expiry - now
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
}

const isUsingDefaultNS = (domain) => {
  return domain.use_default_nameservers
}

const manageDNS = (domain) => {
  if (!isUsingDefaultNS(domain)) {
    toast.warning('DNS management is only available for domains using default nameservers')
    return
  }
  router.push(`/domains/${domain.id}/dns`)
}

// Nameserver Management
const openModifyNS = (domain) => {
  selectedDomain.value = domain

  // Parse existing nameservers from domain
  if (domain.nameservers) {
    try {
      const ns = JSON.parse(domain.nameservers)
      nameservers.value = ns.length > 0 ? ns : ['', '']
      useCustomNS.value = !domain.use_default_nameservers
    } catch (e) {
      nameservers.value = ['', '']
      useCustomNS.value = true
    }
  } else {
    nameservers.value = ['', '']
    useCustomNS.value = true
  }

  showNSModal.value = true
}

const addNameserver = () => {
  nameservers.value.push('')
}

const removeNameserver = (index) => {
  nameservers.value.splice(index, 1)
}

const saveNameservers = async () => {
  let nsToSave = []

  if (useCustomNS.value) {
    // Use custom nameservers
    nsToSave = nameservers.value.filter(ns => ns.trim() !== '')

    if (nsToSave.length === 0) {
      toast.warning('Please add at least one nameserver')
      return
    }
  } else {
    // Use default nameservers
    nsToSave = ['ns1.nodelook.com', 'ns2.nodelook.com']
  }

  submitting.value = true
  try {
    await axios.put(`/api/domains/${selectedDomain.value.id}/nameservers`, {
      nameservers: nsToSave
    })
    toast.success('Nameservers updated successfully')
    closeNSModal()
    await fetchDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to update nameservers')
  } finally {
    submitting.value = false
  }
}

const closeNSModal = () => {
  showNSModal.value = false
  selectedDomain.value = null
  nameservers.value = ['', '']
  useCustomNS.value = false
}

// Renew Domain
const openRenew = (domain) => {
  selectedDomain.value = domain
  renewYears.value = 1
  showRenewModal.value = true
}

const calculateNewExpiry = (currentExpiry, years) => {
  const expiry = new Date(currentExpiry)
  expiry.setFullYear(expiry.getFullYear() + years)
  return formatDate(expiry)
}

const calculateRenewPrice = () => {
  if (!selectedDomain.value?.root_domain?.price_per_year) return 0
  return selectedDomain.value.root_domain.price_per_year * renewYears.value
}

const renewDomain = async () => {
  submitting.value = true
  try {
    await axios.post(`/api/domains/${selectedDomain.value.id}/renew`, {
      years: renewYears.value
    })
    toast.success(`Domain renewed for ${renewYears.value} year(s)`)
    closeRenewModal()
    await fetchDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to renew domain')
  } finally {
    submitting.value = false
  }
}

const closeRenewModal = () => {
  showRenewModal.value = false
  selectedDomain.value = null
  renewYears.value = 1
}

// Transfer Domain
const openTransfer = (domain) => {
  selectedDomain.value = domain
  transferTarget.value = ''
  showTransferModal.value = true
}

const transferDomain = async () => {
  if (!transferTarget.value.trim()) {
    toast.warning('Please enter recipient email or username')
    return
  }

  if (!confirm(`Are you sure you want to transfer ${selectedDomain.value.full_domain} to ${transferTarget.value}?`)) {
    return
  }

  submitting.value = true
  try {
    await axios.post(`/api/domains/${selectedDomain.value.id}/transfer`, {
      target: transferTarget.value
    })
    toast.success('Domain transferred successfully')
    closeTransferModal()
    await fetchDomains()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to transfer domain')
  } finally {
    submitting.value = false
  }
}

const closeTransferModal = () => {
  showTransferModal.value = false
  selectedDomain.value = null
  transferTarget.value = ''
}

// Delete Domain
const confirmDelete = async (domain) => {
  if (confirm(`Are you sure you want to delete ${domain.full_domain}? This action cannot be undone.`)) {
    try {
      await axios.delete(`/api/domains/${domain.id}`)
      toast.success('Domain deleted successfully')
      await fetchDomains()
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to delete domain')
    }
  }
}

// FOSSBilling Sync
const openSyncModal = () => {
  // 使用配置中的默认 URL（如果有）
  syncForm.value = {
    fossbilling_url: siteConfig.fossbilling?.url || '',
    fossbilling_api_key: ''
  }
  syncResult.value = null
  showSyncModal.value = true
}

const syncFromFOSSBilling = async () => {
  // 验证必填字段
  if (!syncForm.value.fossbilling_url && !siteConfig.fossbilling?.url) {
    toast.warning('Please enter FOSSBilling server URL')
    return
  }

  if (!syncForm.value.fossbilling_api_key) {
    toast.warning('Please enter FOSSBilling API Key')
    return
  }

  syncInProgress.value = true
  syncResult.value = null

  try {
    const response = await axios.post('/api/user/sync-from-fossbilling', syncForm.value)
    syncResult.value = response.data

    if (response.data.synced_count > 0) {
      toast.success(`Successfully synced ${response.data.synced_count} domain(s)`)
      // Refresh domain list after sync
      await fetchDomains()
    } else {
      toast.info('No new domains to sync')
    }
  } catch (error) {
    toast.error(error.response?.data?.error || 'Failed to sync from FOSSBilling')
    syncResult.value = {
      message: 'Sync failed',
      error_count: 1,
      synced_count: 0,
      existing_count: 0,
      skipped_count: 0,
      details: [error.response?.data?.error || 'An error occurred during sync']
    }
  } finally {
    syncInProgress.value = false
  }
}

const closeSyncModal = () => {
  if (syncInProgress.value) {
    return
  }
  showSyncModal.value = false
  syncForm.value = {
    fossbilling_url: '',
    fossbilling_api_key: ''
  }
  syncResult.value = null
}
</script>
