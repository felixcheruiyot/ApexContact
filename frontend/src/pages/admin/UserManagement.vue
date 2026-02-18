<template>
  <div class="space-y-6">
    <h1 class="text-white font-bold text-2xl">User Management</h1>

    <div class="card overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-white/5">
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">User</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Role</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Phone</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Joined</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Status</th>
            <th class="px-6 py-4" />
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="user in users" :key="user.id" class="hover:bg-bg-elevated transition-colors">
            <td class="px-6 py-4">
              <p class="text-white font-medium text-sm">{{ user.full_name }}</p>
              <p class="text-text-muted text-xs">{{ user.email }}</p>
            </td>
            <td class="px-6 py-4">
              <span class="badge-upcoming text-xs capitalize">{{ user.role }}</span>
            </td>
            <td class="px-6 py-4">
              <span class="text-text-muted text-sm">{{ user.phone || '—' }}</span>
            </td>
            <td class="px-6 py-4">
              <span class="text-text-muted text-sm">
                {{ format(new Date(user.created_at), 'MMM d, yyyy') }}
              </span>
            </td>
            <td class="px-6 py-4">
              <span :class="user.is_locked ? 'text-status-error' : 'text-status-success'" class="text-sm font-medium">
                {{ user.is_locked ? 'Locked' : 'Active' }}
              </span>
            </td>
            <td class="px-6 py-4">
              <button v-if="!user.is_locked" @click="lockUser(user)"
                class="text-status-error text-sm hover:underline">
                Lock
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { format } from 'date-fns'
import client from '@/api/client'
import type { User } from '@/types'

const users = ref<User[]>([])

onMounted(async () => {
  const res = await client.get('/admin/users')
  users.value = res.data.data ?? []
})

async function lockUser(user: User) {
  await client.post(`/admin/users/${user.id}/lock`)
  user.is_locked = true
}
</script>
