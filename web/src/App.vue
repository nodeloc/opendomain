<template>
  <div id="app" class="min-h-screen flex flex-col">
    <Navbar />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer />
    <Toast />
  </div>
</template>

<script setup>
import { watch } from 'vue'
import Navbar from './components/Navbar.vue'
import Footer from './components/Footer.vue'
import Toast from './components/Toast.vue'
import { useSiteConfigStore } from './stores/siteConfig'

const siteConfigStore = useSiteConfigStore()
siteConfigStore.fetch()

watch(() => siteConfigStore.siteName, (name) => {
  if (name) {
    document.title = `${name} - ${siteConfigStore.siteDescription}`
  }
})
</script>
