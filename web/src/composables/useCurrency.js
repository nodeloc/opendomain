import { computed } from 'vue'
import { useSiteConfigStore } from '@/stores/siteConfig'

export function useCurrency() {
  const siteConfig = useSiteConfigStore()

  const currencySymbol = computed(() => siteConfig.currencySymbol || 'NL')

  const formatPrice = (amount) => {
    if (amount === null || amount === undefined) return '-'
    const formatted = Number(amount).toFixed(2)
    return `${currencySymbol.value}${formatted}`
  }

  return {
    currencySymbol,
    formatPrice
  }
}
