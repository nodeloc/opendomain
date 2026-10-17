<template>
  <component :is="layout">
    <router-view />
  </component>
  <Toast />
</template>

<script setup>
import { computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import SiteLayout from './layouts/SiteLayout.vue'
import ConsoleLayout from './layouts/ConsoleLayout.vue'
import Toast from './components/Toast.vue'
import { useSiteConfigStore } from './stores/siteConfig'

const route = useRoute()
const siteConfigStore = useSiteConfigStore()
siteConfigStore.fetch()

const layout = computed(() => (route.meta.layout === 'console' ? ConsoleLayout : SiteLayout))

watch(
  () => [siteConfigStore.siteName, siteConfigStore.siteDescription],
  ([name, description]) => {
    if (name) {
      document.title = description ? `${name} · ${description}` : name
    }
  }
)
</script>
