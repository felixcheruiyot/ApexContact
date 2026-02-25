<template>
  <div>
    <div v-if="eventsStore.loading" class="flex items-center justify-center min-h-[60vh]">
      <div class="w-10 h-10 border-4 border-accent-red border-t-transparent rounded-full animate-spin" />
    </div>

    <template v-else-if="event">
      <!-- Waiting for host banner: event time has passed but broadcaster hasn't gone live yet -->
      <div
        v-if="event.status === 'live' && event.event_type === 'video' && !event.stream_active"
        class="bg-accent-orange/10 border-b border-accent-orange/20 px-4 py-3"
      >
        <div class="max-w-7xl mx-auto flex items-center gap-3">
          <div class="w-2 h-2 rounded-full bg-accent-orange animate-pulse shrink-0" />
          <p class="text-accent-orange text-sm font-medium">
            The event has started — waiting for the host to begin the stream. Please stay on this page and refresh shortly.
          </p>
        </div>
      </div>

      <EventHero :event="event" @buy-ticket="handleCta" />

      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-12">
          <!-- Description -->
          <div class="lg:col-span-2">
            <h2 class="section-heading mb-6">About This Event</h2>
            <p class="text-text-muted leading-relaxed whitespace-pre-wrap">{{ event.description }}</p>
          </div>

          <!-- Sidebar info -->
          <div class="space-y-4">
            <div class="card p-5">
              <h3 class="text-white font-semibold mb-4">Event Details</h3>
              <dl class="space-y-3">
                <div>
                  <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Date & Time</dt>
                  <dd class="text-white text-sm">{{ formattedDate }}</dd>
                </div>
                <div>
                  <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Sport</dt>
                  <dd class="text-white text-sm capitalize">{{ event.sport_type }}</dd>
                </div>
                <div>
                  <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Status</dt>
                  <dd>
                    <span v-if="event.status === 'live'" class="badge-live text-xs">Live</span>
                    <span v-else class="badge-upcoming text-xs capitalize">{{ event.status }}</span>
                  </dd>
                </div>
                <div v-if="!auth.isAdmin">
                  <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Ticket Price</dt>
                  <dd class="text-accent-red font-bold text-xl">
                    {{ event.currency }} {{ event.price.toLocaleString() }}
                  </dd>
                </div>
                <div v-else>
                  <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Access</dt>
                  <dd class="text-status-success font-semibold text-sm">Admin — Free</dd>
                </div>
              </dl>

              <!-- CTA button — changes based on access state -->
              <button @click="handleCta" class="btn-primary w-full mt-6">
                {{ ctaLabel }}
              </button>
            </div>

            <!-- Countdown -->
            <div v-if="event.status === 'scheduled'" class="card p-5">
              <h3 class="text-white font-semibold mb-4">Event Starts In</h3>
              <CountdownTimer :target="event.scheduled_at" />
            </div>
          </div>
        </div>
      </div>
    </template>

    <MpesaModal
      v-if="event && showPaymentModal"
      :event="event"
      @close="showPaymentModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { format } from 'date-fns'
import { useEventsStore } from '@/stores/events'
import { useAuthStore } from '@/stores/auth'
import { usePaymentStore } from '@/stores/payment'
import EventHero from '@/components/events/EventHero.vue'
import CountdownTimer from '@/components/events/CountdownTimer.vue'
import MpesaModal from '@/components/payment/MpesaModal.vue'

const route = useRoute()
const router = useRouter()
const eventsStore = useEventsStore()
const auth = useAuthStore()
const payment = usePaymentStore()
const showPaymentModal = ref(false)

const event = computed(() => eventsStore.currentEvent)
const formattedDate = computed(() =>
  event.value ? format(new Date(event.value.scheduled_at), 'EEEE, MMMM d, yyyy · h:mm a') : ''
)

// Does the authenticated viewer already own a ticket for this event?
const existingToken = computed(() =>
  event.value ? payment.getStoredToken(event.value.id) : null
)

const waitingForHost = computed(() =>
  event.value?.status === 'live' &&
  event.value?.event_type === 'video' &&
  !event.value?.stream_active
)

const ctaLabel = computed(() => {
  if (auth.isAdmin) return event.value?.status === 'live' ? 'Watch Free (Admin)' : 'Event Details (Admin)'
  if (waitingForHost.value) return existingToken.value ? 'Stream Starting Soon' : 'Reserve Your Spot'
  if (existingToken.value) return event.value?.status === 'live' ? 'Watch Now (Live)' : 'Watch Now'
  return event.value?.status === 'live' ? 'Buy Ticket — Event is Live Now' : 'Buy Ticket'
})

onMounted(() => eventsStore.fetchEvent(route.params.id as string))

function handleCta() {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: route.fullPath } })
    return
  }

  // Admins go straight to Watch — no token needed
  if (auth.isAdmin) {
    router.push({ name: 'watch', params: { eventId: event.value!.id } })
    return
  }

  // Viewer with existing subscription — go straight to Watch with their token
  if (existingToken.value) {
    router.push({
      name: 'watch',
      params: { eventId: event.value!.id },
      query: { token: existingToken.value },
    })
    return
  }

  // New purchase
  showPaymentModal.value = true
}
</script>
