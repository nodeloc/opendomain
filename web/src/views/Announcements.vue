<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <div>
      <h1 class="text-3xl font-bold">{{ $t('announcement.title') }}</h1>
      <p class="text-lg opacity-70 mt-2">{{ $t('announcement.subtitle') }}</p>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="announcements.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 opacity-40 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
        </svg>
        <h3 class="card-title text-2xl mb-2">{{ $t('announcement.noAnnouncements') }}</h3>
        <p>{{ $t('announcement.checkLater') }}</p>
      </div>
    </div>

    <div v-else class="grid gap-4">
      <router-link
        v-for="announcement in announcements"
        :key="announcement.id"
        :to="`/announcements/${announcement.id}`"
        class="card bg-base-200 shadow-xl hover:shadow-2xl transition-all"
      >
        <div class="card-body">
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <span class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1" :class="getTypeBadgeClass(announcement.type)">
                  <component :is="getTypeIcon(announcement.type)" class="h-3 w-3" />
                  {{ formatType(announcement.type) }}
                </span>
                <span v-if="announcement.priority > 0" class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1 bg-amber-100 text-amber-800">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                  </svg>
                  {{ $t('announcement.priority') }}
                </span>
              </div>
              <h2 class="card-title text-2xl">{{ announcement.title }}</h2>
              <p class="text-sm opacity-70 mt-2">
                {{ formatDate(announcement.published_at || announcement.created_at) }}
                <span v-if="announcement.author_name"> · {{ $t('announcement.by') }} {{ announcement.author_name }}</span>
                <span> · {{ announcement.views }} {{ $t('announcement.views') }}</span>
              </p>
            </div>
            <component :is="getTypeIcon(announcement.type)" class="h-12 w-12 flex-shrink-0" />
          </div>
          <div class="mt-3 prose prose-sm max-w-none line-clamp-3" v-html="formatContent(announcement.content)"></div>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from '../utils/axios'

const { t } = useI18n()
const announcements = ref([])
const loading = ref(true)

onMounted(async () => {
  await fetchAnnouncements()
})

const fetchAnnouncements = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/public/announcements')
    announcements.value = response.data.announcements
  } catch (error) {
    console.error('Failed to fetch announcements:', error)
  } finally {
    loading.value = false
  }
}

const getTypeBadgeClass = (type) => {
  const classes = {
    general: 'bg-blue-100 text-blue-800',
    maintenance: 'bg-amber-100 text-amber-800',
    update: 'bg-green-100 text-green-800',
    important: 'bg-red-100 text-red-800',
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getTypeIcon = (type) => {
  const icons = {
    general: 'svg',
    maintenance: 'svg',
    update: 'svg',
    important: 'svg',
  }
  // General - Megaphone
  if (type === 'general') {
    return { 
      template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" /></svg>'
    }
  }
  // Maintenance - Wrench
  if (type === 'maintenance') {
    return {
      template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>'
    }
  }
  // Update - Sparkles
  if (type === 'update') {
    return {
      template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" /></svg>'
    }
  }
  // Important - Exclamation
  if (type === 'important') {
    return {
      template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>'
    }
  }
  // Default
  return {
    template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" /></svg>'
  }
}

const formatType = (type) => {
  const typeKey = `announcement.${type}`
  return t(typeKey)
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const formatContent = (content) => {
  // Basic markdown-like formatting
  return content
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
}
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
