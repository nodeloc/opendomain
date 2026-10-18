// Procedural projected glass, filaments and light: no image textures.
export function createEnergyCubeRenderer(canvas) {
  const ctx = canvas.getContext('2d', { alpha: true })
  if (!ctx) return null
  let width = 0
  let height = 0
  let lightX = 0
  let lightY = 0
  // Two hexagonal shoulders and pointed ends form 24 individual crystal facets.
  const points = [{ x: .06, y: -1.35, z: 0 }]
  for (const [y, radius, twist] of [[-.44, .86, 0], [.48, .73, .13]]) {
    for (let i = 0; i < 6; i++) {
      const a = i * Math.PI / 3 + twist
      points.push({ x: Math.cos(a) * radius, y, z: Math.sin(a) * radius })
    }
  }
  points.push({ x: -.04, y: 1.4, z: .02 })
  const faces = []
  for (let i = 0; i < 6; i++) {
    const upper = i + 1
    const nextUpper = (i + 1) % 6 + 1
    const lower = i + 7
    const nextLower = (i + 1) % 6 + 7
    faces.push([0, upper, nextUpper], [upper, lower, nextUpper],
      [nextUpper, lower, nextLower], [lower, 13, nextLower])
  }
  const facets = faces.map((indices, index) => {
    const [a, b, c] = indices.map(i => points[i])
    const ab = { x: b.x-a.x, y: b.y-a.y, z: b.z-a.z }
    const ac = { x: c.x-a.x, y: c.y-a.y, z: c.z-a.z }
    const normal = { x: ab.y*ac.z-ab.z*ac.y, y: ab.z*ac.x-ab.x*ac.z, z: ab.x*ac.y-ab.y*ac.x }
    const center = { x: (a.x+b.x+c.x)/3, y: (a.y+b.y+c.y)/3, z: (a.z+b.z+c.z)/3 }
    const sign = Math.sign(normal.x*center.x + normal.y*center.y + normal.z*center.z)
    const length = Math.hypot(normal.x, normal.y, normal.z)
    return { indices, index, center, normal: { x: normal.x*sign/length, y: normal.y*sign/length, z: normal.z*sign/length } }
  })
  const dust = Array.from({ length: 34 }, (_, i) => ({
    offset: ((i * 73 + 19) % 101) / 101,
    spread: Math.sin(i * 23.7),
    speed: .045 + (i % 5) * .012,
    radius: .6 + (i % 3) * .45,
  }))

  function rotate(point, angle, tilt) {
    const x = point.x * Math.cos(angle) + point.z * Math.sin(angle)
    const z = -point.x * Math.sin(angle) + point.z * Math.cos(angle)
    const y = point.y * Math.cos(tilt) - z * Math.sin(tilt)
    return { x, y, z: point.y * Math.sin(tilt) + z * Math.cos(tilt) }
  }
  function project(point, scale, cx, cy) {
    const perspective = 5.8 / (5.8 - point.z)
    return { x: cx + point.x * scale * perspective, y: cy + point.y * scale * perspective }
  }
  function trace(projected, indices) {
    ctx.beginPath()
    indices.forEach((index, i) => {
      const point = projected[index]
      if (i === 0) ctx.moveTo(point.x, point.y)
      else ctx.lineTo(point.x, point.y)
    })
    ctx.closePath()
  }
  function glow(x, y, radius, color) {
    const gradient = ctx.createRadialGradient(x, y, 0, x, y, radius)
    gradient.addColorStop(0, color)
    gradient.addColorStop(1, 'rgba(62,230,180,0)')
    ctx.fillStyle = gradient
    ctx.fillRect(x - radius, y - radius, radius * 2, radius * 2)
  }
  function drawRing(cx, cy, rx, ry, angle, opacity) {
    ctx.save()
    ctx.translate(cx, cy)
    ctx.rotate(angle)
    ctx.strokeStyle = `rgba(101,235,199,${opacity})`
    ctx.lineWidth = .8
    ctx.beginPath()
    ctx.ellipse(0, 0, rx, ry, 0, 0, Math.PI * 2)
    ctx.stroke()
    ctx.restore()
  }

  function render(time) {
    if (!ctx || !width || !height) return
    ctx.clearRect(0, 0, width, height)
    const unit = Math.min(width, height)
    const scale = unit * .255
    const cx = width * .5
    const cy = height * .43 + Math.sin(time * .85) * unit * .014
    const baseY = height * .84
    const pulse = .88 + Math.sin(time * 1.1) * .12
    const angle = time * .11 + .65 + lightX * .22
    const tilt = -.12 + Math.sin(time * .17) * .06 + lightY * .1

    glow(cx, cy, unit * .46, `rgba(18,158,116,${.22 * pulse})`)
    glow(cx, baseY, unit * .27, 'rgba(62,230,180,.12)')
    drawRing(cx, baseY, unit * .32, unit * .052, 0, .12)
    const beam = ctx.createLinearGradient(cx, cy + scale * .9, cx, baseY)
    beam.addColorStop(0, 'rgba(150,255,219,0)')
    beam.addColorStop(.6, 'rgba(150,255,219,.11)')
    beam.addColorStop(1, 'rgba(150,255,219,0)')
    ctx.fillStyle = beam
    ctx.fillRect(cx - 1, cy + scale * .9, 2, baseY - cy - scale * .9)

    // The rear orbit is drawn first, with its near side in front of the crystal.
    const orbit = []
    for (let i = 0; i <= 80; i++) {
      const a = i / 80 * Math.PI * 2
      orbit.push({ x: cx + Math.cos(a) * unit * .32, y: cy + Math.sin(a) * unit * .115 - Math.cos(a) * unit * .12 })
    }
    const orbitLine = (start, end, alpha) => {
      ctx.beginPath()
      for (let i = start; i <= end; i++) {
        if (i === start) ctx.moveTo(orbit[i].x, orbit[i].y)
        else ctx.lineTo(orbit[i].x, orbit[i].y)
      }
      ctx.lineWidth = .7
      ctx.strokeStyle = `rgba(106,241,203,${alpha})`
      ctx.stroke()
    }
    orbitLine(40, 80, .19)

    dust.slice(0, width < 400 ? 20 : 34).forEach((particle, index) => {
      const progress = (particle.offset + time * particle.speed) % 1
      const x = cx + particle.spread * unit * (.25 - progress * .07) + Math.sin(time * .4 + index) * unit * .015
      const y = baseY - progress * unit * .7
      const alpha = Math.sin(progress * Math.PI) * .5
      ctx.beginPath()
      ctx.arc(x, y, particle.radius, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(166,255,226,${alpha})`
      ctx.fill()
    })

    const rotated = points.map(point => rotate(point, angle, tilt))
    const projected = rotated.map(point => project(point, scale, cx, cy))
    const orderedFaces = facets.map(face => {
      const center = rotate(face.center, angle, tilt)
      const normal = rotate(face.normal, angle, tilt)
      return { ...face, normal, z: center.z,
        front: normal.x * -center.x + normal.y * -center.y + normal.z * (5.8-center.z) > 0 }
    }).sort((a, b) => a.z - b.z)

    const drawFace = (face, front) => {
      const shine = Math.max(0, face.normal.x * -.5 + face.normal.y * -.65 + face.normal.z * .57)
      trace(projected, face.indices)
      const gradient = ctx.createLinearGradient(cx - scale, cy - scale * 1.4, cx + scale, cy + scale * 1.4)
      const strength = front ? .08 + shine * .24 : .025
      gradient.addColorStop(0, `rgba(198,255,235,${strength + (front ? .08 : 0)})`)
      gradient.addColorStop(.48, `rgba(${face.index % 3 === 0 ? '52,218,167' : '20,130,96'},${strength})`)
      gradient.addColorStop(1, `rgba(8,79,58,${strength * .9})`)
      ctx.fillStyle = gradient
      ctx.fill()
      ctx.strokeStyle = front ? `rgba(195,255,233,${.15 + shine * .5})` : 'rgba(98,230,188,.08)'
      ctx.lineWidth = front ? .9 : .5
      ctx.save()
      if (front) { ctx.shadowColor = '#3ee6b4'; ctx.shadowBlur = 8 }
      ctx.stroke()
      ctx.restore()
      if (front) {
        ctx.save()
        const center = face.indices.reduce((p, i) => ({ x: p.x + projected[i].x / face.indices.length, y: p.y + projected[i].y / face.indices.length }), {x:0,y:0})
        const inset = projected.map(p => ({x:p.x*.982+center.x*.018, y:p.y*.982+center.y*.018}))
        trace(inset, face.indices)
        ctx.strokeStyle = `rgba(191,255,232,${.07 + shine*.1})`
        ctx.lineWidth = .6
        ctx.stroke()
        ctx.restore()
      }
      // A narrow reflection gives one facet a moving glass edge.
      if (front && shine > .75) {
        ctx.save()
        trace(projected, face.indices)
        ctx.clip()
        const sweep = Math.sin(time * .35 + face.index * .12) * scale
        const reflection = ctx.createLinearGradient(cx + sweep - scale * .4, cy - scale, cx + sweep + scale * .4, cy + scale)
        reflection.addColorStop(0, 'rgba(223,255,244,0)')
        reflection.addColorStop(.48, `rgba(223,255,244,${(shine - .75) * .6})`)
        reflection.addColorStop(.55, 'rgba(223,255,244,0)')
        ctx.fillStyle = reflection
        ctx.fillRect(cx - scale, cy - scale * 1.7, scale * 2, scale * 3.4)
        ctx.restore()
      }
    }
    orderedFaces.filter(face => !face.front).forEach(face => drawFace(face, false))

    // Clip the drifting energy to the crystal's projected silhouette.
    const sorted = projected.map((point, index) => ({ ...point, index })).sort((a, b) => a.x - b.x || a.y - b.y)
    const cross = (a, b, c) => (b.x - a.x) * (c.y - a.y) - (b.y - a.y) * (c.x - a.x)
    const halfHull = (list) => {
      const hull = []
      for (const point of list) {
        while (hull.length >= 2 && cross(hull.at(-2), hull.at(-1), point) <= 0) hull.pop()
        hull.push(point)
      }
      return hull.slice(0, -1)
    }
    const hull = [...halfHull(sorted), ...halfHull([...sorted].reverse())].map(point => point.index)
    ctx.save()
    trace(projected, hull)
    ctx.clip()
    ctx.globalCompositeOperation = 'screen'
    glow(cx, cy, scale * 1.35, `rgba(26,191,140,${.5 * pulse})`)
    for (let i = 0; i < 8; i++) {
      const a = time * (.25 + i * .015) + i * .82
      const x = cx + Math.cos(a) * scale * .34
      const y = cy + Math.sin(a * 1.3) * scale * .42
      glow(x, y, scale * (.32 + i % 3 * .1), `rgba(${i % 2 ? '37,225,165' : '161,255,221'},.13)`)
    }
    glow(cx, cy, scale * .56, `rgba(186,255,232,${.65 * pulse})`)
    glow(cx, cy, scale * .2, `rgba(234,255,249,${.7 * pulse})`)
    for (let i = 0; i < 6; i++) {
      const a = i * Math.PI / 3 + time * .09
      const bend = Math.sin(time * .45 + i) * .38
      const endX = cx + Math.cos(a) * scale * 1.1
      const endY = cy + Math.sin(a) * scale * .9
      const controlX = cx + Math.cos(a + bend) * scale * .65
      const controlY = cy + Math.sin(a + bend) * scale * .65
      ctx.beginPath()
      ctx.moveTo(cx, cy)
      ctx.quadraticCurveTo(controlX, controlY, endX, endY)
      const filament = ctx.createLinearGradient(cx, cy, endX, endY)
      filament.addColorStop(0, 'rgba(215,255,239,.24)')
      filament.addColorStop(.65, 'rgba(119,255,209,.12)')
      filament.addColorStop(1, 'rgba(119,255,209,0)')
      ctx.strokeStyle = filament
      ctx.lineWidth = .65
      ctx.stroke()
      const t = (time * .12 + i / 6) % 1
      const q = 1 - t
      const x = q*q*cx + 2*q*t*controlX + t*t*endX
      const y = q*q*cy + 2*q*t*controlY + t*t*endY
      glow(x, y, 4, `rgba(190,255,231,${Math.sin(t*Math.PI)*.5})`)
    }
    ctx.restore()
    orderedFaces.filter(face => face.front).forEach(face => drawFace(face, true))
    orbitLine(0, 40, .27)

    for (let i = 0; i < 3; i++) {
      const a = time * .34 + i * Math.PI * 2 / 3
      const x = cx + Math.cos(a) * unit * .32
      const y = cy + Math.sin(a) * unit * .115 - Math.cos(a) * unit * .12
      glow(x, y, 10, 'rgba(130,255,215,.25)')
      ctx.beginPath()
      ctx.arc(x, y, i ? 1.4 : 2, 0, Math.PI * 2)
      ctx.fillStyle = '#dbfff2'
      ctx.fill()
    }
    for (const tip of [projected[0], projected[13]]) {
      glow(tip.x, tip.y, unit * .025, 'rgba(210,255,237,.35)')
    }
  }
  return {
    resize(w, h, dpr) {
      width = w
      height = h
      canvas.width = Math.max(1, Math.round(w * dpr))
      canvas.height = Math.max(1, Math.round(h * dpr))
      ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
    },
    draw(time, x, y) {
      lightX = x
      lightY = y
      render(time)
    },
    dispose() { ctx.clearRect(0, 0, width, height) },
  }
}
