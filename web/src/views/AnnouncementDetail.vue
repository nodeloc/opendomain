<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-4xl">
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="error" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 text-red-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="card-title text-2xl mb-2">Error</h3>
        <p>{{ error }}</p>
        <router-link to="/announcements" class="btn btn-primary mt-4">
          Back to Announcements
        </router-link>
      </div>
    </div>

    <div v-else-if="announcement" class="space-y-6">
      <!-- Back Button -->
      <div>
        <router-link to="/announcements" class="btn btn-ghost btn-sm gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back to Announcements
        </router-link>
      </div>

      <!-- Announcement Content -->
      <div class="card bg-base-200 shadow-xl">
        <div class="card-body">
          <!-- Badges -->
          <div class="flex items-center gap-2 mb-4">
            <span class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1" :class="getTypeBadgeClass(announcement.type)">
              <component :is="getTypeIcon(announcement.type)" class="h-3 w-3" />
              {{ formatType(announcement.type) }}
            </span>
            <span v-if="announcement.priority > 0" class="px-2 py-0.5 rounded text-xs font-medium inline-flex items-center gap-1 bg-amber-100 text-amber-800">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
              Priority
            </span>
          </div>

          <!-- Title -->
          <h1 class="text-4xl font-bold mb-4">{{ announcement.title }}</h1>

          <!-- Metadata -->
          <div class="flex items-center gap-4 text-sm opacity-70 mb-6 pb-6 border-b border-base-300">
            <div class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              {{ formatDate(announcement.published_at || announcement.created_at) }}
            </div>
            <div v-if="announcement.author_name" class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              {{ announcement.author_name }}
            </div>
            <div class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              {{ announcement.views || 0 }} views
            </div>
          </div>

          <!-- Content -->
          <div class="prose prose-lg max-w-none" v-html="formatContent(announcement.content)"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../utils/axios'

const route = useRoute()
const router = useRouter()

const announcement = ref(null)
const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  await fetchAnnouncement()
})

const fetchAnnouncement = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await axios.get(`/api/public/announcements/${route.params.id}`)
    announcement.value = response.data
  } catch (err) {
    console.error('Failed to fetch announcement:', err)
    error.value = err.response?.data?.error || 'Failed to load announcement'
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
  const types = {
    general: 'General',
    maintenance: 'Maintenance',
    update: 'Update',
    important: 'Important',
  }
  return types[type] || type
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
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
