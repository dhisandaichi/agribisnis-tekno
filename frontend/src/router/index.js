import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'
import { homeForRole } from '../services/format'

const routes = [
  { path: '/', name: 'Home', redirect: '/login' },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { guest: true } },
  { path: '/register', name: 'Register', component: () => import('../views/Register.vue'), meta: { guest: true } },
  { path: '/sales', name: 'SalesVisit', component: () => import('../views/SalesVisit.vue'), meta: { requiresAuth: true, role: 'sales' } },
  { path: '/supervisor', name: 'Supervisor', component: () => import('../views/SupervisorDashboard.vue'), meta: { requiresAuth: true, role: 'supervisor' } },
  { path: '/farmer', name: 'Farmer', component: () => import('../views/FarmerEducation.vue'), meta: { requiresAuth: true, role: 'farmer' } },
  { path: '/distributor', name: 'Distributor', component: () => import('../views/DistributorPortal.vue'), meta: { requiresAuth: true, role: 'distributor' } },
  { path: '/:pathMatch(.*)*', redirect: '/login' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

let initPromise = null

router.beforeEach(async (to) => {
  if (!store.state.ready) {
    initPromise = initPromise || store.dispatch('init')
    await initPromise
  }
  const authed = store.getters.isAuthenticated
  const role = store.getters.userRole

  if (to.meta.requiresAuth && !authed) return '/login'
  if (to.meta.guest && authed) return homeForRole(role)
  if (to.meta.role && to.meta.role !== role) return homeForRole(role)
  return true
})

export default router
