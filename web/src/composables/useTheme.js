import { ref, computed } from 'vue'

const STORAGE_KEY = 'theme'
const DARK = 'yunamo'
const LIGHT = 'yunamo-light'

const readTheme = () =>
  document.documentElement.getAttribute('data-theme') === LIGHT ? 'light' : 'dark'

const theme = ref(readTheme())

const applyTheme = (value) => {
  document.documentElement.setAttribute('data-theme', value === 'light' ? LIGHT : DARK)
  const meta = document.querySelector('meta[name="theme-color"]')
  if (meta) meta.setAttribute('content', value === 'light' ? '#f6f6f4' : '#0a0b0d')
  try {
    localStorage.setItem(STORAGE_KEY, value)
  } catch (e) {
    // Storage can be unavailable in private mode; the theme still applies for this visit.
  }
}

export function useTheme() {
  const isDark = computed(() => theme.value === 'dark')

  const setTheme = (value) => {
    theme.value = value
    applyTheme(value)
  }

  const toggleTheme = () => setTheme(isDark.value ? 'light' : 'dark')

  return { theme, isDark, setTheme, toggleTheme }
}
