<template>
  <div class="y-page">
    <YLoader :label="$t('uiUser.common.loading')" />
  </div>
</template>

<script setup>
// DNS records are managed inline on /domains now; this route only forwards
// bookmarked /domains/:domainId/dns links to /domains?d=<id>.
import { useRoute, useRouter } from 'vue-router'
import YLoader from '../components/brand/YLoader.vue'

const route = useRoute()
const router = useRouter()

router.replace({ path: '/domains', query: { d: route.params.domainId } })
</script>
