<template>
  <div class="card p-8">
    <h2 class="text-white font-bold text-2xl mb-2">Sign in</h2>
    <p class="text-text-muted text-sm mb-8">Welcome back to ApexContact</p>

    <form @submit.prevent="handleLogin" class="space-y-5">
      <div>
        <label class="block text-text-muted text-sm mb-2">Email</label>
        <input v-model="form.email" type="email" placeholder="you@example.com" class="input" required />
      </div>
      <div>
        <label class="block text-text-muted text-sm mb-2">Password</label>
        <input v-model="form.password" type="password" placeholder="••••••••" class="input" required />
      </div>

      <div v-if="errorMsg" class="bg-status-error/10 border border-status-error/30 text-status-error
                                   text-sm rounded-lg px-4 py-3">
        {{ errorMsg }}
      </div>

      <button type="submit" class="btn-primary w-full" :disabled="loading">
        <span v-if="loading" class="flex items-center justify-center gap-2">
          <span class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
          Signing in...
        </span>
        <span v-else>Sign in</span>
      </button>
    </form>

    <p class="text-text-muted text-sm text-center mt-6">
      Don't have an account?
      <RouterLink to="/register" class="text-accent-red hover:text-accent-red-hover font-medium">
        Create one
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

async function handleLogin() {
  errorMsg.value = ''
  loading.value = true
  try {
    await auth.login(form.value.email, form.value.password)
    const redirect = (route.query.redirect as string) ?? '/'
    router.push(redirect)
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Invalid credentials'
  } finally {
    loading.value = false
  }
}
</script>
