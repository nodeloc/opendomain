<template>
  <div ref="scene" class="energy-art" @pointermove="moveLight" @pointerleave="resetLight">
    <canvas ref="canvas" v-show="ready" aria-hidden="true"></canvas>
    <svg v-if="!ready" class="crystal-fallback" viewBox="0 0 520 520" fill="none" aria-hidden="true">
      <defs>
        <linearGradient id="energy-fallback-glass" x1="140" y1="130" x2="385" y2="385" gradientUnits="userSpaceOnUse">
          <stop stop-color="#c9ffe9" stop-opacity=".7" /><stop offset=".45" stop-color="#3ee6b4" stop-opacity=".45" /><stop offset="1" stop-color="#0b8f6a" stop-opacity=".45" />
        </linearGradient>
        <radialGradient id="energy-fallback-core"><stop stop-color="#c9ffe9" stop-opacity=".8" /><stop offset="1" stop-color="#3ee6b4" stop-opacity="0" /></radialGradient>
      </defs>
      <path d="M264 66L357 188L340 319L254 448L173 314L163 191Z" fill="url(#energy-fallback-glass)" stroke="#bdffe4" stroke-opacity=".55" />
      <circle cx="260" cy="251" r="90" fill="url(#energy-fallback-core)" />
      <path d="M264 66L246 213L163 191L256 306L173 314M357 188L246 213L340 319L256 306L254 448M246 213L256 306" stroke="#bdffe4" stroke-opacity=".45" stroke-linejoin="round" />
    </svg>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { createEnergyCubeRenderer } from './energyCubeRenderer'

const props = defineProps({ active: { type: Boolean, default: true } })
const scene = ref(null)
const canvas = ref(null)
const ready = ref(false)
let renderer
let elapsed = 0
let frame = 0
let lastTime = 0
let lastDraw = 0
let visible = false
let resizeObserver
let visibilityObserver
let targetX = 0
let targetY = 0
let lightX = 0
let lightY = 0

const canAnimate = () => props.active && visible && !document.hidden && ready.value
function tick(timestamp) {
  frame = 0
  if (!canAnimate()) return
  if (lastTime) elapsed += Math.min((timestamp-lastTime)/1000,.1)
  lastTime = timestamp
  if (timestamp-lastDraw >= 1000/30) {
    lightX += (targetX-lightX)*.07
    lightY += (targetY-lightY)*.07
    renderer?.draw(elapsed,lightX,lightY)
    lastDraw = timestamp
  }
  frame = window.requestAnimationFrame(tick)
}
function syncAnimation() {
  window.cancelAnimationFrame(frame)
  frame = 0
  lastTime = 0
  if (canAnimate()) frame = window.requestAnimationFrame(tick)
}
function resize() {
  if (!renderer || !scene.value) return
  const rect = scene.value.getBoundingClientRect()
  const dpr = Math.min(window.devicePixelRatio || 1,1.5)
  renderer.resize(rect.width,rect.height,dpr)
  renderer.draw(elapsed,lightX,lightY)
}
function moveLight(event) {
  if (!canAnimate() || event.pointerType !== 'mouse') return
  const rect = scene.value.getBoundingClientRect()
  targetX = (event.clientX-rect.left)/rect.width*2-1
  targetY = (event.clientY-rect.top)/rect.height*2-1
}
function resetLight() { targetX=0; targetY=0 }
function initializeRenderer() {
  renderer?.dispose()
  renderer = null
  try { renderer = createEnergyCubeRenderer(canvas.value) } catch { /* Keep the vector fallback. */ }
  ready.value = Boolean(renderer)
  if (renderer) resize()
  syncAnimation()
}
function onContextLost() {
  ready.value = false
  syncAnimation()
}

watch(() => props.active,syncAnimation)
onMounted(() => {
  initializeRenderer()
  if ('ResizeObserver' in window) {
    resizeObserver = new ResizeObserver(resize)
    resizeObserver.observe(scene.value)
  } else window.addEventListener('resize',resize)
  if ('IntersectionObserver' in window) {
    visibilityObserver = new IntersectionObserver(([entry]) => {
      visible = entry.isIntersecting && entry.intersectionRatio >= .01
      syncAnimation()
    }, { threshold: .01 })
    visibilityObserver.observe(scene.value)
  } else visible=true
  canvas.value.addEventListener('contextlost',onContextLost)
  canvas.value.addEventListener('contextrestored',initializeRenderer)
  document.addEventListener('visibilitychange',syncAnimation)
  syncAnimation()
})
onUnmounted(() => {
  window.cancelAnimationFrame(frame)
  resizeObserver?.disconnect()
  visibilityObserver?.disconnect()
  document.removeEventListener('visibilitychange',syncAnimation)
  window.removeEventListener('resize',resize)
  canvas.value?.removeEventListener('contextlost',onContextLost)
  canvas.value?.removeEventListener('contextrestored',initializeRenderer)
  renderer?.dispose()
})
</script>

<style scoped>
.energy-art { position: relative; width: 100%; aspect-ratio: 1 / 1; min-width: 0; }
.energy-art canvas, .crystal-fallback { position: absolute; inset: 0; width: 100%; height: 100%; }
</style>
