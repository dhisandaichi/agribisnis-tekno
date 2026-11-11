// Klien API backend Go. Default: same-origin "/api/v1" (Vercel). Untuk dev lokal set VITE_API_URL.
const BASE = (import.meta.env.VITE_API_URL || '/api/v1').replace(/\/$/, '')
const TOKEN_KEY = 'agritec_token'

export const auth = {
  get token() {
    try { return localStorage.getItem(TOKEN_KEY) } catch { return null }
  },
  set token(v) {
    try { v ? localStorage.setItem(TOKEN_KEY, v) : localStorage.removeItem(TOKEN_KEY) } catch { /* abaikan */ }
  },
}

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
  }
}

let onUnauthorized = () => {}
export function setUnauthorizedHandler(fn) { onUnauthorized = fn }

export async function api(path, { method = 'GET', body, raw = false } = {}) {
  const headers = {}
  if (body !== undefined) headers['Content-Type'] = 'application/json'
  if (auth.token) headers.Authorization = `Bearer ${auth.token}`
  let res
  try {
    res = await fetch(BASE + path, { method, headers, body: body !== undefined ? JSON.stringify(body) : undefined })
  } catch {
    throw new ApiError('Tidak ada koneksi internet / server tidak terjangkau.', 0)
  }
  if (res.status === 401 && auth.token && !path.startsWith('/auth/')) onUnauthorized()
  if (raw && res.ok) return res
  const text = await res.text()
  let data = null
  try { data = text ? JSON.parse(text) : null } catch { data = null }
  if (!res.ok) throw new ApiError((data && data.error) || `Terjadi kesalahan (${res.status})`, res.status)
  return data
}

export const get = (p) => api(p)
export const post = (p, body = {}) => api(p, { method: 'POST', body })
export const put = (p, body = {}) => api(p, { method: 'PUT', body })
export const patch = (p, body = {}) => api(p, { method: 'PATCH', body })
export const del = (p) => api(p, { method: 'DELETE' })

// Ambil foto kunjungan (butuh header Authorization) sebagai object URL
export async function photoUrl(visitId) {
  const res = await api(`/visits/${visitId}/photo`, { raw: true })
  return URL.createObjectURL(await res.blob())
}
