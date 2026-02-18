<template>
  <!-- Backdrop -->
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="$emit('close')" />

      <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-md
                  shadow-2xl animate-slide-up">
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-5 border-b border-white/5">
          <div>
            <h2 class="text-white font-bold text-lg">Buy Ticket</h2>
            <p class="text-text-muted text-sm">{{ event.title }}</p>
          </div>
          <button @click="$emit('close')" class="text-text-muted hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Content -->
        <div class="px-6 py-6">
          <!-- Idle / Input state -->
          <template v-if="payment.status === 'idle'">
            <div class="bg-bg-surface rounded-xl p-4 mb-6 flex items-center justify-between">
              <div>
                <p class="text-text-muted text-xs mb-1">Total Amount</p>
                <p class="text-white font-bold text-2xl">{{ event.currency }} {{ event.price.toLocaleString() }}</p>
              </div>
              <span class="font-bold text-sm tracking-widest px-2 py-1 rounded bg-[#00A651]/20 text-[#00A651]">M-PESA</span>
            </div>

            <div class="mb-5">
              <label class="block text-text-muted text-sm mb-2">M-Pesa Phone Number</label>
              <input
                v-model="phoneNumber"
                type="tel"
                placeholder="+254712345678"
                class="input"
              />
              <p class="text-text-muted text-xs mt-1.5">
                You'll receive an STK push prompt on this number
              </p>
            </div>

            <button @click="handlePay" class="btn-primary w-full" :disabled="!phoneNumber.trim()">
              Pay with M-Pesa
            </button>
          </template>

          <!-- Polling / Pending state -->
          <template v-else-if="payment.status === 'pending' || payment.status === 'polling'">
            <div class="text-center py-8">
              <div class="w-16 h-16 border-4 border-accent-red border-t-transparent rounded-full
                          animate-spin mx-auto mb-4" />
              <h3 class="text-white font-semibold mb-2">Check your phone</h3>
              <p class="text-text-muted text-sm">
                Enter your M-Pesa PIN to complete the payment of
                <span class="text-white font-bold">{{ event.currency }} {{ event.price.toLocaleString() }}</span>
              </p>
            </div>
          </template>

          <!-- Success state -->
          <template v-else-if="payment.status === 'success'">
            <div class="text-center py-8">
              <div class="w-16 h-16 rounded-full bg-status-success/20 flex items-center justify-center
                          mx-auto mb-4">
                <svg class="w-8 h-8 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h3 class="text-white font-bold text-lg mb-2">Payment Successful!</h3>
              <p class="text-text-muted text-sm mb-6">Your ticket is confirmed. Enjoy the event!</p>
              <button @click="handleWatchNow" class="btn-primary w-full">Watch Now</button>
            </div>
          </template>

          <!-- Failed state -->
          <template v-else-if="payment.status === 'failed' || payment.status === 'cancelled'">
            <div class="text-center py-8">
              <div class="w-16 h-16 rounded-full bg-status-error/20 flex items-center justify-center
                          mx-auto mb-4">
                <svg class="w-8 h-8 text-status-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </div>
              <h3 class="text-white font-bold text-lg mb-2">Payment Failed</h3>
              <p class="text-text-muted text-sm mb-6">
                {{ payment.error || 'The payment was not completed. Please try again.' }}
              </p>
              <button @click="payment.reset()" class="btn-primary w-full">Try Again</button>
            </div>
          </template>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePaymentStore } from '@/stores/payment'
import type { Event } from '@/types'

const props = defineProps<{ event: Event }>()
const emit = defineEmits<{ close: [] }>()

const payment = usePaymentStore()
const router = useRouter()
const phoneNumber = ref('')

async function handlePay() {
  await payment.initiate(props.event.id, phoneNumber.value)
}

async function handleWatchNow() {
  const fingerprint = await generateFingerprint()
  await payment.createSubscription(props.event.id, fingerprint)
  if (payment.streamToken) {
    emit('close')
    router.push({ name: 'watch', params: { eventId: props.event.id }, query: { token: payment.streamToken } })
  }
}

// Basic device fingerprint using canvas and screen properties
async function generateFingerprint(): Promise<string> {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  if (ctx) {
    ctx.textBaseline = 'top'
    ctx.font = '14px Arial'
    ctx.fillText('LiveStreamify', 2, 2)
  }
  const canvasData = canvas.toDataURL()
  const screenData = `${screen.width}x${screen.height}x${screen.colorDepth}`
  const tzData = Intl.DateTimeFormat().resolvedOptions().timeZone
  const raw = `${canvasData}|${screenData}|${tzData}|${navigator.language}`

  // Simple hash using SubtleCrypto
  const buf = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(raw))
  return Array.from(new Uint8Array(buf)).map(b => b.toString(16).padStart(2, '0')).join('')
}
</script>
