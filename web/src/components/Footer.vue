<template>
  <footer class="footer p-10 bg-base-200 text-base-content mt-20">
    <aside>
      <p class="font-bold text-lg">{{ siteConfigStore.siteName || 'OpenDomain' }}</p>
      <p>{{ siteConfigStore.siteDescription || $t('footer.description') }}</p>
      <p class="text-sm opacity-70">© {{ new Date().getFullYear() }} {{ siteConfigStore.siteName || 'OpenDomain' }}. All rights reserved.</p>
    </aside>
    <nav>
      <h6 class="footer-title">{{ $t('footer.services') }}</h6>
      <router-link to="/" class="link link-hover">Domain Registration</router-link>
      <router-link to="/domains" class="link link-hover">DNS Management</router-link>
      <router-link to="/domain-health" class="link link-hover">Domain Scanner</router-link>
    </nav>
    <nav>
      <h6 class="footer-title">{{ $t('footer.company') }}</h6>
      <router-link 
        v-for="page in companyPages" 
        :key="page.id" 
        :to="`/pages/${page.slug}`" 
        class="link link-hover"
      >
        {{ page.title }}
      </router-link>
      <a v-if="companyPages.length === 0" class="link link-hover opacity-50">Loading...</a>
    </nav>
    <nav>
      <h6 class="footer-title">{{ $t('footer.resources') }}</h6>
      <router-link 
        v-for="page in resourcePages" 
        :key="page.id" 
        :to="`/pages/${page.slug}`" 
        class="link link-hover"
      >
        {{ page.title }}
      </router-link>
      <a v-if="resourcePages.length === 0" class="link link-hover opacity-50">Loading...</a>
    </nav>
  </footer>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from '../utils/axios'
import { useSiteConfigStore } from '../stores/siteConfig'

const siteConfigStore = useSiteConfigStore()

const pages = ref([])

const companyPages = computed(() => {
  return pages.value.filter(page => page.category === 'company')
})

const resourcePages = computed(() => {
  return pages.value.filter(page => page.category === 'resources')
})

const fetchPages = async () => {
  try {
    const response = await axios.get('/api/public/pages')
    pages.value = response.data.pages || []
  } catch (error) {
    console.error('Failed to fetch pages:', error)
    // Silently fail - footer will show loading state
  }
}

onMounted(() => {
  fetchPages()
})
</script>

