import { createRouter, createWebHistory } from 'vue-router'
import store, { homeFor } from '../store'

const routes = [
    {
        path: '/',
        name: 'Home',
        redirect: '/login'
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/Register.vue'),
        meta: { guestOnly: true }
    },
    {
        path: '/sales',
        name: 'SalesVisit',
        component: () => import('../views/SalesVisit.vue'),
        meta: { requiresAuth: true, role: 'sales' }
    },
    {
        path: '/supervisor',
        name: 'Supervisor',
        component: () => import('../views/SupervisorDashboard.vue'),
        meta: { requiresAuth: true, role: 'supervisor' }
    },
    {
        path: '/farmer',
        name: 'Farmer',
        component: () => import('../views/FarmerEducation.vue'),
        meta: { requiresAuth: true, role: 'farmer' }
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/login'
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Route guard: tunggu sesi Supabase dipulihkan, lalu cek login & role.
router.beforeEach(async to => {
    await store.dispatch('init')

    const user = store.state.user
    const home = user ? homeFor(user.role) : null

    if (user && !home) {
        // Role tanpa halaman (mis. distributor) belum didukung di PoC.
        await store.dispatch('logout')
        return { path: '/login', query: { reason: 'role' } }
    }
    if (to.meta.guestOnly && user) return home
    if (to.meta.requiresAuth && !user) return { path: '/login', query: { redirect: to.fullPath } }
    if (to.meta.role && user && to.meta.role !== user.role) return home
    return true
})

export default router
