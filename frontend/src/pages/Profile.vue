<template>
  <div class="max-w-2xl mx-auto py-12 px-4">

    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-white font-bold text-2xl">My Profile</h1>
      <p class="text-text-muted text-sm mt-1">Manage your personal details and preferences</p>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="card p-8 space-y-5 animate-pulse">
      <div class="h-5 bg-white/10 rounded w-1/3" />
      <div class="h-10 bg-white/10 rounded" />
      <div class="h-10 bg-white/10 rounded" />
      <div class="grid grid-cols-2 gap-4">
        <div class="h-10 bg-white/10 rounded" />
        <div class="h-10 bg-white/10 rounded" />
      </div>
    </div>

    <!-- Form -->
    <form v-else @submit.prevent="handleSave" class="space-y-6">

      <!-- Avatar + name section -->
      <div class="card p-6 flex items-center gap-5">
        <div class="w-16 h-16 rounded-full bg-accent-red/20 flex items-center justify-center shrink-0">
          <span class="text-accent-red text-xl font-bold">{{ initials }}</span>
        </div>
        <div>
          <p class="text-white font-semibold text-lg">{{ form.full_name || auth.user?.full_name }}</p>
          <p class="text-text-muted text-sm">{{ auth.user?.email }}</p>
          <span class="inline-block mt-1.5 text-xs font-medium px-2 py-0.5 rounded-full capitalize"
            :class="roleBadgeClass">{{ auth.isAdmin ? 'Admin' : 'Member' }}</span>
        </div>
      </div>

      <!-- Personal info -->
      <div class="card p-6 space-y-5">
        <h2 class="text-white font-semibold text-base">Personal Information</h2>

        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Full Name</label>
          <input v-model="form.full_name" type="text" class="input w-full" placeholder="Your full name" required />
        </div>

        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Email</label>
          <input :value="auth.user?.email" type="email" class="input w-full opacity-50 cursor-not-allowed"
            disabled title="Email cannot be changed" />
          <p class="text-text-muted text-xs mt-1">Email address cannot be changed.</p>
        </div>

        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">
            Phone Number
            <span class="text-text-muted font-normal normal-case ml-1">(M-Pesa number for payments)</span>
          </label>
          <input v-model="form.phone" type="tel" class="input w-full" placeholder="+254712345678" />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Age</label>
            <input v-model.number="form.age" type="number" min="13" max="120" class="input w-full" placeholder="—" />
          </div>
          <div>
            <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Gender</label>
            <select v-model="form.gender" class="input w-full">
              <option value="">Prefer not to say</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="other">Other</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Country</label>
          <select v-model="form.country" class="input w-full">
            <option value="">Select country</option>
            <option value="Kenya">Kenya</option>
            <option value="Uganda">Uganda</option>
            <option value="Tanzania">Tanzania</option>
            <option value="Rwanda">Rwanda</option>
            <option value="Ethiopia">Ethiopia</option>
            <option value="Nigeria">Nigeria</option>
            <option value="South Africa">South Africa</option>
            <option value="Ghana">Ghana</option>
            <option value="United Kingdom">United Kingdom</option>
            <option value="United States">United States</option>
            <option value="Other">Other</option>
          </select>
        </div>
      </div>

      <!-- Feedback messages -->
      <div v-if="saveError" class="bg-status-error/10 border border-status-error/20 rounded-xl px-4 py-3 text-status-error text-sm">
        {{ saveError }}
      </div>
      <div v-if="saveSuccess" class="bg-status-success/10 border border-status-success/20 rounded-xl px-4 py-3 text-status-success text-sm flex items-center gap-2">
        <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        Profile updated successfully.
      </div>

      <div class="flex gap-3">
        <button type="submit" :disabled="saving" class="btn-primary py-2.5 px-6 flex items-center gap-2">
          <svg v-if="saving" class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ saving ? 'Saving…' : 'Save Changes' }}
        </button>
        <RouterLink to="/" class="btn-ghost py-2.5 px-6">Back to Home</RouterLink>
      </div>
    </form>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { profileApi } from '@/api/profile'

const auth = useAuthStore()

const loading = ref(true)
const saving = ref(false)
const saveError = ref<string | null>(null)
const saveSuccess = ref(false)

const form = ref({
  full_name: '',
  phone: '',
  age: null as number | null,
  gender: '',
  country: '',
})

const initials = computed(() =>
  (form.value.full_name || auth.user?.full_name || '')
    .split(' ')
    .slice(0, 2)
    .map((n: string) => n[0])
    .join('')
    .toUpperCase() || '?',
)

const roleBadgeClass = computed(() =>
  auth.isAdmin ? 'bg-accent-red/20 text-accent-red' : 'bg-accent-orange/20 text-accent-orange',
)

async function loadProfile() {
  loading.value = true
  try {
    const res = await profileApi.get()
    const u = res.data.data
    if (u) {
      form.value = {
        full_name: u.full_name,
        phone: u.phone ?? '',
        age: u.age ?? null,
        gender: u.gender ?? '',
        country: u.country ?? '',
      }
    }
  } catch {
    // fall back to stored user data
    if (auth.user) {
      form.value = {
        full_name: auth.user.full_name,
        phone: auth.user.phone ?? '',
        age: auth.user.age ?? null,
        gender: auth.user.gender ?? '',
        country: auth.user.country ?? '',
      }
    }
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  saveError.value = null
  saveSuccess.value = false
  try {
    await profileApi.update({
      full_name: form.value.full_name,
      phone: form.value.phone,
      age: form.value.age ?? undefined,
      gender: form.value.gender,
      country: form.value.country,
    })
    // Update local user store so the header/avatar reflects the name change
    if (auth.user) {
      auth.user.full_name = form.value.full_name
      auth.user.phone = form.value.phone
      localStorage.setItem('ls_user', JSON.stringify(auth.user))
    }
    saveSuccess.value = true
    setTimeout(() => { saveSuccess.value = false }, 4000)
  } catch (e: any) {
    saveError.value = e.response?.data?.error ?? 'Failed to save profile'
  } finally {
    saving.value = false
  }
}

onMounted(loadProfile)
</script>
