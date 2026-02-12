import { defineStore } from 'pinia'
import axios from '../utils/axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user,
    isAdmin: (state) => state.user?.is_admin === true,
  },

  actions: {
    async register(userData) {
      try {
        const response = await axios.post('/api/auth/register', userData)
        return { success: true, message: response.data.message }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.error || 'Registration failed',
        }
      }
    },

    async login(credentials) {
      try {
        const response = await axios.post('/api/auth/login', credentials)
        this.token = response.data.token
        this.user = response.data.user
        localStorage.setItem('token', this.token)
        return { success: true }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.error || 'Login failed',
        }
      }
    },

    async fetchProfile() {
      try {
        const response = await axios.get('/api/user/profile')
        this.user = response.data
      } catch (error) {
        console.error('Failed to fetch profile:', error)
        // If we get 401, the token is invalid, logout
        if (error.response?.status === 401) {
          this.logout()
        }
      }
    },

    // Check if token is valid and has admin privileges
    async validateAdminToken() {
      if (!this.token) return false

      try {
        if (!this.user) {
          await this.fetchProfile()
        }
        return this.isAdmin
      } catch (error) {
        return false
      }
    },

    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    },
  },
})
