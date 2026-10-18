<template>
  <div ref="scene" class="energy-art" @pointermove="moveLight" @pointerleave="resetLight">
    <canvas ref="canvas" v-show="ready" aria-hidden="true"></canvas>
    <img
      ref="artwork"
      v-show="!ready"
      class="crystal-fallback"
      src="/brand/energy-cube-mint-v2.png"
      alt=""
      aria-hidden="true"
      width="1254"
      height="1254"
      loading="lazy"
      decoding="async"
      @load="onArtworkReady"
    />
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps({ active: { type: Boolean, default: true } })
const scene = ref(null)
const canvas = ref(null)
const artwork = ref(null)
const ready = ref(false)
let ctx
let width = 0
let height = 0
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

function glow(x, y, radius, color) {
  const gradient = ctx.createRadialGradient(x, y, 0, x, y, radius)
  gradient.addColorStop(0, color)
  gradient.addColorStop(1, 'rgba(62,230,180,0)')
  ctx.fillStyle = gradient
  ctx.fillRect(x - radius, y - radius, radius * 2, radius * 2)
}

function draw(time) {
  if (!ctx || !width || !height || !artwork.value?.complete || !artwork.value.naturalWidth) return
  ctx.clearRect(0, 0, width, height)
  const unit = Math.min(width, height)
  const cx = width * .5
  const cy = height * .435 + Math.sin(time * .8) * unit * .012
  const pulse = .82 + Math.sin(time * .65) * .18
  const size = unit * .83

  // Keep the surrounding light soft, with no wireframe or orbital decoration.
  glow(cx, cy, unit * .44, `rgba(33,159,118,${.14 * pulse})`)
  ctx.save()
  ctx.translate(cx, height * .84)
  ctx.scale(1, .14)
  glow(0, 0, unit * .29, `rgba(62,230,180,${.15 * pulse})`)
  ctx.restore()

  ctx.save()
  ctx.translate(cx + lightX * unit * .009, cy + lightY * unit * .007)
  ctx.rotate(Math.sin(time * .26) * .018 + lightX * .025)
  ctx.transform(1, lightY * .012, lightX * .018, 1, 0, 0)
  ctx.drawImage(artwork.value, -size / 2, -size / 2, size, size)

  // A slow, broad light pass breathes across the existing glass reflections.
  ctx.globalCompositeOperation = 'screen'
  glow(Math.sin(time * .35) * unit * .08, Math.cos(time * .24) * unit * .06,
    unit * .16, `rgba(176,255,225,${.045 * pulse})`)
  ctx.restore()

  for (let i = 0; i < 7; i++) {
    const phase = (i * .137 + time * .018) % 1
    const x = cx + Math.sin(i * 17.1) * unit * .32
    const y = height * .8 - phase * unit * .61
    ctx.beginPath()
    ctx.arc(x, y, i % 2 ? .65 : 1, 0, Math.PI * 2)
    ctx.fillStyle = `rgba(149,246,210,${Math.sin(phase * Math.PI) * .27})`
    ctx.fill()
  }
}

function onArtworkReady() {
  if (!ctx || !artwork.value?.naturalWidth) return
  ready.value = true
  resize()
  syncAnimation()
}

const canAnimate = () => props.active && visible && !document.hidden && ready.value
function tick(timestamp) {
  frame = 0
  if (!canAnimate()) return
  if (lastTime) elapsed += Math.min((timestamp - lastTime) / 1000, .1)
  lastTime = timestamp
  if (timestamp - lastDraw >= 1000 / 30) {
    lightX += (targetX - lightX) * .07
    lightY += (targetY - lightY) * .07
    draw(elapsed)
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
  if (!ctx || !scene.value) return
  const rect = scene.value.getBoundingClientRect()
  width = rect.width
  height = rect.height
  const dpr = Math.min(window.devicePixelRatio || 1, 1.5)
  canvas.value.width = Math.round(width * dpr)
  canvas.value.height = Math.round(height * dpr)
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  draw(elapsed)
}
function moveLight(event) {
  if (!canAnimate() || event.pointerType !== 'mouse') return
  const rect = scene.value.getBoundingClientRect()
  targetX = (event.clientX - rect.left) / rect.width * 2 - 1
  targetY = (event.clientY - rect.top) / rect.height * 2 - 1
}
function resetLight() { targetX = 0; targetY = 0 }

watch(() => props.active, syncAnimation)
onMounted(() => {
  try { ctx = canvas.value?.getContext('2d', { alpha: true }) } catch { return }
  if (!ctx) return
  resize()
  if (artwork.value?.complete && artwork.value.naturalWidth) onArtworkReady()
  if ('ResizeObserver' in window) {
    resizeObserver = new ResizeObserver(resize)
    resizeObserver.observe(scene.value)
  } else window.addEventListener('resize', resize)
  if ('IntersectionObserver' in window) {
    visibilityObserver = new IntersectionObserver(([entry]) => {
      visible = entry.isIntersecting && entry.intersectionRatio >= .01
      syncAnimation()
    }, { threshold: .01 })
    visibilityObserver.observe(scene.value)
  } else visible = true
  document.addEventListener('visibilitychange', syncAnimation)
  syncAnimation()
})
onUnmounted(() => {
  window.cancelAnimationFrame(frame)
  resizeObserver?.disconnect()
  visibilityObserver?.disconnect()
  document.removeEventListener('visibilitychange', syncAnimation)
  window.removeEventListener('resize', resize)
})
</script>

<style scoped>
.energy-art { position: relative; width: 100%; aspect-ratio: 1 / 1; min-width: 0; }
.energy-art canvas { position: absolute; inset: 0; width: 100%; height: 100%; }
.crystal-fallback { position: absolute; top: 2%; left: 8.5%; width: 83%; height: 83%; object-fit: contain; }
</style>
