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
        { path: 'discover', name: 'discover', component: () => import('@/pages/Discover.vue') },
        { path: 'events/:id', name: 'event-detail', component: () => import('@/pages/EventDetail.vue') },
        { path: 'commentary/:id', name: 'commentary-detail', component: () => import('@/pages/commentary/CommentaryDetail.vue') },
        { path: 'promoters', name: 'for-promoters', component: () => import('@/pages/ForPromoters.vue') },
        { path: 'use-cases', name: 'use-cases', component: () => import('@/pages/UseCases.vue') },
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

    // OAuth callback — standalone (no AuthLayout chrome)
    {
      path: '/auth/callback',
      name: 'auth-callback',
      component: () => import('@/pages/auth/OAuthCallback.vue'),
    },

    // ── Public: Try stream (no auth required) ─────────────────────────────────
    {
      path: '/try',
      name: 'try-stream',
      component: () => import('@/pages/TryStream.vue'),
    },

    // ── Public: Guest LiveKit room viewer (audio / audio+video) ──────────────
    // Must be defined BEFORE /guest/:guestId to avoid the dynamic segment
    // swallowing the static 'room' segment in some router edge cases.
    {
      path: '/guest/room/:guestId',
      name: 'guest-room-view',
      component: () => import('@/pages/GuestRoomView.vue'),
    },

    // ── Public: Guest stream viewer (shareable link) ───────────────────────────
    {
      path: '/guest/:guestId',
      name: 'guest-view',
      component: () => import('@/pages/GuestView.vue'),
    },

    // ── Protected: Watch ──────────────────────────────────────────────────────
    {
      path: '/watch/:eventId',
      name: 'watch',
      component: () => import('@/pages/Watch.vue'),
      meta: { requiresAuth: true },
    },

    // ── Protected: Commentary Room ────────────────────────────────────────────
    {
      path: '/commentary/:id/room',
      name: 'commentary-room',
      component: () => import('@/pages/commentary/CommentaryRoom.vue'),
      meta: { requiresAuth: true },
    },

    // ── Protected: Member dashboard ───────────────────────────────────────────
    {
      path: '/dashboard',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'dashboard', component: () => import('@/pages/dashboard/PromoterDashboard.vue') },
        { path: 'create', name: 'create-event', component: () => import('@/pages/dashboard/CreateEvent.vue') },
        { path: 'edit/:eventId', name: 'edit-event', component: () => import('@/pages/dashboard/EditEvent.vue') },
        { path: 'analytics/:eventId', name: 'analytics', component: () => import('@/pages/dashboard/Analytics.vue') },
        { path: 'revenue', name: 'revenue', component: () => import('@/pages/dashboard/Revenue.vue') },
        { path: 'streams', name: 'my-streams', component: () => import('@/pages/dashboard/MyStreams.vue') },
        { path: 'withdrawal', name: 'withdrawal', component: () => import('@/pages/dashboard/Withdrawal.vue') },
        // Redirect legacy commentary/create to the unified event creation wizard
        { path: 'commentary/create', redirect: '/dashboard/create' },
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

// Recover from stale PWA chunk failures (common after new deployments).
// When a lazy-loaded route chunk can't be fetched, Vue Router fails the
// navigation and the catch-all would redirect to '/'. Instead, force a
// full page reload so the service worker fetches fresh assets.
router.onError((error, to) => {
  const isChunkError =
    error?.message?.includes('Failed to fetch dynamically imported module') ||
    error?.message?.includes('error loading dynamically imported module') ||
    error?.message?.includes('Importing a module script failed')
  if (isChunkError) {
    const reloadKey = `chunk-reload:${to.fullPath}`
    if (!sessionStorage.getItem(reloadKey)) {
      sessionStorage.setItem(reloadKey, '1')
      window.location.href = to.fullPath
    }
  }
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

  next()
})

export default router
