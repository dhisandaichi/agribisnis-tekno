import { createStore } from 'vuex'
import { auth, get, post, ApiError, setUnauthorizedHandler } from '../services/api'

const QUEUE_KEY = 'agritec_offline_visits'
const USER_KEY = 'agritec_user'

function load(key, fallback) {
  try { return JSON.parse(localStorage.getItem(key)) ?? fallback } catch { return fallback }
}
function save(key, value) {
  try { localStorage.setItem(key, JSON.stringify(value)); return true } catch { return false }
}

const toPayload = (v) => ({
  client_ref: v.clientRef,
  farmer_name: v.farmerName,
  farmer_phone: v.farmerPhone || '',
  village: v.village || '',
  crop_type: v.cropType,
  land_area: Number(v.landArea),
  fertilizer_usage: v.fertilizerUsage === '' || v.fertilizerUsage == null ? null : Number(v.fertilizerUsage),
  competitor_product: v.competitorProduct || '',
  notes: v.notes || '',
  latitude: v.latitude,
  longitude: v.longitude,
  visit_time: v.createdAt,
  photo: v.photoDataUrl,
})

const isNetworkError = (err) => err instanceof ApiError && (err.status === 0 || err.status >= 502)

const store = createStore({
  state: {
    user: auth.token ? load(USER_KEY, null) : null,
    ready: false,
    online: typeof navigator === 'undefined' ? true : navigator.onLine,
    offlineVisits: load(QUEUE_KEY, []),
    syncing: false,
  },
  mutations: {
    SET_USER(state, user) {
      state.user = user
      if (user) save(USER_KEY, user)
      else { try { localStorage.removeItem(USER_KEY) } catch { /* abaikan */ } }
    },
    SET_READY(state) { state.ready = true },
    SET_ONLINE(state, v) { state.online = v },
    SET_SYNCING(state, v) { state.syncing = v },
    SET_QUEUE(state, q) { state.offlineVisits = q },
  },
  actions: {
    async init({ commit, dispatch, state }) {
      setUnauthorizedHandler(() => {
        auth.token = null
        commit('SET_USER', null)
        if (!['/', '/login', '/register'].includes(window.location.pathname)) window.location.assign('/login')
      })
      window.addEventListener('online', () => { commit('SET_ONLINE', true); dispatch('syncOfflineVisits') })
      window.addEventListener('offline', () => commit('SET_ONLINE', false))
      // Perbarui data user (role terbaru) bila online; saat offline pakai data tersimpan
      if (auth.token && state.online) {
        try {
          commit('SET_USER', await get('/me'))
        } catch (err) {
          if (!isNetworkError(err)) { auth.token = null; commit('SET_USER', null) }
        }
      }
      commit('SET_READY')
      dispatch('syncOfflineVisits')
    },

    async register(_, { email, password, role, fullName, phone }) {
      return post('/auth/register', { email, password, role, full_name: fullName, phone })
    },

    async login({ commit }, { email, password }) {
      const data = await post('/auth/login', { email, password })
      auth.token = data.token
      commit('SET_USER', data.user)
      return data.user
    },

    logout({ commit }) {
      auth.token = null
      commit('SET_USER', null)
    },

    // Kirim langsung bila online; simpan di antrean perangkat bila offline / jaringan gagal
    async submitVisit({ state, commit }, visit) {
      if (state.online) {
        try {
          await post('/sales/visits', toPayload(visit))
          return { queued: false }
        } catch (err) {
          if (!isNetworkError(err)) throw err
        }
      }
      const queue = [...state.offlineVisits, { ...visit, userId: state.user?.id }]
      if (!save(QUEUE_KEY, queue)) throw new Error('Penyimpanan perangkat penuh. Kirim data tertunda terlebih dahulu.')
      commit('SET_QUEUE', queue)
      return { queued: true }
    },

    async syncOfflineVisits({ state, commit }) {
      if (!state.user || state.user.role !== 'sales' || state.syncing || !state.online || !state.offlineVisits.length) {
        return { synced: 0, failed: 0 }
      }
      commit('SET_SYNCING', true)
      let synced = 0
      const remaining = []
      for (const visit of state.offlineVisits) {
        if (visit.userId && visit.userId !== state.user.id) { remaining.push(visit); continue }
        try {
          await post('/sales/visits', toPayload(visit))
          synced++
        } catch (err) {
          // data tidak valid (4xx) tidak akan berhasil walau dicoba ulang -> tetap disimpan agar tidak hilang
          remaining.push(err instanceof ApiError && err.status >= 400 && err.status < 500 ? { ...visit, error: err.message } : visit)
        }
      }
      save(QUEUE_KEY, remaining)
      commit('SET_QUEUE', remaining)
      commit('SET_SYNCING', false)
      return { synced, failed: remaining.length }
    },
  },
  getters: {
    isAuthenticated: (state) => !!state.user && !!auth.token,
    userRole: (state) => state.user?.role || null,
    userId: (state) => state.user?.id || null,
    pendingCount: (state) => state.offlineVisits.length,
  },
})

export default store
