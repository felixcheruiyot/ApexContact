<template>
  <div class="space-y-6">

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <h1 class="text-white font-bold text-2xl">Fraud Monitor</h1>
        <span v-if="openCount > 0"
          class="bg-status-warning/20 text-status-warning text-xs font-bold px-2.5 py-1 rounded-full">
          {{ openCount }} open
        </span>
      </div>
      <button @click="load" class="btn-ghost text-sm py-2 px-4 flex items-center gap-2">
        <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh
      </button>
    </div>

    <!-- Error -->
    <div v-if="error" class="bg-status-error/10 border border-status-error/20 rounded-xl px-5 py-4 text-status-error text-sm">
      {{ error }}
    </div>

    <!-- Filter tabs -->
    <div class="flex gap-1 bg-bg-surface rounded-xl p-1 w-fit">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="activeTab = tab.value"
        class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
        :class="activeTab === tab.value
          ? 'bg-bg-elevated text-white'
          : 'text-text-muted hover:text-white'"
      >
        {{ tab.label }}
        <span v-if="tab.value === 'open' && openCount > 0"
          class="ml-1.5 bg-status-warning/30 text-status-warning text-xs rounded-full px-1.5 py-0.5">
          {{ openCount }}
        </span>
      </button>
    </div>

    <!-- Skeleton -->
    <div v-if="loading" class="card overflow-hidden">
      <div v-for="i in 4" :key="i" class="px-6 py-4 border-b border-white/5 last:border-0 animate-pulse flex gap-6">
        <div class="h-4 bg-white/10 rounded w-1/4" />
        <div class="h-4 bg-white/10 rounded w-1/5" />
        <div class="h-4 bg-white/10 rounded w-1/6" />
        <div class="h-4 bg-white/10 rounded w-1/6" />
      </div>
    </div>

    <!-- Table -->
    <div v-else-if="visibleFlags.length" class="card overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-white/5">
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Reason</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">User ID</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Detected</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Status</th>
            <th class="px-6 py-3" />
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="flag in visibleFlags" :key="flag.id" class="hover:bg-bg-elevated transition-colors">
            <td class="px-6 py-4">
              <span class="text-status-warning text-sm font-medium">{{ flag.reason }}</span>
            </td>
            <td class="px-6 py-4">
              <code class="text-text-muted text-xs bg-bg-surface px-2 py-0.5 rounded">
                {{ flag.user_id.slice(0, 8) }}…
              </code>
            </td>
            <td class="px-6 py-4 text-text-muted text-sm">
              {{ formatDate(flag.detected_at) }}
            </td>
            <td class="px-6 py-4">
              <span
                class="text-xs font-semibold px-2 py-0.5 rounded-full"
                :class="flag.resolved
                  ? 'bg-status-success/15 text-status-success'
                  : 'bg-status-warning/15 text-status-warning'"
              >
                {{ flag.resolved ? 'Resolved' : 'Open' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <button
                v-if="!flag.resolved"
                @click="confirmLock(flag)"
                class="text-xs text-status-error hover:text-white hover:bg-status-error/20
                       px-3 py-1.5 rounded-lg transition-colors font-medium"
              >
                Lock User
              </button>
              <span v-else class="text-text-muted text-xs">—</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty state -->
    <div v-else class="card p-14 text-center">
      <template v-if="activeTab === 'open'">
        <div class="w-14 h-14 rounded-full bg-status-success/10 flex items-center justify-center mx-auto mb-4">
          <svg class="w-7 h-7 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <p class="text-white font-semibold mb-1">All clear</p>
        <p class="text-text-muted text-sm">No open fraud flags at this time.</p>
      </template>
      <template v-else>
        <FolderOpen class="w-10 h-10 mx-auto mb-3 text-text-muted" />
        <p class="text-text-muted text-sm">No resolved flags to display.</p>
      </template>
    </div>

    <!-- Confirm Lock Modal -->
    <Teleport to="body">
      <div v-if="lockTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="lockTarget = null" />
        <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-sm p-6 shadow-2xl">
          <h3 class="text-white font-bold text-lg mb-2">Lock this user?</h3>
          <p class="text-text-muted text-sm mb-1">
            User ID: <code class="text-white bg-bg-surface px-2 py-0.5 rounded text-xs">{{ lockTarget.user_id.slice(0, 8) }}…</code>
          </p>
          <p class="text-text-muted text-sm mb-6">
            Reason: <span class="text-status-warning">{{ lockTarget.reason }}</span>
          </p>
          <p class="text-text-muted text-xs mb-6">
            This will immediately revoke all active sessions for this user. They will not be able to log in
            or stream until unlocked from User Management.
          </p>
          <div class="flex gap-3">
            <button @click="lockTarget = null" class="btn-ghost flex-1 text-sm py-2.5">Cancel</button>
            <button
              @click="executeLock"
              :disabled="locking"
              class="btn-primary flex-1 text-sm py-2.5 flex items-center justify-center gap-2"
            >
              <svg v-if="locking" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
              </svg>
              {{ locking ? 'Locking…' : 'Yes, Lock User' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { format } from 'date-fns'
import { FolderOpen } from 'lucide-vue-next'
import { adminApi } from '@/api/admin'
import type { FraudFlag } from '@/types'

const flags = ref<FraudFlag[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const activeTab = ref<'open' | 'resolved'>('open')
const lockTarget = ref<FraudFlag | null>(null)
const locking = ref(false)

const tabs = [
  { label: 'Open', value: 'open' as const },
  { label: 'History', value: 'resolved' as const },
]

const openCount = computed(() => flags.value.filter((f) => !f.resolved).length)

const visibleFlags = computed(() =>
  activeTab.value === 'open'
    ? flags.value.filter((f) => !f.resolved)
    : flags.value.filter((f) => f.resolved),
)

function formatDate(d: string) {
  return format(new Date(d), 'MMM d, h:mm a')
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await adminApi.listFraudFlags()
    flags.value = res.data.data ?? []
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to load fraud flags'
  } finally {
    loading.value = false
  }
}

function confirmLock(flag: FraudFlag) {
  lockTarget.value = flag
}

async function executeLock() {
  if (!lockTarget.value) return
  locking.value = true
  try {
    await adminApi.lockUser(lockTarget.value.user_id)
    // Mark matching flags as resolved locally so UI updates instantly
    flags.value = flags.value.map((f) =>
      f.user_id === lockTarget.value!.user_id ? { ...f, resolved: true } : f,
    )
    lockTarget.value = null
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to lock user'
    lockTarget.value = null
  } finally {
    locking.value = false
  }
}

onMounted(load)
</script>
