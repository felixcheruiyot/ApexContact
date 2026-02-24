<template>
  <div class="min-h-screen bg-bg flex items-center justify-center">
    <div class="flex flex-col items-center gap-4 text-center px-4">
      <template v-if="!errorMsg">
        <div class="w-12 h-12 border-4 border-accent-red border-t-transparent rounded-full animate-spin" />
        <p class="text-text-muted text-sm">Signing you in…</p>
      </template>
      <template v-else>
        <div class="w-14 h-14 rounded-full bg-white/5 flex items-center justify-center">
          <svg class="w-7 h-7 text-accent-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10" stroke-width="2" />
            <line x1="15" y1="9" x2="9" y2="15" stroke-width="2.5" stroke-linecap="round" />
            <line x1="9" y1="9" x2="15" y2="15" stroke-width="2.5" stroke-linecap="round" />
          </svg>
        </div>
        <div>
          <p class="text-white font-semibold mb-1">Sign-in failed</p>
          <p class="text-text-muted text-sm max-w-xs">{{ errorMsg }}</p>
        </div>
        <RouterLink
          to="/login"
          class="px-5 py-2.5 rounded-lg border border-white/20 hover:border-white/40
                 text-white text-sm font-medium transition-all hover:bg-white/5"
        >
          Back to login
        </RouterLink>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const errorMsg = ref('')

onMounted(async () => {
  const code = route.query.code as string | undefined
  const error = route.query.error as string | undefined

  if (error || !code) {
    errorMsg.value = error === 'access_denied'
      ? 'You cancelled the Google sign-in.'
      : 'Google sign-in was unsuccessful. Please try again.'
    return
  }

  const redirectUri = `${window.location.origin}/auth/callback`
  const destination = sessionStorage.getItem('oauth_redirect') || '/'
  sessionStorage.removeItem('oauth_redirect')

  try {
    await auth.loginWithGoogle(code, redirectUri)
    router.replace(destination)
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Could not complete sign-in. Please try again.'
  }
})
</script>
