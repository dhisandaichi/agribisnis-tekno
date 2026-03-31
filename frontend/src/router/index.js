import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'

const routes = [
    {
        path: '/',
        name: 'Home',
        redirect: '/login'
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue')
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/Register.vue')
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
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Route Guard to check credentials
router.beforeEach((to, from, next) => {
    if (to.meta.requiresAuth && !store.getters.isAuthenticated) {
        next('/login')
    } else if (to.meta.requiresAuth && store.getters.isAuthenticated) {
        const role = store.getters.userRole
        // Check if role is allowed (Optional Strict Validation)
        if (to.meta.role && to.meta.role !== role) {
            // Basic fallback
            if (role === 'sales') next('/sales')
            else if (role === 'supervisor') next('/supervisor')
            else if (role === 'farmer') next('/farmer')
            else next('/login')
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router
