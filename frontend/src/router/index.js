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
    },
    {
        path: '/distributor',
        name: 'Distributor',
        component: () => import('../views/DistributorManagement.vue'),
        meta: { requiresAuth: true, role: 'distributor' }
    },
    {
        // Catch-all 404
        path: '/:pathMatch(.*)*',
        redirect: '/login'
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
        if (to.meta.role && to.meta.role !== role) {
            // Role-based redirect fallback
            if (role === 'sales') next('/sales')
            else if (role === 'supervisor') next('/supervisor')
            else if (role === 'farmer') next('/farmer')
            else if (role === 'distributor') next('/distributor')
            else next('/login')
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router
