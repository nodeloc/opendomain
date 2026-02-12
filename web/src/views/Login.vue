<template>
  <div class="min-h-[70vh] flex items-center justify-center px-4 sm:px-6 py-8">
    <div class="card w-full max-w-lg shadow-2xl bg-base-200">
      <div class="card-body">
        <h2 class="card-title text-3xl font-bold justify-center mb-6">{{ $t('login.title', { siteName: siteConfigStore.siteName }) }}</h2>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('login.email') }}</span>
            </label>
            <input
              v-model="form.email"
              type="email"
              :placeholder="$t('login.emailPlaceholder')"
              class="input input-bordered"
              required
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">{{ $t('login.password') }}</span>
            </label>
            <input
              v-model="form.password"
              type="password"
              :placeholder="$t('login.passwordPlaceholder')"
              class="input input-bordered"
              required
            />
            <label class="label">
              <a href="#" class="label-text-alt link link-hover">{{ $t('login.forgotPassword') }}</a>
            </label>
          </div>

          <div v-if="error" class="alert alert-error">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ error }}</span>
          </div>

          <div class="form-control mt-6">
            <button type="submit" class="btn btn-primary" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              <span v-else>{{ $t('login.loginButton') }}</span>
            </button>
          </div>
        </form>

        <div v-if="hasOAuth" class="divider">{{ $t('common.or') }}</div>

        <div v-if="hasOAuth" class="flex flex-col gap-3">
          <a v-if="siteConfig.oauth?.github" :href="`${apiBase}/api/auth/github`" class="btn btn-outline gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            {{ $t('login.githubLogin') }}
          </a>
          <a v-if="siteConfig.oauth?.google" :href="`${apiBase}/api/auth/google`" class="btn btn-outline gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            {{ $t('login.googleLogin') }}
          </a>
          <a v-if="siteConfig.oauth?.nodeloc" :href="`${apiBase}/api/auth/nodeloc`" class="btn btn-outline gap-2">
            <svg viewBox="0 0 102 64" version="1.1" xmlns="http://www.w3.org/2000/svg" width="20" height="13">
              <g>
                <path d="M51.504 19.32C51.504 23.224 51.312 28.888 50.928 36.312C50.544 43.608 50.352 49.144 50.352 52.92C50.352 54.008 50.416 55.64 50.544 57.816C50.672 59.992 50.736 61.688 50.736 62.904C50.608 63.288 49.36 63.64 46.992 63.96C44.624 64.28 42.992 64.44 42.096 64.44C41.648 64.44 40.944 63.832 39.984 62.616C39.024 61.4 37.808 59.704 36.336 57.528C33.904 54.072 32.368 52.088 31.728 51.576C31.728 51.576 31.216 50.808 30.192 49.272C29.104 47.544 27.76 45.496 26.16 43.128C24.624 40.76 23.856 39.544 23.856 39.48L15.504 29.592C15.44 29.464 15.28 29.112 15.024 28.536C14.768 27.96 14.416 27.672 13.968 27.672C13.776 27.672 13.584 27.864 13.392 28.248C13.264 28.568 13.2 28.92 13.2 29.304C13.2 36.024 13.264 41.304 13.392 45.144C13.456 46.04 13.488 47.16 13.488 48.504C13.488 48.696 13.456 48.92 13.392 49.176C13.328 49.368 13.264 49.528 13.2 49.656V64.44L12.432 65.208L10.32 65.304L0.528 64.536C0.528 62.744 0.592 60.632 0.72 58.2C0.912 55.768 1.008 54.328 1.008 53.88C1.008 52.728 0.944 51.128 0.816 49.08C0.688 46.776 0.624001 44.792 0.624001 43.128L0.816 39.864C1.008 36.216 1.104 34.2 1.104 33.816C1.104 32.28 1.04 29.976 0.912 26.904C0.784 23.832 0.72 21.56 0.72 20.088C0.72 19.192 0.848 17.912 1.104 16.248C1.36 14.712 1.488 13.592 1.488 12.888C1.488 12.696 1.36 11.736 1.104 10.008C0.848 8.216 0.72 6.68 0.72 5.4L0.336 4.72799C0.464 4.088 0.688 3.672 1.008 3.48C1.456 3.416 2.672 3.224 4.656 2.904C6.64 2.52 8.24 2.328 9.456 2.328C10.544 2.328 11.28 2.456 11.664 2.712C12.112 3.096 12.816 3.768 13.776 4.72799C14.736 5.688 15.184 6.264 15.12 6.456L20.016 14.808C20.016 14.872 20.912 16.152 22.704 18.648L24.24 20.856C27.44 25.592 30.256 29.656 32.688 33.048C35.12 36.376 36.688 38.04 37.392 38.04C37.84 38.04 38.096 37.88 38.16 37.56C38.16 36.536 38.096 35.128 37.968 33.336L37.872 29.592C37.872 28.44 37.936 26.648 38.064 24.216L38.16 18.84C38.16 17.368 38.096 15.16 37.968 12.216C37.84 9.4 37.776 7.224 37.776 5.688C37.776 4.152 37.872 3.064 38.064 2.424C38.32 1.784 38.8 1.336 39.504 1.08C40.208 0.759996 41.36 0.599997 42.96 0.599997C46.16 0.599997 48.144 0.759996 48.912 1.08C49.68 1.336 50.096 1.816 50.16 2.51999C50.224 3.16 50.416 3.704 50.736 4.152V10.2L51.024 11.064L50.928 14.328C50.928 16.76 51.12 18.424 51.504 19.32ZM100.475 49.368C100.795 49.368 101.083 49.464 101.339 49.656C101.659 49.784 101.819 50.008 101.819 50.328V60.6C101.819 61.24 101.595 61.784 101.147 62.232C100.699 62.616 100.251 62.808 99.8025 62.808L94.7145 62.616C92.0905 62.616 90.2025 62.712 89.0505 62.904L86.7465 63C86.7465 62.872 85.8185 62.744 83.9625 62.616L78.9705 62.232C77.6265 62.232 76.2825 62.424 74.9385 62.808C73.5945 63.192 72.9225 63.416 72.9225 63.48C71.8345 63.48 70.5865 63.32 69.1785 63C67.5145 62.744 66.2985 62.616 65.5305 62.616C65.5305 61.4 65.4025 59.512 65.1465 56.952C64.8905 54.264 64.7625 52.216 64.7625 50.808C64.7625 50.808 64.7625 50.68 64.7625 50.424C64.8265 50.168 64.8585 49.528 64.8585 48.504C64.8585 46.968 64.8265 45.528 64.7625 44.184L64.2825 43.704V27.96C63.9625 26.04 63.8025 24.344 63.8025 22.872C63.8025 21.464 63.9625 20.568 64.2825 20.184V11.256C64.0905 11.128 63.9945 10.584 63.9945 9.62399V8.184C63.9945 6.2 63.8345 4.536 63.5145 3.192L63.2265 2.808C63.2265 2.168 64.2185 1.656 66.2025 1.272C68.1865 0.823997 70.9065 0.599997 74.3625 0.599997C74.8745 0.599997 75.2905 0.759996 75.6105 1.08C75.9305 1.336 76.0905 1.656 76.0905 2.04C76.0905 2.744 75.9945 3.576 75.8025 4.536C75.3545 6.904 75.1305 9.144 75.1305 11.256C75.1305 11.32 75.0665 12.44 74.9385 14.616C74.8105 16.728 74.6185 18.584 74.3625 20.184L74.9385 20.664V32.184C75.3225 32.888 75.6425 34.488 75.8985 36.984C76.2185 39.416 76.4105 40.92 76.4745 41.496C76.7945 44.568 77.0825 46.936 77.3385 48.6C77.5945 50.2 77.9145 51.192 78.2985 51.576C78.5545 51.704 78.8745 51.832 79.2585 51.96C79.6425 52.088 79.9625 52.184 80.2185 52.248L83.6745 52.152C86.3625 51.96 90.0745 51.864 94.8105 51.864C95.0665 51.8 95.3865 51.736 95.7705 51.672C96.2185 51.608 96.5705 51.576 96.8265 51.576C98.0425 51.512 98.9385 51.352 99.5145 51.096C100.091 50.84 100.411 50.264 100.475 49.368Z" fill="url(#paint0_login)"/>
                <defs>
                  <linearGradient id="paint0_login" x1="-6" y1="35" x2="104" y2="35" gradientUnits="userSpaceOnUse">
                    <stop stop-color="#009966"/>
                    <stop offset="0.790909" stop-color="#FF9933"/>
                  </linearGradient>
                </defs>
              </g>
            </svg>
            {{ $t('login.nodelocLogin') }}
          </a>
        </div>

        <div class="text-center mt-4">
          <p class="text-sm">
            {{ $t('login.noAccount') }}
            <router-link to="/register" class="link link-primary font-semibold">{{ $t('login.signUp') }}</router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSiteConfigStore } from '../stores/siteConfig'

const router = useRouter()
const authStore = useAuthStore()
const siteConfigStore = useSiteConfigStore()
const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000'

const siteConfig = computed(() => ({
  oauth: siteConfigStore.oauth,
}))

const hasOAuth = computed(() => {
  const o = siteConfigStore.oauth
  return o?.github || o?.google || o?.nodeloc
})

const form = ref({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  const result = await authStore.login(form.value)

  if (result.success) {
    await authStore.fetchProfile()
    router.push('/dashboard')
  } else {
    error.value = result.message
  }

  loading.value = false
}
</script>
