import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

// One owner for the homepage's timers, observers, and motion preference.
export function useHomeMotion() {
  const homePage = ref(null)
  const hero = ref(null)
  const motionPaused = ref(false)
  const reducedMotion = ref(false)
  const heroVisible = ref(true)
  const pageVisible = ref(true)
  const exampleNames = ['hello-world', 'my-app', 'portfolio']
  const exampleIndex = ref(0)
  const exampleName = computed(() => exampleNames[exampleIndex.value])
  const motionEnabled = computed(() => !motionPaused.value && !reducedMotion.value)
  const motionActive = computed(() => motionEnabled.value && heroVisible.value && pageVisible.value)
  let mediaQuery
  let heroObserver
  let revealObserver
  let exampleTimer
  let pointerFrame

  const clearExampleTimer = () => {
    window.clearInterval(exampleTimer)
    exampleTimer = undefined
  }
  const syncExamples = () => {
    clearExampleTimer()
    if (motionActive.value) {
      exampleTimer = window.setInterval(() => {
        exampleIndex.value = (exampleIndex.value + 1) % exampleNames.length
      }, 3600)
    }
  }
  const onVisibilityChange = () => { pageVisible.value = !document.hidden }
  const onMotionChange = (event) => { reducedMotion.value = event.matches }
  const toggleMotion = () => {
    motionPaused.value = !motionPaused.value
    try { localStorage.setItem('yunamo-home-motion', motionPaused.value ? 'paused' : 'on') } catch {}
  }

  const moveSpotlight = (event) => {
    if (!motionEnabled.value || event.pointerType === 'touch') return
    const element = event.currentTarget
    const { clientX, clientY } = event
    window.cancelAnimationFrame(pointerFrame)
    pointerFrame = window.requestAnimationFrame(() => {
      const rect = element.getBoundingClientRect()
      element.style.setProperty('--spot-x', `${clientX - rect.left}px`)
      element.style.setProperty('--spot-y', `${clientY - rect.top}px`)
    })
  }

  watch(motionActive, syncExamples)
  onMounted(() => {
    try { motionPaused.value = localStorage.getItem('yunamo-home-motion') === 'paused' } catch {}
    mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
    reducedMotion.value = mediaQuery.matches
    mediaQuery.addEventListener('change', onMotionChange)
    onVisibilityChange()
    document.addEventListener('visibilitychange', onVisibilityChange)

    if ('IntersectionObserver' in window) {
      heroObserver = new IntersectionObserver(([entry]) => {
        heroVisible.value = entry.isIntersecting && entry.intersectionRatio >= 0.01
      }, { threshold: 0.01, rootMargin: '-64px 0px 0px 0px' })
      if (hero.value) heroObserver.observe(hero.value)

      revealObserver = new IntersectionObserver((entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            entry.target.classList.add('is-revealed')
            revealObserver.unobserve(entry.target)
          }
        }
      }, { threshold: 0.08, rootMargin: '0px 0px -24px 0px' })
      homePage.value?.querySelectorAll('[data-reveal]').forEach((element) => revealObserver.observe(element))
      homePage.value?.classList.add('motion-observed')
    }
    syncExamples()
  })

  onUnmounted(() => {
    clearExampleTimer()
    window.cancelAnimationFrame(pointerFrame)
    heroObserver?.disconnect()
    revealObserver?.disconnect()
    mediaQuery?.removeEventListener('change', onMotionChange)
    document.removeEventListener('visibilitychange', onVisibilityChange)
  })

  return { homePage, hero, motionPaused, reducedMotion, motionEnabled, motionActive,
    exampleNames, exampleName, toggleMotion, moveSpotlight }
}
