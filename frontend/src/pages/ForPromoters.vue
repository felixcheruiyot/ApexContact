<template>
  <div class="bg-bg text-white">

    <!-- ============================================================ -->
    <!-- SECTION 1: HERO                                               -->
    <!-- ============================================================ -->
    <section class="relative py-24 lg:py-32 overflow-hidden">
      <!-- Subtle background glow -->
      <div class="absolute inset-0 pointer-events-none">
        <div
          class="absolute top-1/2 left-1/4 w-[500px] h-[500px] rounded-full opacity-8"
          style="background: radial-gradient(circle, #E8002D 0%, transparent 70%); filter: blur(120px);"
        />
      </div>

      <div class="relative max-w-5xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="max-w-2xl space-y-8">
          <!-- Eyebrow -->
          <div class="inline-flex items-center gap-2 bg-accent-red/10 border border-accent-red/20
                      text-accent-red text-xs font-bold uppercase tracking-widest px-3 py-1.5 rounded">
            For Creators
          </div>

          <h1 class="font-display text-6xl md:text-7xl uppercase text-white leading-none tracking-wide">
            Turn Your Expertise<br />
            Into Income.<br />
            <span class="text-accent-red">Live.</span>
          </h1>

          <p class="text-text-muted text-lg leading-relaxed">
            Host live sessions on any topic. Set a ticket price.
            Collect payments via M-Pesa. 70% goes directly to you.
          </p>

          <div class="flex flex-wrap items-center gap-4 pt-2">
            <RouterLink
              to="/try"
              class="inline-flex items-center gap-2 px-8 py-4 rounded-lg bg-accent-red
                     hover:bg-accent-red-hover text-white font-bold text-base
                     transition-all duration-200 active:scale-95"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              Try a Free Test Stream
            </RouterLink>
            <a
              href="#how-it-works"
              class="inline-flex items-center gap-2 px-6 py-4 rounded-lg border border-white/20
                     hover:border-white/40 text-white text-base font-medium
                     transition-all duration-200 hover:bg-white/5"
            >
              See How Payouts Work
            </a>
          </div>

          <p class="text-text-muted text-sm">
            No account needed to test &nbsp;·&nbsp; M-Pesa payouts &nbsp;·&nbsp; 70% revenue share
          </p>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 2: REVENUE CALCULATOR                                 -->
    <!-- ============================================================ -->
    <section class="py-24 bg-bg-surface border-y border-white/5">
      <div class="max-w-5xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="mb-14">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Earnings</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Calculate your earnings
          </h2>
          <p class="text-text-muted text-base mt-4">
            70% of every ticket goes straight to you
          </p>
        </div>

        <div class="bg-bg-elevated rounded-xl p-8 md:p-10 max-w-2xl border border-white/5">
          <div class="space-y-8">
            <!-- Slider 1: Ticket Price -->
            <div>
              <div class="flex items-center justify-between mb-3">
                <label class="text-text-muted text-xs uppercase tracking-widest font-semibold">
                  Ticket Price
                </label>
                <span class="font-display text-3xl text-white leading-none">
                  KES {{ calc.price.toLocaleString() }}
                </span>
              </div>
              <input
                v-model.number="calc.price"
                type="range"
                min="100"
                max="5000"
                step="50"
                class="w-full h-2 appearance-none bg-white/10 cursor-pointer"
                style="accent-color: #E8002D;"
              />
              <div class="flex justify-between text-text-muted text-xs mt-2">
                <span>KES 100</span>
                <span>KES 5,000</span>
              </div>
            </div>

            <!-- Slider 2: Expected Viewers -->
            <div>
              <div class="flex items-center justify-between mb-3">
                <label class="text-text-muted text-xs uppercase tracking-widest font-semibold">
                  Expected Viewers
                </label>
                <span class="font-display text-3xl text-white leading-none">
                  {{ calc.viewers.toLocaleString() }}
                </span>
              </div>
              <input
                v-model.number="calc.viewers"
                type="range"
                min="10"
                max="10000"
                step="10"
                class="w-full h-2 appearance-none bg-white/10 cursor-pointer"
                style="accent-color: #E8002D;"
              />
              <div class="flex justify-between text-text-muted text-xs mt-2">
                <span>10 viewers</span>
                <span>10,000 viewers</span>
              </div>
            </div>
          </div>

          <!-- Results -->
          <div class="mt-8 pt-8 border-t border-white/10 space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-text-muted text-sm">Gross revenue</span>
              <span class="text-white font-semibold">KES {{ grossRevenue.toLocaleString() }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-text-muted text-sm">Platform fee (30%)</span>
              <span class="text-text-muted text-sm">− KES {{ platformFee.toLocaleString() }}</span>
            </div>
            <div class="flex items-center justify-between pt-4 border-t border-white/10">
              <span class="text-white font-bold text-base">Your Cut</span>
              <span class="font-display text-4xl text-status-success leading-none">
                KES {{ yourEarnings.toLocaleString() }}
              </span>
            </div>
          </div>

          <p class="text-text-muted text-xs text-center mt-6 flex items-center justify-center gap-1.5">
            <Smartphone class="w-4 h-4 shrink-0" />
            Paid directly to your M-Pesa within 24h
          </p>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 3: HOW PAYOUTS WORK                                   -->
    <!-- ============================================================ -->
    <section id="how-it-works" class="py-24">
      <div class="max-w-5xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="mb-16">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Payouts</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Simple, transparent payouts
          </h2>
          <p class="text-text-muted text-sm mt-4">Under 24 hours end-to-end</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-0">
          <div
            v-for="(step, i) in payoutSteps"
            :key="step.title"
            class="relative flex flex-col items-center text-center px-8"
          >
            <div class="w-16 h-16 rounded-full bg-accent-red/10 border-2 border-accent-red/40
                        flex items-center justify-center mb-5 shrink-0 z-10 relative">
              <span class="font-display text-2xl text-accent-red">{{ i + 1 }}</span>
            </div>

            <div
              v-if="i < payoutSteps.length - 1"
              class="hidden md:block absolute top-8 left-[calc(50%+2rem)] right-0 h-px"
              style="background: linear-gradient(to right, rgba(232,0,45,0.5), transparent);"
            />

            <div class="mb-4 text-accent-red">
              <component :is="step.icon" class="w-8 h-8 mx-auto" />
            </div>

            <h3 class="text-white font-bold text-base mb-2">{{ step.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed max-w-[200px]">{{ step.desc }}</p>
          </div>
        </div>

        <p class="text-center text-text-muted text-sm mt-12">
          End-to-end payout time:
          <span class="text-status-success font-bold ml-1">Under 24 hours</span>
        </p>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 4: PLATFORM FEATURES                                  -->
    <!-- ============================================================ -->
    <section class="py-24 bg-bg-surface border-y border-white/5">
      <div class="max-w-5xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="mb-14">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Platform</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Everything you need
          </h2>
          <p class="text-text-muted text-base mt-4 max-w-lg">
            A complete toolkit for monetizing live knowledge sessions.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div
            v-for="feature in platformFeatures"
            :key="feature.title"
            class="bg-bg-elevated rounded-xl p-6 flex gap-5 border border-white/5
                   hover:border-accent-red/20 transition-colors duration-200 group"
          >
            <div class="w-12 h-12 rounded-lg bg-accent-red/10 border border-accent-red/20
                        flex items-center justify-center shrink-0 group-hover:bg-accent-red/20 transition-colors duration-200">
              <component :is="feature.icon" class="w-6 h-6 text-accent-red" />
            </div>
            <div>
              <h3 class="text-white font-bold text-base mb-1.5">{{ feature.title }}</h3>
              <p class="text-text-muted text-sm leading-relaxed">{{ feature.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 5: CTA                                                -->
    <!-- ============================================================ -->
    <section class="py-24">
      <div class="max-w-5xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="bg-bg-elevated border border-white/10 rounded-xl p-12 md:p-16">
          <h2 class="font-display text-5xl md:text-6xl uppercase text-white leading-none mb-4">
            Ready to go live?
          </h2>
          <p class="text-text-muted text-base mb-8 max-w-lg leading-relaxed">
            Test the platform first — no account needed. When you're confident it works,
            create an account and start charging your audience.
          </p>

          <div class="flex flex-wrap items-center gap-4 mb-10">
            <RouterLink
              to="/try"
              class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg bg-accent-red
                     hover:bg-accent-red-hover text-white font-bold text-sm
                     transition-all duration-200 active:scale-95"
            >
              Try a Free Stream First
            </RouterLink>
            <RouterLink
              to="/register"
              class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg border border-white/20
                     hover:border-white/40 text-white text-sm font-medium
                     transition-all hover:bg-white/5"
            >
              Create Account
            </RouterLink>
          </div>

          <div class="flex flex-wrap items-center gap-6 text-text-muted text-sm">
            <span class="flex items-center gap-1.5">
              <Check class="w-4 h-4 text-status-success shrink-0" /> No monthly fees
            </span>
            <span class="flex items-center gap-1.5">
              <Check class="w-4 h-4 text-status-success shrink-0" /> M-Pesa payouts
            </span>
            <span class="flex items-center gap-1.5">
              <Check class="w-4 h-4 text-status-success shrink-0" /> 70% revenue share
            </span>
          </div>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Ticket, CheckCircle, Smartphone, Monitor, Shield, BarChart2, Target, MessageCircle, Check } from 'lucide-vue-next'

const payoutSteps = [
  {
    icon: Ticket,
    title: 'Ticket Sold',
    desc: 'A viewer pays via M-Pesa STK Push and gets instant access.',
  },
  {
    icon: CheckCircle,
    title: 'Payment Confirmed',
    desc: 'IntaSend confirms the payment and credits the platform.',
  },
  {
    icon: Smartphone,
    title: 'M-Pesa Sent',
    desc: 'Your 70% share is sent directly to your M-Pesa account.',
  },
]

const platformFeatures = [
  {
    icon: Monitor,
    title: 'HD Live Streaming',
    desc: 'Stream from OBS, your phone, or any encoder. We handle HLS conversion and global delivery automatically.',
  },
  {
    icon: Smartphone,
    title: 'M-Pesa Instant Payouts',
    desc: 'Viewers pay via M-Pesa STK Push. Your earnings land in your M-Pesa within 24 hours — no bank account needed.',
  },
  {
    icon: Shield,
    title: 'Device Anti-Piracy',
    desc: 'Unique stream tokens, device fingerprinting, and IP locking ensure only paying viewers can watch.',
  },
  {
    icon: BarChart2,
    title: 'Real-time Analytics',
    desc: 'Live viewer counts, peak attendance, ticket sales, and revenue — all visible in your host dashboard.',
  },
  {
    icon: Target,
    title: 'Custom Ticket Pricing',
    desc: 'Set any price from KES 0 (free events) to KES 100,000. Full pricing control is in your hands.',
  },
  {
    icon: MessageCircle,
    title: '24/7 Support',
    desc: 'WhatsApp, email, and in-app support. Our team is online during your events to make sure everything runs smoothly.',
  },
]

const calc = ref({ price: 500, viewers: 500 })
const grossRevenue = computed(() => calc.value.price * calc.value.viewers)
const platformFee = computed(() => Math.round(grossRevenue.value * 0.3))
const yourEarnings = computed(() => grossRevenue.value - platformFee.value)
</script>
