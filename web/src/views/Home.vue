<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section class="w-full min-h-[80vh] bg-gradient-to-b from-base-100 via-base-200 to-base-100 relative overflow-hidden flex items-center justify-center">
      <!-- Animated background elements -->
      <div class="absolute inset-0 opacity-10">
        <div class="absolute top-20 left-10 w-72 h-72 bg-primary rounded-full blur-3xl animate-blob"></div>
        <div class="absolute top-40 right-10 w-72 h-72 bg-secondary rounded-full blur-3xl animate-blob animation-delay-2000"></div>
        <div class="absolute -bottom-20 left-1/2 w-72 h-72 bg-accent rounded-full blur-3xl animate-blob animation-delay-4000"></div>
      </div>

      <div class="w-full text-center relative z-10 py-20">
        <div class="max-w-5xl mx-auto px-4">
          <div class="mb-10 space-y-6 animate-fade-in">
            <h1 class="text-5xl md:text-7xl font-extrabold leading-tight">
              {{ $t('home.heroTitle') }}
              <span class="block mt-2 bg-gradient-to-r from-primary via-secondary to-accent bg-clip-text text-transparent animate-gradient">
                {{ $t('home.heroSubtitle') }}
              </span>
              {{ $t('home.heroToday') }}
            </h1>
            <p class="text-xl md:text-2xl opacity-70 max-w-3xl mx-auto">
              {{ $t('home.heroDescription') }}
            </p>
            <div class="flex flex-wrap items-center justify-center gap-4 text-sm md:text-base opacity-60">
              <span class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                {{ $t('home.featureNoCredit') }}
              </span>
              <span class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                {{ $t('home.featureInstant') }}
              </span>
              <span class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                {{ $t('home.featureFree') }}
              </span>
            </div>
          </div>

          <!-- Domain Search Card -->
          <div class="card bg-base-100 shadow-2xl p-6 md:p-10 border border-base-300 backdrop-blur-sm hover:shadow-3xl transition-all duration-300 animate-slide-up">
            <form @submit.prevent="searchDomain" class="space-y-6">
              <div class="flex flex-col md:flex-row gap-3">
                <div class="flex-1 relative group">
                  <input
                    v-model="searchQuery"
                    type="text"
                    :placeholder="$t('home.searchPlaceholder')"
                    class="input input-bordered input-lg w-full pr-12 focus:input-primary focus:scale-[1.02] transition-all duration-200"
                    required
                  />
                  <div class="absolute right-4 top-1/2 transform -translate-y-1/2 opacity-40 group-focus-within:opacity-70 transition-opacity">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                    </svg>
                  </div>
                </div>
                <select v-model="selectedRootDomain" class="select select-bordered select-lg md:w-52 focus:select-primary focus:scale-[1.02] transition-all duration-200">
                  <option value="" disabled>{{ $t('home.chooseExtension') }}</option>
                  <option v-for="domain in rootDomains" :key="domain.id" :value="domain.id">
                    .{{ domain.domain }}
                  </option>
                </select>
                <button
                  type="submit"
                  class="btn btn-primary btn-lg md:w-40 gap-2 hover:scale-105 active:scale-95 transition-transform"
                  :disabled="loading"
                >
                  <span v-if="loading" class="loading loading-spinner"></span>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  {{ $t('home.search') }}
                </button>
              </div>
            </form>

            <!-- Search Result -->
            <transition name="slide-fade">
              <div v-if="searchResult" class="mt-6">
                <div v-if="searchResult.available" class="alert alert-success shadow-lg border-2 border-success/20">
                  <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <div class="flex-1">
                    <h3 class="font-bold text-lg flex items-center gap-2">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      {{ searchResult.message }}
                    </h3>
                    <p class="text-sm opacity-80 mt-1">
                      {{ searchResult.isFree ? $t('home.claimNow') : `${$t('home.startingFrom')} ${currencySymbol}${searchResult.pricePerYear || searchResult.lifetimePrice}/${$t('home.year')}` }}
                    </p>
                  </div>
                  <button
                    v-if="!isAuthenticated"
                    @click="$router.push('/register')"
                    class="btn btn-sm btn-success gap-2 hover:scale-105 transition-transform"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                    </svg>
                    {{ $t('home.signupToRegister') }}
                  </button>
                  <button
                    v-else
                    @click="registerDomain"
                    class="btn btn-sm btn-success gap-2 hover:scale-105 transition-transform"
                  >
                    <svg v-if="searchResult.isFree" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    {{ searchResult.isFree ? $t('home.registerFree') : $t('home.proceedToCheckout') }}
                  </button>
                </div>
                <div v-else class="alert alert-warning shadow-lg border-2 border-warning/20">
                  <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <div class="flex-1">
                    <h3 class="font-bold">{{ searchResult.message }}</h3>
                    <p class="text-sm opacity-80 mt-1">{{ $t('home.tryDifferent') }}</p>
                  </div>
                </div>
              </div>
            </transition>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="w-full py-20 bg-base-100">
      <div class="max-w-7xl mx-auto px-4">
        <div class="text-center mb-16 space-y-4">
          <h2 class="text-4xl md:text-5xl font-bold">{{ $t('home.whyChoose', { siteName: siteConfigStore.siteName }) }}</h2>
          <p class="text-xl md:text-2xl opacity-60 max-w-2xl mx-auto">
            {{ $t('home.whyChooseDesc') }}
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Lightning Fast -->
          <div class="card bg-base-200/50 shadow-xl hover:shadow-2xl hover:scale-105 transition-all duration-300 border border-base-300 backdrop-blur-sm group">
            <div class="card-body items-center text-center">
              <div class="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 class="card-title text-2xl mb-2">{{ $t('home.fastTitle') }}</h3>
              <p class="opacity-70">{{ $t('home.fastDesc') }}</p>
            </div>
          </div>

          <!-- Secure & Reliable -->
          <div class="card bg-base-200/50 shadow-xl hover:shadow-2xl hover:scale-105 transition-all duration-300 border border-base-300 backdrop-blur-sm group">
            <div class="card-body items-center text-center">
              <div class="w-20 h-20 rounded-full bg-secondary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
              <h3 class="card-title text-2xl mb-2">{{ $t('home.secureTitle') }}</h3>
              <p class="opacity-70">{{ $t('home.secureDesc') }}</p>
            </div>
          </div>

          <!-- Easy Management -->
          <div class="card bg-base-200/50 shadow-xl hover:shadow-2xl hover:scale-105 transition-all duration-300 border border-base-300 backdrop-blur-sm group">
            <div class="card-body items-center text-center">
              <div class="w-20 h-20 rounded-full bg-accent/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                </svg>
              </div>
              <h3 class="card-title text-2xl mb-2">{{ $t('home.easyTitle') }}</h3>
              <p class="opacity-70">{{ $t('home.easyDesc') }}</p>
            </div>
          </div>

          <!-- Multiple Domains -->
          <div class="card bg-base-200/50 shadow-xl hover:shadow-2xl hover:scale-105 transition-all duration-300 border border-base-300 backdrop-blur-sm group">
            <div class="card-body items-center text-center">
              <div class="w-20 h-20 rounded-full bg-primary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="card-title text-2xl mb-2">{{ $t('home.multipleTitle') }}</h3>
              <p class="opacity-70">{{ $t('home.multipleDesc') }}</p>
            </div>
          </div>

          <!-- DNS Management -->
          <div class="card bg-base-200/50 shadow-xl hover:shadow-2xl hover:scale-105 transition-all duration-300 border border-base-300 backdrop-blur-sm group">
            <div class="card-body items-center text-center">
              <div class="w-20 h-20 rounded-full bg-secondary/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-secondary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
                </svg>
              </div>
              <h3 class="card-title text-2xl mb-2">{{ $t('home.dnsTitle') }}</h3>
              <p class="opacity-70">{{ $t('home.dnsDesc') }}</p>
            </div>
          </div>

          <!-- Free Forever -->
          <div class="card bg-base-200/50 shadow-xl hover:shadow-2xl hover:scale-105 transition-all duration-300 border border-base-300 backdrop-blur-sm group">
            <div class="card-body items-center text-center">
              <div class="w-20 h-20 rounded-full bg-accent/10 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                </svg>
              </div>
              <h3 class="card-title text-2xl mb-2">{{ $t('home.freeTitle') }}</h3>
              <p class="opacity-70">{{ $t('home.freeDesc') }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats Section -->
    <section class="w-full py-20 bg-base-200/50">
      <div class="max-w-7xl mx-auto px-4">
        <div class="stats stats-vertical lg:stats-horizontal shadow-2xl w-full bg-base-100 border border-base-300">
          <div class="stat hover:bg-base-200 transition-colors duration-300">
            <div class="stat-figure text-primary">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-10 h-10 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
              </svg>
            </div>
            <div class="stat-title text-lg">{{ $t('home.totalDomains') }}</div>
            <div class="stat-value text-primary text-5xl">10K+</div>
            <div class="stat-desc text-base">{{ $t('home.growing') }}</div>
          </div>

          <div class="stat hover:bg-base-200 transition-colors duration-300">
            <div class="stat-figure text-secondary">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-10 h-10 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
            <div class="stat-title text-lg">{{ $t('home.activeUsers') }}</div>
            <div class="stat-value text-secondary text-5xl">5K+</div>
            <div class="stat-desc text-base">{{ $t('home.newUsers') }}</div>
          </div>

          <div class="stat hover:bg-base-200 transition-colors duration-300">
            <div class="stat-figure text-accent">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-10 h-10 stroke-current">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
              </svg>
            </div>
            <div class="stat-title text-lg">{{ $t('home.uptime') }}</div>
            <div class="stat-value text-accent text-5xl">99.9%</div>
            <div class="stat-desc text-base">{{ $t('home.industryLeading') }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="w-full py-20 bg-gradient-to-r from-primary via-secondary to-accent relative overflow-hidden">
      <div class="absolute inset-0 opacity-10">
        <div class="absolute top-0 left-0 w-96 h-96 bg-white rounded-full blur-3xl animate-blob"></div>
        <div class="absolute bottom-0 right-0 w-96 h-96 bg-white rounded-full blur-3xl animate-blob animation-delay-2000"></div>
      </div>
      <div class="max-w-4xl mx-auto px-4 text-center relative z-10">
        <h2 class="text-4xl md:text-5xl font-bold mb-6 text-white">{{ $t('home.readyTitle') }}</h2>
        <p class="text-xl md:text-2xl mb-10 text-white/90 max-w-2xl mx-auto">
          {{ $t('home.readyDesc', { siteName: siteConfigStore.siteName }) }}
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <router-link
            to="/register"
            class="btn btn-lg bg-white text-primary hover:bg-base-100 hover:scale-105 transition-transform shadow-xl border-0"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            {{ $t('home.getStarted') }}
          </router-link>
          <router-link
            to="/announcements"
            class="btn btn-lg btn-outline text-white border-white hover:bg-white hover:text-primary hover:border-white hover:scale-105 transition-transform"
          >
            {{ $t('home.learnMore') }}
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'
import { useCurrency } from '../composables/useCurrency'

const router = useRouter()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const toast = useToast()
const { currencySymbol } = useCurrency()

const searchQuery = ref('')
const selectedRootDomain = ref('')
const rootDomains = ref([])
const searchResult = ref(null)
const loading = ref(false)

const isAuthenticated = computed(() => authStore.isAuthenticated)

onMounted(async () => {
  await fetchRootDomains()
})

const fetchRootDomains = async () => {
  try {
    const response = await axios.get('/api/public/root-domains')
    rootDomains.value = response.data.root_domains
    if (rootDomains.value.length > 0) {
      selectedRootDomain.value = rootDomains.value[0].id
    }
  } catch (error) {
    console.error('Failed to fetch root domains:', error)
  }
}

const searchDomain = async () => {
  if (!searchQuery.value || !selectedRootDomain.value) return

  loading.value = true
  searchResult.value = null

  try {
    const response = await axios.get('/api/domains/search', {
      params: {
        subdomain: searchQuery.value,
        root_domain_id: selectedRootDomain.value,
      },
    })

    // 获取选中的根域名信息
    const rootDomain = rootDomains.value.find(d => d.id === selectedRootDomain.value)

    searchResult.value = {
      available: response.data.available,
      message: response.data.available
        ? `${response.data.full_domain} is available!`
        : `${response.data.full_domain} is already taken.`,
      fullDomain: response.data.full_domain,
      rootDomain: rootDomain,
      isFree: rootDomain?.is_free ?? true,
      pricePerYear: rootDomain?.price_per_year,
      lifetimePrice: rootDomain?.lifetime_price,
    }
  } catch (error) {
    searchResult.value = {
      available: false,
      message: error.response?.data?.error || 'Search failed',
    }
  } finally {
    loading.value = false
  }
}

const registerDomain = async () => {
  if (!isAuthenticated.value) {
    router.push('/register')
    return
  }

  // 检查是否为付费域名
  if (searchResult.value && !searchResult.value.isFree) {
    // 跳转到结算页面
    router.push({
      path: '/checkout',
      query: {
        subdomain: searchQuery.value,
        root_domain_id: selectedRootDomain.value,
      },
    })
    return
  }

  // 免费域名直接注册
  try {
    const response = await axios.post('/api/domains', {
      subdomain: searchQuery.value,
      root_domain_id: selectedRootDomain.value,
    })

    toast.success('Domain registered successfully!')
    router.push('/domains')
  } catch (error) {
    const errorData = error.response?.data
    // 如果是付费域名错误，跳转到结算页
    if (errorData?.requires_payment) {
      router.push({
        path: '/checkout',
        query: {
          subdomain: searchQuery.value,
          root_domain_id: selectedRootDomain.value,
        },
      })
    } else {
      toast.error(errorData?.error || 'Registration failed')
    }
  }
}
</script>

<style scoped>
/* Animations */
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
  100% {
    transform: translate(0px, 0px) scale(1);
  }
}

@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.animate-fade-in {
  animation: fade-in 0.8s ease-out;
}

.animate-slide-up {
  animation: slide-up 0.8s ease-out 0.2s both;
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}

.animate-gradient {
  background-size: 200% 200%;
  animation: gradient 3s ease infinite;
}

/* Transitions */
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

/* Custom shadows */
.shadow-3xl {
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}
</style>
