import { createStore } from 'vuex'
import { supabase } from '../services/supabase'
import { getQueue, submitVisit, syncQueue } from '../services/visits'

const ROLE_HOME = {
    sales: '/sales',
    supervisor: '/supervisor',
    farmer: '/farmer'
}

export function homeFor(role) {
    return ROLE_HOME[role] || null
}

async function loadUser(authUser) {
    const { data, error } = await supabase
        .from('profiles')
        .select('id, email, full_name, role')
        .eq('id', authUser.id)
        .single()
    if (error) throw error
    return data
}

let initPromise = null

const store = createStore({
    state: {
        user: null, // { id, email, full_name, role }
        pendingVisits: 0 // kunjungan di antrean offline yang belum terkirim
    },
    mutations: {
        SET_USER(state, payload) {
            state.user = payload
        },
        SET_PENDING(state, count) {
            state.pendingVisits = count
        }
    },
    actions: {
        // Pulihkan sesi Supabase sekali saat app dibuka.
        init({ commit, dispatch }) {
            if (initPromise) return initPromise
            initPromise = (async () => {
                localStorage.removeItem('user') // sisa sesi versi backend Go
                supabase.auth.onAuthStateChange(event => {
                    if (event === 'SIGNED_OUT') commit('SET_USER', null)
                })
                const { data } = await supabase.auth.getSession()
                if (data.session) {
                    try {
                        commit('SET_USER', await loadUser(data.session.user))
                        dispatch('syncVisits')
                    } catch (err) {
                        console.error('Gagal memuat profil', err)
                        await supabase.auth.signOut()
                    }
                }
            })()
            return initPromise
        },

        // Mengembalikan user jika langsung login, atau null jika perlu konfirmasi email.
        async register({ commit }, { email, password, role, fullName }) {
            const { data, error } = await supabase.auth.signUp({
                email,
                password,
                options: {
                    data: { role, full_name: fullName },
                    emailRedirectTo: `${window.location.origin}/login`
                }
            })
            if (error) throw error
            if (!data.session) return null
            const user = await loadUser(data.user)
            commit('SET_USER', user)
            return user
        },

        async login({ commit, dispatch }, { email, password }) {
            const { data, error } = await supabase.auth.signInWithPassword({ email, password })
            if (error) {
                throw new Error(error.message === 'Invalid login credentials'
                    ? 'Email atau password salah.'
                    : error.message)
            }
            const user = await loadUser(data.user)
            commit('SET_USER', user)
            dispatch('syncVisits')
            return user
        },

        async logout({ commit }) {
            await supabase.auth.signOut()
            commit('SET_USER', null)
        },

        async submitVisit({ state, dispatch }, form) {
            const result = await submitVisit(state.user.id, form)
            await dispatch('refreshPending')
            return result
        },

        async refreshPending({ state, commit }) {
            if (!state.user) return
            const queue = await getQueue()
            commit('SET_PENDING', queue.filter(v => v.userId === state.user.id).length)
        },

        async syncVisits({ state, commit }) {
            if (!state.user || state.user.role !== 'sales') return
            commit('SET_PENDING', await syncQueue(state.user.id))
        }
    },
    getters: {
        isAuthenticated: state => !!state.user,
        userRole: state => state.user ? state.user.role : null
    }
})

window.addEventListener('online', () => store.dispatch('syncVisits'))

export default store
