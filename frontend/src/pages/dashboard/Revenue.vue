<template>
  <div class="space-y-6">

    <!-- Header -->
    <div>
      <h1 class="text-2xl font-display text-white tracking-wide">Revenue & Analytics</h1>
      <p class="text-text-muted text-sm mt-1">All-event breakdown — tickets sold, revenue, and your 70% promoter share.</p>
    </div>

    <!-- Summary cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-bg-surface rounded-xl border border-white/5 p-4">
        <p class="text-text-muted text-xs uppercase tracking-widest mb-1">Total Events</p>
        <p class="text-2xl font-bold text-white">{{ rows.length }}</p>
      </div>
      <div class="bg-bg-surface rounded-xl border border-white/5 p-4">
        <p class="text-text-muted text-xs uppercase tracking-widest mb-1">Tickets Sold</p>
        <p class="text-2xl font-bold text-white">{{ totalTickets.toLocaleString() }}</p>
      </div>
      <div class="bg-bg-surface rounded-xl border border-white/5 p-4">
        <p class="text-text-muted text-xs uppercase tracking-widest mb-1">Total Revenue</p>
        <p class="text-2xl font-bold text-accent-orange">{{ formatCurrency(totalRevenue) }}</p>
      </div>
      <div class="bg-bg-surface rounded-xl border border-white/5 p-4">
        <p class="text-text-muted text-xs uppercase tracking-widest mb-1">Your Earnings (70%)</p>
        <p class="text-2xl font-bold text-green-400">{{ formatCurrency(totalPromoterCut) }}</p>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-bg-surface rounded-xl border border-white/5 overflow-hidden">
      <div class="px-5 py-4 border-b border-white/5 flex items-center justify-between">
        <h2 class="text-sm font-semibold text-white">Event Revenue Breakdown</h2>
        <div class="flex items-center gap-2">
          <button
            v-for="f in statusFilters"
            :key="f.value"
            @click="activeFilter = f.value"
            class="px-3 py-1 rounded-full text-xs font-medium transition-colors"
            :class="activeFilter === f.value
              ? 'bg-accent-red text-white'
              : 'text-text-muted hover:text-white hover:bg-white/5'"
          >
            {{ f.label }}
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <svg class="animate-spin w-6 h-6 text-accent-red" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
        </svg>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-16 text-red-400 text-sm">{{ error }}</div>

      <!-- Empty -->
      <div v-else-if="filteredRows.length === 0" class="text-center py-16 text-text-muted text-sm">
        No events match the selected filter.
      </div>

      <!-- Table body -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-white/5">
              <th class="text-left px-5 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Event</th>
              <th class="text-left px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Category</th>
              <th class="text-left px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Date</th>
              <th class="text-left px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Status</th>
              <th class="text-right px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Ticket Price</th>
              <th class="text-right px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Tickets Sold</th>
              <th class="text-right px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Gross Revenue</th>
              <th class="text-right px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Platform (30%)</th>
              <th class="text-right px-4 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Your Share (70%)</th>
              <th class="text-right px-5 py-3 text-text-muted text-xs uppercase tracking-widest font-medium">Peak Viewers</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in filteredRows"
              :key="row.event_id"
              class="border-b border-white/5 hover:bg-white/[0.02] transition-colors"
            >
              <td class="px-5 py-3.5">
                <RouterLink
                  :to="`/dashboard/analytics/${row.event_id}`"
                  class="text-white font-medium hover:text-accent-orange transition-colors line-clamp-1"
                >
                  {{ row.title }}
                </RouterLink>
              </td>
              <td class="px-4 py-3.5">
                <span class="inline-flex items-center gap-1.5 text-xs">
                  <span>{{ categoryIcon(row.sport_type) }}</span>
                  <span class="text-text-muted capitalize">{{ row.sport_type }}</span>
                </span>
              </td>
              <td class="px-4 py-3.5 text-text-muted whitespace-nowrap">
                {{ formatDate(row.scheduled_at) }}
              </td>
              <td class="px-4 py-3.5">
                <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium" :class="statusClass(row.status)">
                  <span class="w-1.5 h-1.5 rounded-full" :class="statusDot(row.status)" />
                  {{ row.status }}
                </span>
              </td>
              <td class="px-4 py-3.5 text-right text-text-muted whitespace-nowrap">
                {{ formatAmount(row.price, row.currency) }}
              </td>
              <td class="px-4 py-3.5 text-right text-white font-semibold">
                {{ row.tickets.toLocaleString() }}
              </td>
              <td class="px-4 py-3.5 text-right text-white whitespace-nowrap">
                {{ formatAmount(row.revenue, row.currency) }}
              </td>
              <td class="px-4 py-3.5 text-right text-text-muted whitespace-nowrap">
                {{ formatAmount(row.revenue * 0.30, row.currency) }}
              </td>
              <td class="px-4 py-3.5 text-right text-green-400 font-semibold whitespace-nowrap">
                {{ formatAmount(row.promoter_cut, row.currency) }}
              </td>
              <td class="px-5 py-3.5 text-right text-text-muted">
                {{ row.peak_viewers.toLocaleString() }}
              </td>
            </tr>
          </tbody>
          <tfoot v-if="filteredRows.length > 1">
            <tr class="bg-white/[0.02] border-t border-white/10">
              <td colspan="5" class="px-5 py-3 text-text-muted text-xs uppercase tracking-widest font-semibold">Totals</td>
              <td class="px-4 py-3 text-right text-white font-bold">{{ filteredTickets.toLocaleString() }}</td>
              <td class="px-4 py-3 text-right text-white font-bold whitespace-nowrap">KES {{ filteredRevenue.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</td>
              <td class="px-4 py-3 text-right text-text-muted font-bold whitespace-nowrap">KES {{ (filteredRevenue * 0.30).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</td>
              <td class="px-4 py-3 text-right text-green-400 font-bold whitespace-nowrap">KES {{ (filteredRevenue * 0.70).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</td>
              <td class="px-5 py-3" />
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Revenue split note -->
    <p class="text-xs text-text-muted text-center">
      Revenue split: <span class="text-green-400 font-medium">70% Promoter</span> / <span class="text-text-muted">30% Platform</span>.
      Figures reflect confirmed M-Pesa payments only.
    </p>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { eventsApi } from '@/api/events'

interface RevenueRow {
  event_id: string
  title: string
  sport_type: string
  status: string
  scheduled_at: string
  price: number
  currency: string
  tickets: number
  revenue: number
  promoter_cut: number
  peak_viewers: number
}

const rows = ref<RevenueRow[]>([])
const loading = ref(true)
const error = ref('')
const activeFilter = ref('all')

const statusFilters = [
  { label: 'All', value: 'all' },
  { label: 'Completed', value: 'completed' },
  { label: 'Live', value: 'live' },
  { label: 'Scheduled', value: 'scheduled' },
]

onMounted(async () => {
  try {
    const res = await eventsApi.revenue()
    rows.value = res.data.data ?? []
  } catch {
    error.value = 'Failed to load revenue data.'
  } finally {
    loading.value = false
  }
})

const filteredRows = computed(() =>
  activeFilter.value === 'all'
    ? rows.value
    : rows.value.filter(r => r.status === activeFilter.value),
)

const totalTickets = computed(() => rows.value.reduce((s, r) => s + r.tickets, 0))
const totalRevenue = computed(() => rows.value.reduce((s, r) => s + r.revenue, 0))
const totalPromoterCut = computed(() => rows.value.reduce((s, r) => s + r.promoter_cut, 0))

const filteredTickets = computed(() => filteredRows.value.reduce((s, r) => s + r.tickets, 0))
const filteredRevenue = computed(() => filteredRows.value.reduce((s, r) => s + r.revenue, 0))

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('en-KE', {
    day: 'numeric', month: 'short', year: 'numeric',
  })
}

function formatCurrency(amount: number) {
  return `KES ${amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

function formatAmount(amount: number, currency = 'KES') {
  return `${currency} ${amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

function statusClass(status: string) {
  switch (status) {
    case 'live': return 'bg-green-500/10 border border-green-500/20 text-green-400'
    case 'completed': return 'bg-white/5 border border-white/10 text-text-muted'
    case 'cancelled': return 'bg-red-500/10 border border-red-500/20 text-red-400'
    default: return 'bg-accent-orange/10 border border-accent-orange/20 text-accent-orange'
  }
}

function categoryIcon(type: string) {
  const icons: Record<string, string> = {
    sales: '🎯', mentoring: '🏫', business: '💼', education: '📚',
    visa: '🌍', legal: '⚖️', fitness: '💪', music: '🎵',
    gaming: '🎮', cooking: '🍳', community: '🙏', other: '📌',
  }
  return icons[type] ?? '📌'
}

function statusDot(status: string) {
  switch (status) {
    case 'live': return 'bg-green-400 animate-pulse'
    case 'completed': return 'bg-text-muted'
    case 'cancelled': return 'bg-red-400'
    default: return 'bg-accent-orange'
  }
}
</script>
