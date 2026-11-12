// Tekstur kanvas untuk sak Pupuk Organik Agrina 99 3D.
// Digambar dengan kode (bukan foto) supaya tulisan tetap tajam di layar HP beresolusi tinggi.

const GREEN = '#0a9a46'
const GREEN_DARK = '#067a37'
const NAVY = '#2e3192'
const RED = '#e1251b'

function loadImage(src) {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => resolve(img)
    img.onerror = () => resolve(null)
    img.src = src
  })
}

// Anyaman karung plastik: garis-garis halus horizontal dan vertikal.
function weave(ctx, w, h, base = '#fbfbf8') {
  ctx.fillStyle = base
  ctx.fillRect(0, 0, w, h)
  ctx.globalAlpha = 0.07
  ctx.fillStyle = '#5b6a5f'
  for (let y = 0; y < h; y += 6) ctx.fillRect(0, y, w, 2)
  ctx.globalAlpha = 0.045
  for (let x = 0; x < w; x += 6) ctx.fillRect(x, 0, 2, h)
  ctx.globalAlpha = 1
}

function centerText(ctx, text, x, y, font, color) {
  ctx.font = font
  ctx.fillStyle = color
  ctx.textAlign = 'center'
  ctx.textBaseline = 'alphabetic'
  ctx.fillText(text, x, y)
}

async function ensureFonts() {
  if (!document.fonts || !document.fonts.load) return
  try {
    await Promise.all([
      document.fonts.load('800 120px Poppins'),
      document.fonts.load('700 40px Poppins'),
      document.fonts.load('600 30px Poppins'),
    ])
  } catch (_) { /* fallback ke font sistem */ }
}

const FONT = "Poppins, 'Segoe UI', Roboto, Arial, sans-serif"

export async function drawFront(W = 1024, H = 1472) {
  await ensureFonts()
  const logo = await loadImage('/brand/logo.png')
  const c = document.createElement('canvas')
  c.width = W
  c.height = H
  const ctx = c.getContext('2d')
  weave(ctx, W, H)

  // Judul
  centerText(ctx, 'Pupuk Organik', W / 2, 250, `800 118px ${FONT}`, GREEN)
  centerText(ctx, 'Agrina 99', W / 2, 390, `800 132px ${FONT}`, GREEN)

  // Pita hijau melengkung di bawah judul
  ctx.fillStyle = GREEN
  ctx.beginPath()
  ctx.moveTo(60, 452)
  ctx.quadraticCurveTo(W / 2, 410, W - 60, 430)
  ctx.lineTo(W - 60, 452)
  ctx.quadraticCurveTo(W / 2, 438, 60, 474)
  ctx.closePath()
  ctx.fill()

  // Logo daur ulang Agrina 99
  if (logo) ctx.drawImage(logo, W / 2 - 230, 500, 460, 460)

  // Lencana merah 2X
  ctx.save()
  ctx.translate(205, 1040)
  ctx.rotate(-0.12)
  ctx.fillStyle = RED
  ctx.beginPath()
  ctx.arc(0, 0, 118, 0, Math.PI * 2)
  ctx.fill()
  ctx.strokeStyle = '#fff'
  ctx.lineWidth = 6
  ctx.beginPath()
  ctx.arc(0, 0, 102, 0, Math.PI * 2)
  ctx.stroke()
  centerText(ctx, '2X', 0, 40, `800 112px ${FONT}`, '#fff')
  ctx.restore()

  // Label berat bersih
  ctx.fillStyle = RED
  roundRect(ctx, 600, 990, 340, 64, 10)
  ctx.fill()
  centerText(ctx, 'Berat Bersih 50 Kg', 770, 1034, `700 36px ${FONT}`, '#fff')

  // Tabel kandungan (dekoratif, garis tipis)
  ctx.strokeStyle = 'rgba(46,49,146,.35)'
  ctx.lineWidth = 2
  roundRect(ctx, 640, 1080, 300, 190, 8)
  ctx.stroke()
  ctx.fillStyle = 'rgba(46,49,146,.22)'
  for (let i = 0; i < 7; i++) {
    ctx.fillRect(660, 1104 + i * 23, 150 - (i % 3) * 20, 8)
    ctx.fillRect(850, 1104 + i * 23, 70, 8)
  }

  // Produsen
  centerText(ctx, 'Diproduksi Oleh :', 330, 1210, `600 26px ${FONT}`, NAVY)
  centerText(ctx, 'PT. Agribisnis Technology Indonesia', 330, 1246, `700 26px ${FONT}`, NAVY)
  centerText(ctx, 'Sumatera Barat - Indonesia', 330, 1280, `600 24px ${FONT}`, NAVY)

  // Pita bawah
  ctx.fillStyle = GREEN_DARK
  ctx.fillRect(0, H - 150, W, 12)
  centerText(ctx, 'JANGAN DIGANCU', W / 2, H - 80, `700 36px ${FONT}`, RED)

  return c
}

export async function drawBack(W = 1024, H = 1472) {
  await ensureFonts()
  const logo = await loadImage('/brand/logo.png')
  const c = document.createElement('canvas')
  c.width = W
  c.height = H
  const ctx = c.getContext('2d')
  weave(ctx, W, H)

  ctx.fillStyle = GREEN
  ctx.fillRect(0, 150, W, 12)
  if (logo) ctx.drawImage(logo, W / 2 - 170, 260, 340, 340)
  centerText(ctx, 'AGRINA 99', W / 2, 720, `800 110px ${FONT}`, GREEN)
  centerText(ctx, 'Pupuk Organik Premium', W / 2, 790, `600 46px ${FONT}`, NAVY)
  centerText(ctx, 'Penyubur Tanah, Daun dan Buah', W / 2, 900, `600 40px ${FONT}`, GREEN_DARK)

  ctx.fillStyle = GREEN
  roundRect(ctx, 130, 1020, W - 260, 230, 22)
  ctx.fill()
  centerText(ctx, 'agritechindonesia.id', W / 2, 1090, `700 40px ${FONT}`, '#fff')
  centerText(ctx, 'IG  @pupukorganikagrina99', W / 2, 1150, `600 36px ${FONT}`, '#fff')
  centerText(ctx, 'FB  Pupuk Organik Agrina', W / 2, 1210, `600 36px ${FONT}`, '#fff')

  ctx.fillStyle = GREEN_DARK
  ctx.fillRect(0, H - 150, W, 12)
  return c
}

export function drawPlain(W = 256, H = 256) {
  const c = document.createElement('canvas')
  c.width = W
  c.height = H
  weave(c.getContext('2d'), W, H, '#f4f4ef')
  return c
}

export function drawShadow(S = 256) {
  const c = document.createElement('canvas')
  c.width = S
  c.height = S
  const ctx = c.getContext('2d')
  const g = ctx.createRadialGradient(S / 2, S / 2, 0, S / 2, S / 2, S / 2)
  g.addColorStop(0, 'rgba(0,0,0,.55)')
  g.addColorStop(0.55, 'rgba(0,0,0,.18)')
  g.addColorStop(1, 'rgba(0,0,0,0)')
  ctx.fillStyle = g
  ctx.fillRect(0, 0, S, S)
  return c
}

function roundRect(ctx, x, y, w, h, r) {
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.arcTo(x + w, y, x + w, y + h, r)
  ctx.arcTo(x + w, y + h, x, y + h, r)
  ctx.arcTo(x, y + h, x, y, r)
  ctx.arcTo(x, y, x + w, y, r)
  ctx.closePath()
}
