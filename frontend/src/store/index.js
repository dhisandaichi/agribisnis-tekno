import { createStore } from 'vuex'

const API_URL = 'http://localhost:8080/api/v1'

export default createStore({
    state: {
        user: JSON.parse(localStorage.getItem('user')) || null, // { id, email, role, token }
        offlineVisits: [], // store pending visits when offline
    },
    mutations: {
        SET_USER(state, payload) {
            state.user = payload
            if (payload) {
                localStorage.setItem('user', JSON.stringify(payload))
            } else {
                localStorage.removeItem('user')
            }
        },
        ADD_OFFLINE_VISIT(state, visit) {
            state.offlineVisits.push(visit)
        },
        CLEAR_OFFLINE_VISITS(state) {
            state.offlineVisits = []
        }
    },
    actions: {
        async register(_, userData) {
            const response = await fetch(`${API_URL}/auth/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: userData.email,
                    password: userData.password,
                    role: userData.role
                })
            })

            const data = await response.json()
            if (!response.ok) {
                throw new Error(data.error || 'Terjadi kesalahan saat pendaftaran')
            }
            return data
        },

        async login({ commit }, credentials) {
            const response = await fetch(`${API_URL}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: credentials.email,
                    password: credentials.password
                })
            })

            const data = await response.json()
            if (!response.ok) {
                throw new Error(data.error || 'Login gagal. Cek kembali email/password.')
            }

            commit('SET_USER', {
                id: data.user.id,
                email: data.user.email,
                role: data.user.role,
                token: data.token
            })

            return data.user
        },

        logout({ commit }) {
            commit('SET_USER', null)
        },

        saveVisitOffline({ commit }, visit) {
            commit('ADD_OFFLINE_VISIT', visit)
        }
    },
    getters: {
        isAuthenticated: state => !!state.user,
        userRole: state => state.user ? state.user.role : null,
        token: state => state.user ? state.user.token : null,
    }
})
