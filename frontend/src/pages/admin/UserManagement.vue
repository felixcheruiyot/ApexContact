<template>
  <div class="space-y-6">

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <h1 class="text-white font-bold text-2xl">User Management</h1>
        <span v-if="!loading" class="text-text-muted text-sm">{{ filteredUsers.length }} shown</span>
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

    <!-- Filters -->
    <div class="flex flex-wrap gap-3">
      <div class="relative flex-1 min-w-48 max-w-sm">
        <svg class="w-4 h-4 text-text-muted absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="search"
          type="text"
          placeholder="Search by name or email…"
          class="input pl-9 text-sm"
        />
      </div>

      <select v-model="roleFilter" class="input text-sm w-auto pr-8">
        <option value="">All roles</option>
        <option value="viewer">Viewer</option>
        <option value="promoter">Promoter</option>
        <option value="broadcaster">Broadcaster</option>
        <option value="admin">Admin</option>
      </select>

      <select v-model="statusFilter" class="input text-sm w-auto pr-8">
        <option value="">All statuses</option>
        <option value="active">Active</option>
        <option value="locked">Locked</option>
      </select>
    </div>

    <!-- Skeleton -->
    <div v-if="loading" class="card overflow-hidden">
      <div v-for="i in 6" :key="i" class="px-6 py-4 border-b border-white/5 last:border-0 animate-pulse flex items-center gap-6">
        <div class="w-8 h-8 rounded-full bg-white/10 shrink-0" />
        <div class="flex-1 space-y-2">
          <div class="h-3.5 bg-white/10 rounded w-1/3" />
          <div class="h-3 bg-white/10 rounded w-1/4" />
        </div>
        <div class="h-5 bg-white/10 rounded-full w-16" />
        <div class="h-3 bg-white/10 rounded w-20" />
        <div class="h-3 bg-white/10 rounded w-20" />
        <div class="h-6 bg-white/10 rounded w-14" />
      </div>
    </div>

    <!-- Table -->
    <div v-else-if="filteredUsers.length" class="card overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-white/5">
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">User</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Role</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Phone</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Joined</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Status</th>
            <th class="px-6 py-3" />
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="user in filteredUsers" :key="user.id" class="hover:bg-bg-elevated transition-colors">
            <!-- Name + email -->
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full bg-accent-red/20 flex items-center justify-center shrink-0">
                  <span class="text-accent-red text-xs font-bold">{{ initials(user.full_name) }}</span>
                </div>
                <div class="min-w-0">
                  <p class="text-white text-sm font-medium truncate">{{ user.full_name }}</p>
                  <p class="text-text-muted text-xs truncate">{{ user.email }}</p>
                </div>
              </div>
            </td>

            <!-- Role badge -->
            <td class="px-6 py-4">
              <span class="text-xs font-semibold px-2.5 py-1 rounded-full capitalize" :class="roleBadgeClass(user.role)">
                {{ user.role }}
              </span>
            </td>

            <!-- Phone -->
            <td class="px-6 py-4 text-text-muted text-sm">{{ user.phone || '—' }}</td>

            <!-- Joined -->
            <td class="px-6 py-4 text-text-muted text-sm whitespace-nowrap">
              {{ formatDate(user.created_at) }}
            </td>

            <!-- Status -->
            <td class="px-6 py-4">
              <span
                class="text-xs font-semibold px-2.5 py-1 rounded-full"
                :class="user.is_locked
                  ? 'bg-status-error/15 text-status-error'
                  : 'bg-status-success/15 text-status-success'"
              >
                {{ user.is_locked ? 'Locked' : 'Active' }}
              </span>
            </td>

            <!-- Actions -->
            <td class="px-6 py-4 text-right whitespace-nowrap">
              <div class="flex items-center justify-end gap-2">
                <button
                  @click="openRoleEdit(user)"
                  class="text-xs text-text-muted hover:text-white hover:bg-white/10 px-3 py-1.5 rounded-lg transition-colors font-medium"
                >
                  Role
                </button>
                <button
                  v-if="!user.is_locked"
                  @click="openConfirm(user, 'lock')"
                  class="text-xs text-status-error hover:bg-status-error/15 px-3 py-1.5 rounded-lg transition-colors font-medium"
                >
                  Lock
                </button>
                <button
                  v-else
                  @click="openConfirm(user, 'unlock')"
                  class="text-xs text-status-success hover:bg-status-success/15 px-3 py-1.5 rounded-lg transition-colors font-medium"
                >
                  Unlock
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty state -->
    <div v-else class="card p-14 text-center">
      <p class="text-3xl mb-3">🔍</p>
      <p class="text-white font-semibold mb-1">No users found</p>
      <p class="text-text-muted text-sm">Try adjusting your search or filter.</p>
    </div>

    <!-- Role edit modal -->
    <Teleport to="body">
      <div v-if="roleTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="roleTarget = null" />
        <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-sm p-6 shadow-2xl">
          <h3 class="text-white font-bold text-lg mb-2">Change Role</h3>
          <div class="flex items-center gap-3 mb-5 p-3 bg-bg-surface rounded-xl">
            <div class="w-9 h-9 rounded-full bg-accent-red/20 flex items-center justify-center shrink-0">
              <span class="text-accent-red text-sm font-bold">{{ initials(roleTarget.full_name) }}</span>
            </div>
            <div class="min-w-0">
              <p class="text-white text-sm font-medium">{{ roleTarget.full_name }}</p>
              <p class="text-text-muted text-xs">{{ roleTarget.email }}</p>
            </div>
          </div>

          <label class="block text-text-muted text-xs uppercase tracking-wider mb-2">New Role</label>
          <select v-model="selectedRole" class="input w-full mb-5">
            <option value="viewer">Viewer</option>
            <option value="promoter">Promoter</option>
            <option value="broadcaster">Broadcaster</option>
            <option value="admin">Admin</option>
          </select>

          <div v-if="roleError" class="bg-status-error/10 border border-status-error/20 rounded-lg px-3 py-2 text-status-error text-xs mb-4">
            {{ roleError }}
          </div>

          <div class="flex gap-3">
            <button @click="roleTarget = null" class="btn-ghost flex-1 text-sm py-2.5">Cancel</button>
            <button
              @click="executeRoleChange"
              :disabled="acting"
              class="btn-primary flex-1 text-sm py-2.5 flex items-center justify-center gap-2"
            >
              <svg v-if="acting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
              </svg>
              {{ acting ? 'Saving…' : 'Save Role' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Confirm modal -->
    <Teleport to="body">
      <div v-if="confirmTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="confirmTarget = null" />
        <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-sm p-6 shadow-2xl">
          <h3 class="text-white font-bold text-lg mb-2">
            {{ confirmAction === 'lock' ? 'Lock this account?' : 'Unlock this account?' }}
          </h3>
          <div class="flex items-center gap-3 mb-4 p-3 bg-bg-surface rounded-xl">
            <div class="w-9 h-9 rounded-full bg-accent-red/20 flex items-center justify-center shrink-0">
              <span class="text-accent-red text-sm font-bold">{{ initials(confirmTarget.full_name) }}</span>
            </div>
            <div class="min-w-0">
              <p class="text-white text-sm font-medium">{{ confirmTarget.full_name }}</p>
              <p class="text-text-muted text-xs">{{ confirmTarget.email }}</p>
            </div>
          </div>
          <p class="text-text-muted text-sm mb-6">
            <template v-if="confirmAction === 'lock'">
              This will immediately revoke all active sessions and prevent this user from logging in.
            </template>
            <template v-else>
              This will restore the user's access. They will be able to log in and stream again.
            </template>
          </p>
          <div class="flex gap-3">
            <button @click="confirmTarget = null" class="btn-ghost flex-1 text-sm py-2.5">Cancel</button>
            <button
              @click="executeAction"
              :disabled="acting"
              class="flex-1 text-sm py-2.5 rounded-xl font-semibold transition-colors flex items-center justify-center gap-2"
              :class="confirmAction === 'lock'
                ? 'bg-status-error/20 text-status-error hover:bg-status-error/30'
                : 'bg-status-success/20 text-status-success hover:bg-status-success/30'"
            >
              <svg v-if="acting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
              </svg>
              {{ acting ? 'Please wait…' : (confirmAction === 'lock' ? 'Yes, Lock' : 'Yes, Unlock') }}
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
import { adminApi } from '@/api/admin'
import type { User, UserRole } from '@/types'

const users = ref<User[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const search = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const confirmTarget = ref<User | null>(null)
const confirmAction = ref<'lock' | 'unlock'>('lock')
const acting = ref(false)

// Role editing
const roleTarget = ref<User | null>(null)
const selectedRole = ref<UserRole>('viewer')
const roleError = ref<string | null>(null)

const filteredUsers = computed(() => {
  const q = search.value.toLowerCase()
  return users.value.filter((u) => {
    if (q && !u.full_name.toLowerCase().includes(q) && !u.email.toLowerCase().includes(q)) return false
    if (roleFilter.value && u.role !== roleFilter.value) return false
    if (statusFilter.value === 'active' && u.is_locked) return false
    if (statusFilter.value === 'locked' && !u.is_locked) return false
    return true
  })
})

function initials(name: string) {
  return name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase()
}

function formatDate(d: string) {
  return format(new Date(d), 'MMM d, yyyy')
}

function roleBadgeClass(role: UserRole) {
  return {
    viewer: 'bg-blue-500/15 text-blue-400',
    promoter: 'bg-accent-orange/15 text-accent-orange',
    broadcaster: 'bg-purple-500/15 text-purple-400',
    admin: 'bg-accent-red/15 text-accent-red',
  }[role]
}

function openRoleEdit(user: User) {
  roleTarget.value = user
  selectedRole.value = user.role
  roleError.value = null
}

async function executeRoleChange() {
  if (!roleTarget.value) return
  acting.value = true
  roleError.value = null
  try {
    await adminApi.updateUserRole(roleTarget.value.id, selectedRole.value)
    const u = users.value.find((x) => x.id === roleTarget.value!.id)
    if (u) u.role = selectedRole.value
    roleTarget.value = null
  } catch (e: any) {
    roleError.value = e.response?.data?.error ?? 'Failed to update role'
  } finally {
    acting.value = false
  }
}

function openConfirm(user: User, action: 'lock' | 'unlock') {
  confirmTarget.value = user
  confirmAction.value = action
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await adminApi.listUsers()
    users.value = res.data.data ?? []
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to load users'
  } finally {
    loading.value = false
  }
}

async function executeAction() {
  if (!confirmTarget.value) return
  acting.value = true
  const target = confirmTarget.value
  try {
    if (confirmAction.value === 'lock') {
      await adminApi.lockUser(target.id)
      const u = users.value.find((x) => x.id === target.id)
      if (u) u.is_locked = true
    } else {
      await adminApi.unlockUser(target.id)
      const u = users.value.find((x) => x.id === target.id)
      if (u) u.is_locked = false
    }
    confirmTarget.value = null
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Action failed'
    confirmTarget.value = null
  } finally {
    acting.value = false
  }
}

onMounted(load)
</script>
