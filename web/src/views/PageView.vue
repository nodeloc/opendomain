<template>
  <div class="y-container py-14 sm:py-16">
    <YLoader v-if="loading" class="min-h-[400px]" :label="$t('page.loading')" />

    <div v-else-if="error" class="mx-auto max-w-xl">
      <div class="y-panel flex flex-col items-center px-6 py-14 text-center">
        <span class="inline-flex h-10 w-10 items-center justify-center rounded-md border border-y-line bg-y-surface-2 text-y-danger-text">
          <CircleX :size="20" :stroke-width="1.5" />
        </span>
        <h1 class="y-h3 mt-5">{{ $t('uiPublic.page.errorTitle') }}</h1>
        <p class="mt-2 max-w-md text-sm leading-6 text-y-muted">{{ error }}</p>
        <router-link to="/" class="btn btn-primary mt-6">{{ $t('uiPublic.backHome') }}</router-link>
      </div>
    </div>

    <article v-else-if="page" class="mx-auto max-w-3xl">
      <header>
        <p class="y-eyebrow">{{ siteName }}</p>
        <h1 class="y-h1 mt-4 break-words">{{ page.title }}</h1>
        <p class="mt-4 text-sm text-y-faint">
          {{ $t('page.lastUpdated') }}: <span class="font-mono">{{ formatDate(page.updated_at) }}</span>
        </p>
      </header>

      <div class="y-panel mt-10 p-6 sm:p-10">
        <div class="prose max-w-none sm:prose-lg" v-html="renderContent(page.content)"></div>
      </div>
    </article>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { CircleX } from 'lucide-vue-next'
import axios from '../utils/axios'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useSiteConfigStore } from '../stores/siteConfig'
import YLoader from '../components/brand/YLoader.vue'

const route = useRoute()
const { t } = useI18n()
const siteConfigStore = useSiteConfigStore()
const siteName = computed(() => siteConfigStore.siteName || 'Yunamo')
const page = ref(null)
const loading = ref(true)
const error = ref(null)

const fetchPage = async () => {
  try {
    loading.value = true
    error.value = null
    const response = await axios.get(`/api/public/pages/${route.params.slug}`)
    page.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || t('page.loadFailed')
  } finally {
    loading.value = false
  }
}

const renderContent = (content) => {
  if (!content) return ''
  return DOMPurify.sanitize(marked.parse(content))
}

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

// Watch route params change
watch(() => route.params.slug, () => {
  fetchPage()
})

onMounted(() => {
  fetchPage()
})
</script>
