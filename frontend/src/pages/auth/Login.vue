<template>
  <div class="w-full max-w-md mx-auto">
    <!-- Heading -->
    <h1 class="text-3xl font-bold text-white">Welcome back</h1>
    <p class="text-text-muted text-sm mt-2">Sign in to continue watching and learning</p>

    <!-- Form -->
    <form @submit.prevent="handleLogin" class="mt-10 space-y-5">
      <!-- Email -->
      <div>
        <label for="login-email" class="block text-sm font-medium text-text-muted mb-1.5">
          Email
        </label>
        <input
          id="login-email"
          v-model="form.email"
          type="email"
          placeholder="you@example.com"
          class="input"
          required
          autocomplete="email"
        />
      </div>

      <!-- Password -->
      <div>
        <label for="login-password" class="block text-sm font-medium text-text-muted mb-1.5">
          Password
        </label>
        <div class="relative">
          <input
            id="login-password"
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="••••••••"
            class="input pr-12"
            required
            autocomplete="current-password"
          />
          <!-- Show / hide toggle -->
          <button
            type="button"
            @click="showPassword = !showPassword"
            class="absolute right-3.5 top-1/2 -translate-y-1/2 text-text-muted hover:text-white
                   transition-colors focus:outline-none"
            :aria-label="showPassword ? 'Hide password' : 'Show password'"
          >
            <!-- Eye open -->
            <svg
              v-if="!showPassword"
              xmlns="http://www.w3.org/2000/svg"
              class="w-5 h-5"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
              <circle cx="12" cy="12" r="3" />
            </svg>
            <!-- Eye closed -->
            <svg
              v-else
              xmlns="http://www.w3.org/2000/svg"
              class="w-5 h-5"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
              <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
              <line x1="1" y1="1" x2="23" y2="23" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Error message -->
      <div
        v-if="errorMsg"
        class="flex items-center gap-2.5 bg-accent-red/10 border border-accent-red/20
               text-accent-red rounded-xl p-3 text-sm"
      >
        <!-- X icon -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="w-4 h-4 flex-shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="10" />
          <line x1="15" y1="9" x2="9" y2="15" />
          <line x1="9" y1="9" x2="15" y2="15" />
        </svg>
        {{ errorMsg }}
      </div>

      <!-- Submit -->
      <button
        type="submit"
        class="btn-primary w-full justify-center"
        :disabled="loading"
      >
        <span
          v-if="loading"
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin flex-shrink-0"
        />
        {{ loading ? 'Signing in…' : 'Sign in' }}
      </button>
    </form>

    <!-- Divider -->
    <div class="flex items-center gap-3 my-6">
      <div class="flex-1 h-px bg-white/10" />
      <span class="text-text-muted text-xs font-medium">or</span>
      <div class="flex-1 h-px bg-white/10" />
    </div>

    <!-- Google social login (coming soon) -->
    <div class="relative group">
      <button
        type="button"
        disabled
        class="btn-ghost w-full justify-center opacity-50 cursor-not-allowed"
        aria-disabled="true"
      >
        <!-- Google G icon -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="w-5 h-5 flex-shrink-0"
          viewBox="0 0 24 24"
        >
          <path
            fill="#4285F4"
            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
          />
          <path
            fill="#34A853"
            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
          />
          <path
            fill="#FBBC05"
            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"
          />
          <path
            fill="#EA4335"
            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
          />
        </svg>
        Continue with Google
      </button>
      <!-- Tooltip -->
      <div
        class="absolute -top-9 left-1/2 -translate-x-1/2 whitespace-nowrap
               bg-bg-elevated border border-white/10 text-white/80 text-xs
               rounded-lg px-3 py-1.5 pointer-events-none
               opacity-0 group-hover:opacity-100 transition-opacity duration-200"
      >
        Coming soon
      </div>
    </div>

    <!-- Footer -->
    <p class="text-text-muted text-sm text-center mt-8">
      Don't have an account?
      <RouterLink
        to="/register"
        class="text-accent-red hover:text-accent-red-hover font-medium transition-colors ml-1"
      >
        Sign up free
      </RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const form = ref({ email: '', password: '' })
const errorMsg = ref('')
const loading = ref(false)
const showPassword = ref(false)

async function handleLogin() {
  errorMsg.value = ''
  loading.value = true
  try {
    await auth.login(form.value.email, form.value.password)
    const redirect = (route.query.redirect as string) ?? '/'
    router.push(redirect)
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Invalid credentials. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
