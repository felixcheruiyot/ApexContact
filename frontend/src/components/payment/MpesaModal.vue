<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="handleClose" />

      <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-md
                  shadow-2xl animate-slide-up">
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-5 border-b border-white/5">
          <div>
            <h2 class="text-white font-bold text-lg">Buy Ticket</h2>
            <p class="text-text-muted text-sm">{{ event.title }}</p>
          </div>
          <button @click="handleClose" class="text-text-muted hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Content -->
        <div class="px-6 py-6">

          <!-- ── Idle: phone input ─────────────────────────────────────── -->
          <template v-if="payment.status === 'idle'">
            <div class="bg-bg-surface rounded-xl p-4 mb-6 flex items-center justify-between">
              <div>
                <p class="text-text-muted text-xs mb-1">Total Amount</p>
                <p class="text-white font-bold text-2xl">
                  {{ event.currency }} {{ event.price.toLocaleString() }}
                </p>
              </div>
              <span class="font-bold text-sm tracking-widest px-2 py-1 rounded bg-[#00A651]/20 text-[#00A651]">
                M-PESA
              </span>
            </div>

            <div class="mb-5">
              <label class="block text-text-muted text-sm mb-2">M-Pesa Phone Number</label>
              <input
                v-model="phoneNumber"
                type="tel"
                placeholder="254712345678"
                class="input"
                @keyup.enter="phoneNumber.trim() && handlePay()"
              />
              <p class="text-text-muted text-xs mt-1.5">
                Enter your number in international format (e.g. 254712345678)
              </p>
            </div>

            <button
              @click="handlePay"
              class="btn-primary w-full"
              :disabled="!phoneNumber.trim()"
            >
              Pay with M-Pesa
            </button>
          </template>

          <!-- ── Pending / Polling: waiting for PIN ───────────────────── -->
          <template v-else-if="payment.status === 'pending' || payment.status === 'polling'">
            <div class="text-center py-8">
              <div class="w-16 h-16 border-4 border-accent-red border-t-transparent rounded-full
                          animate-spin mx-auto mb-4" />
              <h3 class="text-white font-semibold mb-2">Check your phone</h3>
              <p class="text-text-muted text-sm">
                Enter your M-Pesa PIN to complete the payment of
                <span class="text-white font-bold">
                  {{ event.currency }} {{ event.price.toLocaleString() }}
                </span>
              </p>
              <p class="text-text-muted text-xs mt-4">Prompt expires in 5 minutes</p>
            </div>
          </template>

          <!-- ── Payment success: activating ticket ───────────────────── -->
          <template v-else-if="payment.status === 'success'">
            <div class="text-center py-8">
              <template v-if="activating">
                <div class="w-16 h-16 border-4 border-status-success border-t-transparent rounded-full
                            animate-spin mx-auto mb-4" />
                <h3 class="text-white font-semibold mb-2">Activating your ticket…</h3>
                <p class="text-text-muted text-sm">Just a moment</p>
              </template>
              <template v-else-if="activationError">
                <div class="w-16 h-16 rounded-full bg-status-error/20 flex items-center justify-center mx-auto mb-4">
                  <svg class="w-8 h-8 text-status-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <h3 class="text-white font-bold text-lg mb-2">Payment received!</h3>
                <p class="text-text-muted text-sm mb-2">
                  Your payment was confirmed, but we couldn't activate your ticket automatically.
                </p>
                <p class="text-text-muted text-xs mb-6">
                  Please contact support with reference: <span class="text-white">{{ payment.paymentId }}</span>
                </p>
                <button @click="handleClose" class="btn-primary w-full">Close</button>
              </template>
              <template v-else>
                <div class="w-16 h-16 rounded-full bg-status-success/20 flex items-center justify-center
                            mx-auto mb-4">
                  <svg class="w-8 h-8 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <h3 class="text-white font-bold text-lg mb-2">Payment Successful!</h3>
                <p class="text-text-muted text-sm mb-6">Your ticket is confirmed. Enjoy the event!</p>
                <button @click="handleWatchNow" class="btn-primary w-full">Watch Now</button>
              </template>
            </div>
          </template>

          <!-- ── Failed / Cancelled / Timeout ────────────────────────── -->
          <template v-else-if="['failed', 'cancelled', 'timeout'].includes(payment.status)">
            <div class="text-center py-8">
              <div class="w-16 h-16 rounded-full bg-status-error/20 flex items-center justify-center
                          mx-auto mb-4">
                <svg class="w-8 h-8 text-status-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </div>
              <h3 class="text-white font-bold text-lg mb-2">
                {{ payment.status === 'timeout' ? 'Payment Timed Out' : 'Payment Failed' }}
              </h3>
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
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePaymentStore } from '@/stores/payment'
import type { Event } from '@/types'

const props = defineProps<{ event: Event }>()
const emit = defineEmits<{ close: [] }>()

const payment = usePaymentStore()
const router = useRouter()
const phoneNumber = ref('')
const activating = ref(false)
const activationError = ref(false)

// Auto-activate subscription as soon as payment is confirmed
watch(
  () => payment.status,
  async (status) => {
    if (status === 'success' && !payment.streamToken) {
      activating.value = true
      activationError.value = false
      try {
        const fingerprint = await generateFingerprint()
        await payment.createSubscription(props.event.id, fingerprint)
      } catch {
        activationError.value = true
      } finally {
        activating.value = false
      }
    }
  },
)

async function handlePay() {
  if (!phoneNumber.value.trim()) return
  await payment.initiate(props.event.id, phoneNumber.value.trim())
}

function handleWatchNow() {
  if (!payment.streamToken) return
  emit('close')
  router.push({
    name: 'watch',
    params: { eventId: props.event.id },
    query: { token: payment.streamToken },
  })
}

function handleClose() {
  // Only allow closing when not mid-flow (let user cancel from idle/failed/success)
  if (payment.status === 'pending' || payment.status === 'polling') return
  payment.reset()
  emit('close')
}

async function generateFingerprint(): Promise<string> {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  if (ctx) {
    ctx.textBaseline = 'top'
    ctx.font = '14px Arial'
    ctx.fillText('ApexContact', 2, 2)
  }
  const raw = [
    canvas.toDataURL(),
    `${screen.width}x${screen.height}x${screen.colorDepth}`,
    Intl.DateTimeFormat().resolvedOptions().timeZone,
    navigator.language,
    navigator.hardwareConcurrency,
  ].join('|')

  const buf = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(raw))
  return Array.from(new Uint8Array(buf))
    .map((b) => b.toString(16).padStart(2, '0'))
    .join('')
}
</script>
