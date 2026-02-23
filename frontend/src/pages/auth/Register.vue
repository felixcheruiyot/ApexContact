<template>
  <div class="w-full max-w-md mx-auto">
    <!-- Heading -->
    <h1 class="text-3xl font-bold text-white">Create your account</h1>
    <p class="text-text-muted text-sm mt-2">Free to join. Start watching or hosting today.</p>

    <!-- Form -->
    <form @submit.prevent="handleRegister" class="mt-10 space-y-5">
      <!-- Full Name -->
      <div>
        <label for="reg-name" class="block text-sm font-medium text-text-muted mb-1.5">
          Full Name
        </label>
        <input
          id="reg-name"
          v-model="form.full_name"
          type="text"
          placeholder="Your full name"
          class="input"
          required
          autocomplete="name"
        />
      </div>

      <!-- Email -->
      <div>
        <label for="reg-email" class="block text-sm font-medium text-text-muted mb-1.5">
          Email
        </label>
        <input
          id="reg-email"
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
        <label for="reg-password" class="block text-sm font-medium text-text-muted mb-1.5">
          Password
        </label>
        <div class="relative">
          <input
            id="reg-password"
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="Min. 8 characters"
            class="input pr-12"
            required
            minlength="8"
            autocomplete="new-password"
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
        <!-- Data note -->
        <p class="text-text-muted text-xs mt-2 leading-relaxed">
          &#128274; We'll never share your data. Add phone number and details on your profile after.
        </p>
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

      <!-- Benefits row -->
      <div class="flex flex-wrap gap-x-4 gap-y-1.5">
        <span class="text-text-muted text-xs flex items-center gap-1">
          <span class="text-status-success font-bold">&#10003;</span>
          Free to browse all events
        </span>
        <span class="text-text-muted text-xs flex items-center gap-1">
          <span class="text-status-success font-bold">&#10003;</span>
          Pay only what you watch
        </span>
        <span class="text-text-muted text-xs flex items-center gap-1">
          <span class="text-status-success font-bold">&#10003;</span>
          Earn 70% as a creator
        </span>
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
        {{ loading ? 'Creating account…' : 'Create free account' }}
      </button>
    </form>

    <!-- Footer -->
    <p class="text-text-muted text-sm text-center mt-8">
      Already have an account?
      <RouterLink
        to="/login"
        class="text-accent-red hover:text-accent-red-hover font-medium transition-colors ml-1"
      >
        Sign in
      </RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const form = ref({ full_name: '', email: '', password: '' })
const errorMsg = ref('')
const loading = ref(false)
const showPassword = ref(false)

async function handleRegister() {
  errorMsg.value = ''
  loading.value = true
  try {
    await auth.register(form.value.email, form.value.password, form.value.full_name)
    router.push('/')
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
