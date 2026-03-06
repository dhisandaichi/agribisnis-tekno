import { createStore } from 'vuex'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

// Helper: authed fetch
async function authFetch(url, token, options = {}) {
    const res = await fetch(url, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            ...(options.headers || {})
        }
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Request gagal')
    return data
}

export default createStore({
    state: {
        user: JSON.parse(localStorage.getItem('user')) || null,
        offlineVisits: JSON.parse(localStorage.getItem('offlineVisits')) || [],

        // Dynamic data states
        visits: [],
        distributors: [],
        contents: [],
        notifications: [],
        weather: null,
        commodityPrices: [],

        // UI States
        loading: false,
    },

    mutations: {
        SET_USER(state, payload) {
            state.user = payload
            if (payload) localStorage.setItem('user', JSON.stringify(payload))
            else localStorage.removeItem('user')
        },
        ADD_OFFLINE_VISIT(state, visit) {
            state.offlineVisits.push(visit)
            localStorage.setItem('offlineVisits', JSON.stringify(state.offlineVisits))
        },
        CLEAR_OFFLINE_VISITS(state) {
            state.offlineVisits = []
            localStorage.removeItem('offlineVisits')
        },
        SET_VISITS(state, visits) { state.visits = visits },
        SET_DISTRIBUTORS(state, data) { state.distributors = data },
        SET_CONTENTS(state, data) { state.contents = data },
        SET_NOTIFICATIONS(state, data) { state.notifications = data },
        SET_WEATHER(state, data) { state.weather = data },
        SET_COMMODITY_PRICES(state, data) { state.commodityPrices = data },
        SET_LOADING(state, val) { state.loading = val },
    },

    actions: {
        // ─── AUTH ──────────────────────────────────────────────
        async register(_, userData) {
            const res = await fetch(`${API_URL}/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: userData.email, password: userData.password, role: userData.role })
            })
            const data = await res.json()
            if (!res.ok) throw new Error(data.error || 'Terjadi kesalahan saat pendaftaran')
            return data
        },

        async login({ commit }, credentials) {
            const res = await fetch(`${API_URL}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: credentials.email, password: credentials.password })
            })
            const data = await res.json()
            if (!res.ok) throw new Error(data.error || 'Login gagal. Cek kembali email/password.')
            commit('SET_USER', { id: data.user.id, email: data.user.email, role: data.user.role, token: data.token })
            return data.user
        },

        logout({ commit }) {
            commit('SET_USER', null)
        },

        // ─── VISITS ────────────────────────────────────────────
        saveVisitOffline({ commit }, visit) {
            commit('ADD_OFFLINE_VISIT', visit)
        },

        async fetchVisits({ commit, getters }) {
            commit('SET_LOADING', true)
            try {
                const role = getters.userRole
                const endpoint = role === 'supervisor'
                    ? `${API_URL}/supervisor/visits`
                    : `${API_URL}/sales/visits`
                const data = await authFetch(endpoint, getters.token)
                commit('SET_VISITS', data.data || [])
            } catch (e) { console.error('fetchVisits:', e) }
            finally { commit('SET_LOADING', false) }
        },

        // ─── DISTRIBUTORS ──────────────────────────────────────
        async fetchDistributors({ commit, getters }) {
            commit('SET_LOADING', true)
            try {
                const data = await authFetch(`${API_URL}/distributors/`, getters.token)
                commit('SET_DISTRIBUTORS', data.data || [])
            } catch (e) { console.error('fetchDistributors:', e) }
            finally { commit('SET_LOADING', false) }
        },

        // ─── EDUCATION ─────────────────────────────────────────
        async fetchContents({ commit }, category = '') {
            try {
                const url = category
                    ? `${API_URL}/education/content?category=${category}`
                    : `${API_URL}/education/content`
                const res = await fetch(url)
                const data = await res.json()
                commit('SET_CONTENTS', data.data || [])
            } catch (e) { console.error('fetchContents:', e) }
        },

        // ─── NOTIFICATIONS ─────────────────────────────────────
        async fetchNotifications({ commit, getters }) {
            try {
                const data = await authFetch(`${API_URL}/notifications/`, getters.token)
                commit('SET_NOTIFICATIONS', data.data || [])
            } catch (e) { console.error('fetchNotifications:', e) }
        },

        async markNotificationRead({ getters }, id) {
            try {
                await authFetch(`${API_URL}/notifications/${id}/read`, getters.token, { method: 'PUT' })
            } catch (e) { console.error('markRead:', e) }
        },

        // ─── WEATHER ────────────────────────────────────────────
        async fetchWeather({ commit }, city = 'jakarta') {
            try {
                const res = await fetch(`${API_URL}/integration/weather?city=${city}`)
                const data = await res.json()
                commit('SET_WEATHER', data.data || null)
            } catch (e) { console.error('fetchWeather:', e) }
        },

        // ─── COMMODITY PRICES ──────────────────────────────────
        async fetchCommodityPrices({ commit }) {
            try {
                const res = await fetch(`${API_URL}/integration/commodity-prices`)
                const data = await res.json()
                commit('SET_COMMODITY_PRICES', data.data || [])
            } catch (e) { console.error('fetchCommodityPrices:', e) }
        },
    },

    getters: {
        isAuthenticated: state => !!state.user,
        userRole: state => state.user ? state.user.role : null,
        token: state => state.user ? state.user.token : null,
        unreadNotifications: state => state.notifications.filter(n => n.status === 'pending').length,
    }
})
