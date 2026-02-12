<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 max-w-7xl space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">Announcement Management</h1>
        <p class="text-lg opacity-70 mt-2">Create and manage announcements</p>
      </div>
      <button @click="showCreateModal = true" class="btn btn-primary">
        + Create Announcement
      </button>
    </div>

    <!-- Announcements Table -->
    <div v-if="loading" class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <div v-else-if="announcements.length === 0" class="card bg-base-200 shadow-xl">
      <div class="card-body items-center text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 opacity-40 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
        </svg>
        <h3 class="card-title text-2xl mb-2">No Announcements</h3>
        <p class="mb-4">Create your first announcement to get started</p>
        <button @click="showCreateModal = true" class="btn btn-primary">Create Announcement</button>
      </div>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>Title</th>
            <th>Type</th>
            <th>Priority</th>
            <th>Views</th>
            <th>Status</th>
            <th>Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="announcement in announcements" :key="announcement.id">
            <td class="font-semibold max-w-md truncate">{{ announcement.title }}</td>
            <td><div class="px-2 py-0.5 rounded text-xs font-medium" :class="getTypeBadgeClass(announcement.type)">{{ formatType(announcement.type) }}</div></td>
            <td>{{ announcement.priority }}</td>
            <td>{{ announcement.views }}</td>
            <td>
              <div class="px-2 py-0.5 rounded text-xs font-medium inline-block" :class="announcement.is_published ? 'bg-green-100 text-green-800' : 'bg-amber-100 text-amber-800'">
                {{ announcement.is_published ? 'Published' : 'Draft' }}
              </div>
            </td>
            <td>{{ formatDate(announcement.created_at) }}</td>
            <td>
              <div class="flex gap-2">
                <button @click="editAnnouncement(announcement)" class="btn btn-sm btn-ghost">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button @click="confirmDelete(announcement)" class="btn btn-sm btn-ghost text-error">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <dialog class="modal" :class="{ 'modal-open': showCreateModal || showEditModal }">
      <div class="modal-box w-11/12 max-w-4xl">
        <h3 class="font-bold text-lg">
          {{ showEditModal ? 'Edit' : 'Create' }} Announcement
        </h3>

        <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
          <div class="form-control">
            <label class="label"><span class="label-text">Title</span></label>
            <input
              v-model="form.title"
              type="text"
              placeholder="Announcement title"
              class="input input-bordered"
              required
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label"><span class="label-text">Type</span></label>
              <select v-model="form.type" class="select select-bordered" required>
                <option value="general">General</option>
                <option value="maintenance">Maintenance</option>
                <option value="update">Update</option>
                <option value="important">Important</option>
              </select>
            </div>

            <div class="form-control">
              <label class="label"><span class="label-text">Priority (0-100)</span></label>
              <input
                v-model.number="form.priority"
                type="number"
                min="0"
                max="100"
                class="input input-bordered"
              />
            </div>
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Content</span></label>
            <textarea
              v-model="form.content"
              class="textarea textarea-bordered"
              rows="10"
              placeholder="Announcement content (supports **bold** and *italic*)"
              required
            ></textarea>
            <label class="label">
              <span class="label-text-alt">Tip: Use **text** for bold, *text* for italic</span>
            </label>
          </div>

          <div v-if="showEditModal" class="form-control">
            <label class="label cursor-pointer">
              <span class="label-text">Published</span>
              <input v-model="form.is_published" type="checkbox" class="checkbox" />
            </label>
          </div>

          <div class="modal-action">
            <button type="button" @click="closeModal" class="btn">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="loading loading-spinner"></span>
              <span v-else>{{ showEditModal ? 'Update' : 'Create' }}</span>
            </button>
          </div>
        </form>
      </div>
      <form method="dialog" class="modal-backdrop" @click="closeModal">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from '../utils/axios'
import { useToast } from '../composables/useToast'

const toast = useToast()
const announcements = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const submitting = ref(false)
const editingAnnouncement = ref(null)

const form = ref({
  title: '',
  content: '',
  type: 'general',
  priority: 0,
  is_published: false,
})

onMounted(async () => {
  await fetchAnnouncements()
})

const fetchAnnouncements = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/admin/announcements')
    announcements.value = response.data.announcements
  } catch (error) {
    console.error('Failed to fetch announcements:', error)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (showEditModal.value) {
      await axios.put(`/api/admin/announcements/${editingAnnouncement.value.id}`, form.value)
      toast.success('Announcement updated successfully')
    } else {
      await axios.post('/api/admin/announcements', form.value)
      toast.success('Announcement created successfully')
    }
    closeModal()
    await fetchAnnouncements()
  } catch (error) {
    toast.error(error.response?.data?.error || 'Operation failed')
  } finally {
    submitting.value = false
  }
}

const editAnnouncement = (announcement) => {
  editingAnnouncement.value = announcement
  form.value = {
    title: announcement.title,
    content: announcement.content,
    type: announcement.type,
    priority: announcement.priority,
    is_published: announcement.is_published,
  }
  showEditModal.value = true
}

const confirmDelete = async (announcement) => {
  if (confirm(`Delete announcement "${announcement.title}"?`)) {
    try {
      await axios.delete(`/api/admin/announcements/${announcement.id}`)
      toast.success('Announcement deleted successfully')
      await fetchAnnouncements()
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to delete announcement')
    }
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingAnnouncement.value = null
  form.value = {
    title: '',
    content: '',
    type: 'general',
    priority: 0,
    is_published: false,
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
    month: 'short',
    day: 'numeric',
  })
}
</script>
