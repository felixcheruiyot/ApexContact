<template>
  <div class="card p-8">
    <h2 class="text-white font-bold text-2xl mb-2">Create account</h2>
    <p class="text-text-muted text-sm mb-8">Join Live Streamify and watch live events</p>

    <form @submit.prevent="handleRegister" class="space-y-5">
      <div>
        <label class="block text-text-muted text-sm mb-2">Full Name</label>
        <input v-model="form.full_name" type="text" placeholder="John Doe" class="input" required />
      </div>
      <div>
        <label class="block text-text-muted text-sm mb-2">Email</label>
        <input v-model="form.email" type="email" placeholder="you@example.com" class="input" required />
      </div>
      <div>
        <label class="block text-text-muted text-sm mb-2">Phone (M-Pesa number)</label>
        <input v-model="form.phone" type="tel" placeholder="+254712345678" class="input" required />
      </div>
      <div>
        <label class="block text-text-muted text-sm mb-2">Password</label>
        <input v-model="form.password" type="password" placeholder="At least 8 characters" class="input"
               required minlength="8" />
      </div>

      <div v-if="errorMsg" class="bg-status-error/10 border border-status-error/30 text-status-error
                                   text-sm rounded-lg px-4 py-3">
        {{ errorMsg }}
      </div>

      <button type="submit" class="btn-primary w-full" :disabled="loading">
        <span v-if="loading" class="flex items-center justify-center gap-2">
          <span class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
          Creating account...
        </span>
        <span v-else>Create Account</span>
      </button>
    </form>

    <p class="text-text-muted text-sm text-center mt-6">
      Already have an account?
      <RouterLink to="/login" class="text-accent-red hover:text-accent-red-hover font-medium">
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

const form = ref({ full_name: '', email: '', phone: '', password: '' })
const errorMsg = ref('')
const loading = ref(false)

async function handleRegister() {
  errorMsg.value = ''
  loading.value = true
  try {
    await auth.register(form.value.email, form.value.password, form.value.full_name, form.value.phone)
    router.push('/')
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>
