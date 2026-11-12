<template>
  <div ref="stage" class="sack3d" :class="{ 'sack3d--ready': ready, 'sack3d--grab': dragging }">
    <!-- Gambar biasa tampil dulu, dan tetap dipakai bila WebGL tidak tersedia -->
    <img v-if="!ready" class="sack3d__fallback" :src="fallback" alt="Kemasan Pupuk Organik Agrina 99 50 kg" />
    <div v-if="ready && hint" class="sack3d__hint"><v-icon size="16">mdi-gesture-swipe-horizontal</v-icon> Geser untuk memutar</div>
  </div>
</template>

<script>
// Sak Pupuk Organik Agrina 99 dalam 3D (three.js).
// three.js dimuat terpisah (lazy) supaya teks halaman tetap tampil cepat di HP.
export default {
  name: 'SackViewer',
  props: {
    fallback: { type: String, default: '/brand/pack-agrina99.jpg' },
  },
  data() {
    return { ready: false, dragging: false, hint: true }
  },
  async mounted() {
    if (!this.webglAvailable()) return
    try {
      const THREE = await import('three')
      const { drawFront, drawBack, drawPlain, drawShadow } = await import('./sackTextures.js')
      if (this._destroyed) return
      await this.build(THREE, { drawFront, drawBack, drawPlain, drawShadow })
    } catch (e) {
      console.warn('Sak 3D gagal dimuat, memakai gambar biasa.', e)
      this.teardown()
    }
  },
  beforeUnmount() {
    this._destroyed = true
    this.teardown()
  },
  methods: {
    webglAvailable() {
      try {
        const c = document.createElement('canvas')
        return !!(window.WebGLRenderingContext && (c.getContext('webgl2') || c.getContext('webgl')))
      } catch (_) {
        return false
      }
    },

    // Kotak bersegmen yang "digembungkan" jadi bentuk karung: tebal di tengah, pipih di jahitan tepi.
    sackGeometry(THREE) {
      const W = 1.6, H = 2.3, D = 0.62
      const g = new THREE.BoxGeometry(W, H, D, 48, 64, 10)
      const p = g.attributes.position
      const v = new THREE.Vector3()
      for (let i = 0; i < p.count; i++) {
        v.fromBufferAttribute(p, i)
        const u = Math.min(1, Math.abs(v.x) / (W / 2))
        const t = Math.min(1, Math.abs(v.y) / (H / 2))
        // ketebalan menipis ke tepi (bentuk bantal)
        const puff = Math.pow(1 - Math.pow(u, 3), 0.6) * Math.pow(1 - Math.pow(t, 6), 0.5)
        v.z *= puff
        // jahitan atas-bawah sedikit menyempit, isi turun ke bawah
        v.x *= 1 - 0.05 * Math.pow(t, 8) + (v.y < 0 ? 0.03 * Math.pow(t, 2) : 0)
        // kerutan halus kain karung
        const wr = 0.012 * Math.sin(v.x * 9 + v.y * 4) * Math.sin(v.y * 7) * puff
        v.z += Math.sign(v.z) * wr
        // ujung atas sedikit melengkung
        if (v.y > 0) v.y -= 0.05 * Math.pow(u, 2) * Math.pow(t, 6)
        p.setXYZ(i, v.x, v.y, v.z)
      }
      g.computeVertexNormals()
      return g
    },

    async build(THREE, tex) {
      const stage = this.$refs.stage
      const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true, powerPreference: 'low-power' })
      renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
      renderer.outputColorSpace = THREE.SRGBColorSpace
      renderer.toneMapping = THREE.ACESFilmicToneMapping
      renderer.toneMappingExposure = 1.05
      renderer.domElement.className = 'sack3d__canvas'
      renderer.domElement.setAttribute('role', 'img')
      renderer.domElement.setAttribute('aria-label', 'Sak Pupuk Organik Agrina 99 50 kg dalam tampilan 3D, bisa diputar')

      const scene = new THREE.Scene()
      const camera = new THREE.PerspectiveCamera(30, 1, 0.1, 50)
      camera.position.set(0, 0.25, 6.2)
      camera.lookAt(0, -0.05, 0)

      scene.add(new THREE.HemisphereLight(0xffffff, 0x3c6b48, 1.25))
      const key = new THREE.DirectionalLight(0xffffff, 2.2)
      key.position.set(2.5, 3, 4)
      scene.add(key)
      const rim = new THREE.DirectionalLight(0xbfffd6, 1.1)
      rim.position.set(-3, 1.5, -2.5)
      scene.add(rim)

      const maxAniso = renderer.capabilities.getMaxAnisotropy()
      const mkTex = (canvas) => {
        const t = new THREE.CanvasTexture(canvas)
        t.colorSpace = THREE.SRGBColorSpace
        t.anisotropy = Math.min(8, maxAniso)
        return t
      }
      const [front, back] = await Promise.all([tex.drawFront(), tex.drawBack()])
      if (this._destroyed) { renderer.dispose(); return }
      const plain = mkTex(tex.drawPlain())
      const mat = (map) => new THREE.MeshStandardMaterial({ map, roughness: 0.78, metalness: 0 })
      // urutan sisi BoxGeometry: +x, -x, +y, -y, +z (depan), -z (belakang)
      const backTex = mkTex(back)
      const materials = [mat(plain), mat(plain), mat(plain), mat(plain), mat(mkTex(front)), mat(backTex)]

      const sack = new THREE.Mesh(this.sackGeometry(THREE), materials)
      const group = new THREE.Group()
      group.add(sack)
      scene.add(group)

      const shadow = new THREE.Mesh(
        new THREE.PlaneGeometry(2.6, 1.1),
        new THREE.MeshBasicMaterial({ map: new THREE.CanvasTexture(tex.drawShadow()), transparent: true, depthWrite: false }),
      )
      shadow.rotation.x = -Math.PI / 2
      shadow.position.y = -1.42
      scene.add(shadow)

      stage.appendChild(renderer.domElement)
      this._three = { renderer, scene, camera, group, shadow, materials, sack }

      // ukuran mengikuti wadah
      const resize = () => {
        const w = stage.clientWidth || 1
        const h = stage.clientHeight || 1
        renderer.setSize(w, h, false)
        camera.aspect = w / h
        camera.updateProjectionMatrix()
      }
      resize()
      this._ro = new ResizeObserver(resize)
      this._ro.observe(stage)

      // interaksi geser (mouse dan sentuh)
      const reduce = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches
      const st = { rotY: -0.45, rotX: 0.05, velY: 0, lastX: 0, lastY: 0, idle: 0, pointer: null }
      const el = renderer.domElement
      el.style.touchAction = 'pan-y'
      const down = (e) => {
        st.pointer = e.pointerId
        st.lastX = e.clientX
        st.lastY = e.clientY
        st.velY = 0
        this.dragging = true
        this.hint = false
        el.setPointerCapture && el.setPointerCapture(e.pointerId)
      }
      const move = (e) => {
        if (st.pointer !== e.pointerId) return
        const dx = e.clientX - st.lastX
        const dy = e.clientY - st.lastY
        st.lastX = e.clientX
        st.lastY = e.clientY
        st.velY = dx * 0.012
        st.rotY += st.velY
        st.rotX = Math.max(-0.35, Math.min(0.35, st.rotX + dy * 0.004))
        st.idle = 0
      }
      const up = (e) => {
        if (st.pointer !== e.pointerId) return
        st.pointer = null
        this.dragging = false
      }
      el.addEventListener('pointerdown', down)
      el.addEventListener('pointermove', move)
      el.addEventListener('pointerup', up)
      el.addEventListener('pointercancel', up)
      this._listeners = [['pointerdown', down], ['pointermove', move], ['pointerup', up], ['pointercancel', up]]

      // hanya menggambar saat terlihat di layar (hemat baterai)
      let visible = true
      this._io = new IntersectionObserver(([entry]) => { visible = entry.isIntersecting })
      this._io.observe(stage)

      const clock = new THREE.Clock()
      const loop = () => {
        this._raf = requestAnimationFrame(loop)
        if (!visible || document.hidden) return
        const dt = Math.min(clock.getDelta(), 0.05)
        const t = clock.elapsedTime
        if (st.pointer === null) {
          st.rotY += st.velY
          st.velY *= 0.94
          st.idle += dt
          if (!reduce && st.idle > 1.2 && Math.abs(st.velY) < 0.002) st.rotY += dt * 0.35
          st.rotX += (0.05 - st.rotX) * 0.04
        }
        group.rotation.y = st.rotY
        group.rotation.x = st.rotX
        const bob = reduce ? 0 : Math.sin(t * 1.4) * 0.06
        group.position.y = bob
        shadow.scale.setScalar(1 - bob * 0.8)
        shadow.material.opacity = 0.75 - bob * 2
        renderer.render(scene, camera)
      }
      loop()
      this.ready = true
    },

    teardown() {
      cancelAnimationFrame(this._raf)
      this._ro && this._ro.disconnect()
      this._io && this._io.disconnect()
      const t = this._three
      if (!t) return
      const el = t.renderer.domElement
      ;(this._listeners || []).forEach(([n, f]) => el.removeEventListener(n, f))
      t.sack.geometry.dispose()
      t.materials.forEach((m) => { m.map && m.map.dispose(); m.dispose() })
      t.shadow.geometry.dispose()
      t.shadow.material.map.dispose()
      t.shadow.material.dispose()
      t.renderer.dispose()
      el.remove()
      this._three = null
    },
  },
}
</script>

<style scoped>
.sack3d { position: relative; width: 100%; aspect-ratio: 1 / 1.05; max-height: 560px; display: grid; place-items: center; cursor: grab; }
.sack3d--grab { cursor: grabbing; }
.sack3d :deep(.sack3d__canvas) { position: absolute; inset: 0; width: 100% !important; height: 100% !important; display: block; }
.sack3d__fallback { max-height: 88%; max-width: 80%; object-fit: contain; border-radius: 16px; box-shadow: 0 24px 50px rgba(0, 0, 0, .35); }
.sack3d__hint { position: absolute; top: 4px; left: 50%; transform: translateX(-50%); display: inline-flex; gap: 6px; align-items: center;
  background: rgba(255, 255, 255, .16); border: 1px solid rgba(255, 255, 255, .3); color: #fff; font-size: 13px; font-weight: 600;
  padding: 5px 12px; border-radius: 999px; pointer-events: none; backdrop-filter: blur(4px); white-space: nowrap; }
</style>
