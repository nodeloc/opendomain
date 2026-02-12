<template>
  <div class="container mx-auto px-4 py-8">
    <div v-if="loading" class="flex justify-center items-center min-h-[400px]">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="error" class="alert alert-error">
      <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span>{{ error }}</span>
    </div>

    <article v-else-if="page" class="prose prose-lg max-w-4xl mx-auto">
      <h1>{{ page.title }}</h1>
      <div class="divider"></div>
      <div v-html="renderContent(page.content)"></div>
      <div class="text-sm text-gray-500 mt-8 pt-4 border-t">
        {{ $t('page.lastUpdated') }}: {{ formatDate(page.updated_at) }}
      </div>
    </article>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'

const route = useRoute()
const { t } = useI18n()
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
  // Simple HTML rendering
  // You can add markdown parsing here if needed
  return content
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

<style scoped>
.prose {
  @apply text-base-content;
}

.prose h1 {
  @apply text-4xl font-bold mb-4;
}

.prose h2 {
  @apply text-3xl font-bold mt-8 mb-4;
}

.prose h3 {
  @apply text-2xl font-bold mt-6 mb-3;
}

.prose p {
  @apply mb-4 leading-relaxed;
}

.prose a {
  @apply text-primary hover:underline;
}

.prose ul,
.prose ol {
  @apply ml-6 mb-4;
}

.prose li {
  @apply mb-2;
}

.prose code {
  @apply bg-base-200 px-2 py-1 rounded text-sm;
}

.prose pre {
  @apply bg-base-200 p-4 rounded-lg overflow-x-auto mb-4;
}
</style>
