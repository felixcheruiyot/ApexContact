import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: [
    // ── Public ────────────────────────────────────────────────────────────────
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      children: [
        { path: '', name: 'home', component: () => import('@/pages/Home.vue') },
        { path: 'events/:id', name: 'event-detail', component: () => import('@/pages/EventDetail.vue') },
        { path: 'promoters', name: 'for-promoters', component: () => import('@/pages/ForPromoters.vue') },
        { path: 'privacy', name: 'privacy', component: () => import('@/pages/PrivacyPolicy.vue') },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/pages/Profile.vue'),
          meta: { requiresAuth: true },
        },
      ],
    },

    // ── Auth ──────────────────────────────────────────────────────────────────
    {
      path: '/',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        { path: 'login', name: 'login', component: () => import('@/pages/auth/Login.vue') },
        { path: 'register', name: 'register', component: () => import('@/pages/auth/Register.vue') },
      ],
    },

    // ── Protected: Watch ──────────────────────────────────────────────────────
    {
      path: '/watch/:eventId',
      name: 'watch',
      component: () => import('@/pages/Watch.vue'),
      meta: { requiresAuth: true },
    },

    // ── Protected: Promoter / Broadcaster dashboard ───────────────────────────
    {
      path: '/dashboard',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true, requiresRole: 'promoter' },
      children: [
        { path: '', name: 'dashboard', component: () => import('@/pages/dashboard/PromoterDashboard.vue') },
        { path: 'create', name: 'create-event', component: () => import('@/pages/dashboard/CreateEvent.vue') },
        { path: 'edit/:eventId', name: 'edit-event', component: () => import('@/pages/dashboard/EditEvent.vue') },
        { path: 'analytics/:eventId', name: 'analytics', component: () => import('@/pages/dashboard/Analytics.vue') },
        { path: 'revenue', name: 'revenue', component: () => import('@/pages/dashboard/Revenue.vue') },
      ],
    },

    // ── Protected: Admin ──────────────────────────────────────────────────────
    {
      path: '/admin',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true, requiresRole: 'admin' },
      children: [
        { path: '', name: 'admin', component: () => import('@/pages/admin/AdminDashboard.vue') },
        { path: 'fraud', name: 'fraud', component: () => import('@/pages/admin/FraudMonitor.vue') },
        { path: 'users', name: 'users', component: () => import('@/pages/admin/UserManagement.vue') },
        { path: 'events', name: 'admin-events', component: () => import('@/pages/admin/EventManagement.vue') },
      ],
    },

    // ── Catch-all ─────────────────────────────────────────────────────────────
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

// Route guards
router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next({ name: 'login', query: { redirect: to.fullPath } })
  }

  if (to.meta.requiresRole === 'admin' && !auth.isAdmin) {
    return next({ name: 'home' })
  }

  // 'promoter' role guard: allows promoters, broadcasters, and admins
  if (to.meta.requiresRole === 'promoter' && !auth.isPromoter) {
    return next({ name: 'home' })
  }

  next()
})

export default router
