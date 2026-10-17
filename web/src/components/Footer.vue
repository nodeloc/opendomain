<template>
  <footer class="mt-24 border-t border-y-line">
    <div class="y-container grid gap-10 py-14 md:grid-cols-[1.4fr_1fr_1fr_1fr]">
      <div class="space-y-4">
        <router-link to="/" class="inline-block" :aria-label="siteName">
          <YLogo class="h-7 w-auto" />
        </router-link>
        <p class="max-w-xs text-sm leading-6 text-y-muted">
          {{ siteConfigStore.siteDescription || $t('shell.tagline') }}
        </p>
        <p class="font-mono text-xs text-y-faint">*.yunamo.com</p>
      </div>

      <nav class="space-y-3">
        <h6 class="y-label">{{ $t('footer.services') }}</h6>
        <ul class="space-y-2 text-sm">
          <li><router-link to="/" class="text-y-muted hover:text-y-text">{{ $t('shell.claimName') }}</router-link></li>
          <li><router-link to="/domains" class="text-y-muted hover:text-y-text">{{ $t('shell.dnsRecords') }}</router-link></li>
          <li><router-link to="/domain-health" class="text-y-muted hover:text-y-text">{{ $t('nav.domainHealth') }}</router-link></li>
          <li><router-link to="/announcements" class="text-y-muted hover:text-y-text">{{ $t('nav.announcements') }}</router-link></li>
        </ul>
      </nav>

      <nav class="space-y-3">
        <h6 class="y-label">{{ $t('footer.company') }}</h6>
        <ul class="space-y-2 text-sm">
          <li v-for="page in companyPages" :key="page.id">
            <router-link :to="`/pages/${page.slug}`" class="text-y-muted hover:text-y-text">{{ page.title }}</router-link>
          </li>
          <li v-if="loaded && companyPages.length === 0" class="text-y-faint">—</li>
        </ul>
      </nav>

      <nav class="space-y-3">
        <h6 class="y-label">{{ $t('footer.resources') }}</h6>
        <ul class="space-y-2 text-sm">
          <li v-for="page in resourcePages" :key="page.id">
            <router-link :to="`/pages/${page.slug}`" class="text-y-muted hover:text-y-text">{{ page.title }}</router-link>
          </li>
          <li v-if="loaded && resourcePages.length === 0" class="text-y-faint">—</li>
        </ul>
      </nav>
    </div>

    <div class="border-t border-y-line">
      <div class="y-container flex flex-col gap-2 py-5 text-xs text-y-faint sm:flex-row sm:items-center sm:justify-between">
        <p>© {{ year }} {{ siteName }}. {{ $t('shell.rights') }}</p>
        <p class="font-mono">yunamo.com</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from '../utils/axios'
import { useSiteConfigStore } from '../stores/siteConfig'
import YLogo from './brand/YLogo.vue'

const siteConfigStore = useSiteConfigStore()
const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const year = new Date().getFullYear()

const pages = ref([])
const loaded = ref(false)

const companyPages = computed(() => pages.value.filter((page) => page.category === 'company'))
const resourcePages = computed(() => pages.value.filter((page) => page.category === 'resources'))

const fetchPages = async () => {
  try {
    const response = await axios.get('/api/public/pages')
    pages.value = response.data.pages || []
  } catch (error) {
    console.error('Failed to fetch pages:', error)
  } finally {
    loaded.value = true
  }
}

onMounted(fetchPages)
</script>
