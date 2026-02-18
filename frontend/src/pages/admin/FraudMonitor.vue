<template>
  <div class="space-y-6">
    <h1 class="text-white font-bold text-2xl">Fraud Monitor</h1>

    <div class="card overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-white/5">
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Reason</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">User ID</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Detected</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Status</th>
            <th class="px-6 py-4" />
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="flag in flags" :key="flag.id" class="hover:bg-bg-elevated transition-colors">
            <td class="px-6 py-4">
              <span class="text-status-error text-sm font-medium">{{ flag.reason }}</span>
            </td>
            <td class="px-6 py-4">
              <code class="text-text-muted text-xs">{{ flag.user_id.slice(0, 8) }}...</code>
            </td>
            <td class="px-6 py-4">
              <span class="text-text-muted text-sm">
                {{ format(new Date(flag.detected_at), 'MMM d, h:mm a') }}
              </span>
            </td>
            <td class="px-6 py-4">
              <span :class="flag.resolved ? 'text-status-success' : 'text-status-warning'" class="text-sm font-medium">
                {{ flag.resolved ? 'Resolved' : 'Open' }}
              </span>
            </td>
            <td class="px-6 py-4">
              <button v-if="!flag.resolved" @click="lockUser(flag.user_id)"
                class="text-status-error text-sm hover:underline">
                Lock User
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="!flags.length" class="text-center py-16">
        <p class="text-status-success font-semibold">No open fraud flags</p>
        <p class="text-text-muted text-sm mt-1">Everything looks clean.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { format } from 'date-fns'
import client from '@/api/client'
import type { FraudFlag } from '@/types'

const flags = ref<FraudFlag[]>([])

onMounted(async () => {
  const res = await client.get('/admin/fraud')
  flags.value = res.data.data ?? []
})

async function lockUser(userId: string) {
  await client.post(`/admin/users/${userId}/lock`)
  flags.value = flags.value.filter((f) => f.user_id !== userId)
}
</script>
