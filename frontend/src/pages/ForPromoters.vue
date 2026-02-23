<template>
  <div class="bg-bg text-white">

    <!-- ============================================================ -->
    <!-- SECTION 1: HERO                                               -->
    <!-- ============================================================ -->
    <section class="relative min-h-screen flex items-center overflow-hidden">
      <!-- Background image -->
      <img
        src="https://images.pexels.com/photos/3184340/pexels-photo-3184340.jpeg?auto=compress&cs=tinysrgb&w=1260&q=80"
        alt="Creator speaking"
        class="absolute inset-0 w-full h-full object-cover object-center"
      />
      <!-- Dark overlays -->
      <div class="absolute inset-0 bg-gradient-to-r from-black/95 via-black/70 to-transparent pointer-events-none" />
      <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent pointer-events-none" />

      <!-- Content: left side only -->
      <div class="relative z-10 max-w-7xl mx-auto px-6 sm:px-8 lg:px-16 py-24 w-full">
        <div class="max-w-xl">
          <!-- Pill label -->
          <div class="inline-flex items-center bg-accent-red/10 border border-accent-red/30 text-accent-red text-xs font-bold uppercase tracking-widest px-4 py-2 rounded-full mb-8">
            💰 EARN FROM YOUR KNOWLEDGE
          </div>

          <!-- H1 -->
          <h1 class="font-display text-6xl md:text-7xl uppercase text-white leading-none tracking-wide mb-6">
            Turn Your Expertise<br />
            Into Income.<br />
            <span class="text-accent-red">Live.</span>
          </h1>

          <!-- Body text -->
          <p class="text-white/80 text-lg leading-relaxed mb-8 max-w-lg">
            Join 12,400 African creators earning from business, legal, fitness, education sessions —
            paid directly to M-Pesa.
          </p>

          <!-- Stats row -->
          <div class="flex flex-wrap gap-4 mb-10">
            <span
              v-for="stat in heroStats"
              :key="stat"
              class="bg-white/10 backdrop-blur-sm border border-white/20 text-white text-sm font-semibold px-4 py-2 rounded-full"
            >
              {{ stat }}
            </span>
          </div>

          <!-- CTAs -->
          <div class="flex flex-wrap items-center gap-4">
            <RouterLink
              to="/register"
              class="bg-white text-black font-bold text-base px-8 py-4 rounded-full hover:bg-white/90 transition-colors duration-200 shadow-lg"
            >
              Start Hosting Free
            </RouterLink>
            <a
              href="#how-it-works"
              class="border-2 border-white/40 text-white font-medium text-base px-6 py-4 rounded-full hover:border-white/70 transition-colors duration-200"
            >
              See How It Works ↓
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 2: TRUST STRIP                                        -->
    <!-- ============================================================ -->
    <section class="bg-bg-surface border-y border-white/5 py-8">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="flex flex-wrap items-center gap-4 justify-center md:justify-start">
          <span class="text-text-muted text-sm font-medium shrink-0 mr-2">
            Trusted by creators from:
          </span>
          <span
            v-for="city in cities"
            :key="city.name"
            class="inline-flex items-center gap-1.5 bg-bg-elevated border border-white/10 text-white text-sm font-medium px-4 py-2 rounded-full"
          >
            <span>{{ city.flag }}</span>
            {{ city.name }}
          </span>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 3: CREATOR SUCCESS STORIES                           -->
    <!-- ============================================================ -->
    <section class="py-24">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16">
        <!-- Heading -->
        <div class="mb-14">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Social Proof</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Creators already winning
          </h2>
        </div>

        <!-- 3-column grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Card 1 -->
          <div class="bg-bg-elevated rounded-2xl p-6 flex flex-col gap-5 border border-white/5 hover:border-white/15 transition-colors duration-200">
            <div class="flex items-center gap-4">
              <img
                src="https://images.pexels.com/photos/3769021/pexels-photo-3769021.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="Amara Kioni"
                class="w-16 h-16 rounded-full object-cover ring-2 ring-accent-red/30 shrink-0"
              />
              <div>
                <p class="text-white font-bold text-base leading-tight">Amara Kioni</p>
                <p class="text-text-muted text-sm mt-0.5">Immigration &amp; Visa Consultant</p>
              </div>
            </div>
            <blockquote class="text-white/80 text-sm leading-relaxed flex-1">
              "I made KES 85,000 from one live Q&amp;A session. 400 people paid to watch me explain UK visa applications."
            </blockquote>
            <div class="pt-4 border-t border-white/10">
              <div class="flex flex-wrap gap-2">
                <span class="bg-status-success/10 border border-status-success/20 text-status-success text-xs font-semibold px-3 py-1 rounded-full">
                  KES 85,000 / session
                </span>
                <span class="bg-white/5 border border-white/10 text-text-muted text-xs px-3 py-1 rounded-full">
                  18 sessions
                </span>
              </div>
            </div>
          </div>

          <!-- Card 2 -->
          <div class="bg-bg-elevated rounded-2xl p-6 flex flex-col gap-5 border border-white/5 hover:border-white/15 transition-colors duration-200">
            <div class="flex items-center gap-4">
              <img
                src="https://images.pexels.com/photos/1139743/pexels-photo-1139743.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="David Mwangi"
                class="w-16 h-16 rounded-full object-cover ring-2 ring-accent-red/30 shrink-0"
              />
              <div>
                <p class="text-white font-bold text-base leading-tight">David Mwangi</p>
                <p class="text-text-muted text-sm mt-0.5">Sales Trainer, Nairobi</p>
              </div>
            </div>
            <blockquote class="text-white/80 text-sm leading-relaxed flex-1">
              "My live sales masterclass has an audience of 2,400 repeat viewers. M-Pesa payouts hit in minutes."
            </blockquote>
            <div class="pt-4 border-t border-white/10">
              <div class="flex flex-wrap gap-2">
                <span class="bg-status-success/10 border border-status-success/20 text-status-success text-xs font-semibold px-3 py-1 rounded-full">
                  KES 340,000 total
                </span>
                <span class="bg-white/5 border border-white/10 text-text-muted text-xs px-3 py-1 rounded-full">
                  2,400 viewers
                </span>
              </div>
            </div>
          </div>

          <!-- Card 3 -->
          <div class="bg-bg-elevated rounded-2xl p-6 flex flex-col gap-5 border border-white/5 hover:border-white/15 transition-colors duration-200">
            <div class="flex items-center gap-4">
              <img
                src="https://images.pexels.com/photos/3182812/pexels-photo-3182812.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="Zara Hassan"
                class="w-16 h-16 rounded-full object-cover ring-2 ring-accent-red/30 shrink-0"
              />
              <div>
                <p class="text-white font-bold text-base leading-tight">Zara Hassan</p>
                <p class="text-text-muted text-sm mt-0.5">Fitness Coach, Kampala</p>
              </div>
            </div>
            <blockquote class="text-white/80 text-sm leading-relaxed flex-1">
              "I run weekly fitness sessions. 300 subscribers per session, consistent income every week."
            </blockquote>
            <div class="pt-4 border-t border-white/10">
              <div class="flex flex-wrap gap-2">
                <span class="bg-status-success/10 border border-status-success/20 text-status-success text-xs font-semibold px-3 py-1 rounded-full">
                  300 viewers/week
                </span>
                <span class="bg-white/5 border border-white/10 text-text-muted text-xs px-3 py-1 rounded-full">
                  52 sessions
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 4: REVENUE CALCULATOR                                 -->
    <!-- ============================================================ -->
    <section class="py-24 bg-bg-surface border-y border-white/5">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16">
        <!-- Heading -->
        <div class="text-center mb-14">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Earnings</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Calculate your earnings
          </h2>
          <p class="text-text-muted text-base mt-4 max-w-md mx-auto">
            70% of every ticket goes straight to you
          </p>
        </div>

        <!-- Calculator card -->
        <div class="bg-bg-elevated rounded-3xl p-8 md:p-10 max-w-2xl mx-auto border border-white/5">
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
                class="w-full h-2 rounded-full appearance-none bg-white/10 cursor-pointer slider-accent"
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
                class="w-full h-2 rounded-full appearance-none bg-white/10 cursor-pointer"
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
              <div>
                <span class="text-white font-bold text-base">Your Cut</span>
              </div>
              <span class="font-display text-4xl text-status-success leading-none">
                KES {{ yourEarnings.toLocaleString() }}
              </span>
            </div>
          </div>

          <!-- M-Pesa note -->
          <p class="text-text-muted text-xs text-center mt-6 flex items-center justify-center gap-1.5">
            <span>📱</span>
            Paid directly to your M-Pesa within 24h
          </p>

          <!-- CTA -->
          <RouterLink
            to="/register"
            class="mt-8 bg-accent-red hover:bg-accent-red-hover text-white font-bold text-base px-6 py-4 rounded-full w-full block text-center transition-colors duration-200"
          >
            Start Earning →
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 5: HOW PAYOUTS WORK                                   -->
    <!-- ============================================================ -->
    <section id="how-it-works" class="py-24">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16">
        <!-- Heading -->
        <div class="text-center mb-16">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Payouts</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Simple, transparent payouts
          </h2>
          <p class="text-text-muted text-sm mt-4">
            &lt; 24 hours end-to-end
          </p>
        </div>

        <!-- 3 steps horizontal -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-0">
          <div
            v-for="(step, i) in payoutSteps"
            :key="step.title"
            class="relative flex flex-col items-center text-center px-8"
          >
            <!-- Step number circle -->
            <div class="w-16 h-16 rounded-full bg-accent-red/10 border-2 border-accent-red/40 flex items-center justify-center mb-5 shrink-0 z-10 relative">
              <span class="font-display text-2xl text-accent-red">{{ i + 1 }}</span>
            </div>

            <!-- Arrow connector (between steps) -->
            <div
              v-if="i < payoutSteps.length - 1"
              class="hidden md:block absolute top-8 left-[calc(50%+2rem)] right-0 h-px"
              style="background: linear-gradient(to right, rgba(232,0,45,0.5), transparent);"
            />

            <!-- Icon -->
            <div class="text-3xl mb-4">{{ step.icon }}</div>

            <!-- Step title -->
            <h3 class="text-white font-bold text-base mb-2">{{ step.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed max-w-[200px]">{{ step.desc }}</p>
          </div>
        </div>

        <!-- Timeframe note -->
        <p class="text-center text-text-muted text-sm mt-12">
          End-to-end payout time:
          <span class="text-status-success font-bold ml-1">Under 24 hours</span>
        </p>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 6: PLATFORM FEATURES                                  -->
    <!-- ============================================================ -->
    <section class="py-24 bg-bg-surface border-y border-white/5">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16">
        <!-- Heading -->
        <div class="mb-14">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Platform</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Everything you need
          </h2>
          <p class="text-text-muted text-base mt-4 max-w-lg">
            A complete toolkit for creators ready to monetize their knowledge.
          </p>
        </div>

        <!-- 2-column grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div
            v-for="feature in platformFeatures"
            :key="feature.title"
            class="bg-bg-elevated rounded-2xl p-6 flex gap-5 border border-white/5 hover:border-accent-red/20 transition-colors duration-200 group"
          >
            <!-- Icon -->
            <div class="w-12 h-12 rounded-xl bg-accent-red/10 border border-accent-red/20 flex items-center justify-center shrink-0 group-hover:bg-accent-red/20 transition-colors duration-200">
              <span class="text-2xl">{{ feature.icon }}</span>
            </div>
            <!-- Content -->
            <div>
              <h3 class="text-white font-bold text-base mb-1.5">{{ feature.title }}</h3>
              <p class="text-text-muted text-sm leading-relaxed">{{ feature.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- SECTION 7: CTA                                                -->
    <!-- ============================================================ -->
    <section class="py-24 relative overflow-hidden">
      <!-- Red gradient background -->
      <div class="absolute inset-0" style="background: linear-gradient(135deg, #E8002D 0%, #C0001F 50%, #8B0000 100%);" />
      <!-- Subtle texture overlay -->
      <div class="absolute inset-0 opacity-10" style="background-image: radial-gradient(circle at 20% 50%, rgba(255,255,255,0.3) 0%, transparent 50%), radial-gradient(circle at 80% 20%, rgba(255,255,255,0.2) 0%, transparent 40%);" />

      <div class="relative z-10 max-w-3xl mx-auto px-6 text-center">
        <h2 class="font-display text-5xl md:text-6xl uppercase text-white leading-none mb-4">
          Ready to go live?
        </h2>
        <p class="text-white/80 text-lg mb-10 max-w-xl mx-auto">
          Start for free. No credit card. No technical setup.
        </p>
        <RouterLink
          to="/register"
          class="inline-block bg-white text-black font-bold text-lg px-10 py-5 rounded-full hover:bg-white/90 transition-colors duration-200 shadow-xl shadow-black/30"
        >
          Create Your Host Account
        </RouterLink>

        <!-- Trust indicators below CTA -->
        <div class="flex flex-wrap items-center justify-center gap-6 mt-10 text-white/60 text-sm">
          <span class="flex items-center gap-1.5">
            <span class="text-white/80">✓</span> No monthly fees
          </span>
          <span class="flex items-center gap-1.5">
            <span class="text-white/80">✓</span> M-Pesa payouts
          </span>
          <span class="flex items-center gap-1.5">
            <span class="text-white/80">✓</span> Go live in minutes
          </span>
          <span class="flex items-center gap-1.5">
            <span class="text-white/80">✓</span> 70% revenue share
          </span>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'

// Hero stats strip
const heroStats = [
  'KES 4.2M+ Paid',
  '70% Revenue Share',
  'Zero Setup Cost',
]

// City trust strip
const cities = [
  { name: 'Nairobi', flag: '🇰🇪' },
  { name: 'Lagos', flag: '🇳🇬' },
  { name: 'Kampala', flag: '🇺🇬' },
  { name: 'Accra', flag: '🇬🇭' },
  { name: 'Dar es Salaam', flag: '🇹🇿' },
]

// Payout steps
const payoutSteps = [
  {
    icon: '🎟️',
    title: 'Ticket Sold',
    desc: 'A viewer pays via M-Pesa STK Push and gets instant access.',
  },
  {
    icon: '✅',
    title: 'Payment Confirmed',
    desc: 'IntaSend confirms the payment and credits the platform.',
  },
  {
    icon: '📱',
    title: 'M-Pesa Sent',
    desc: 'Your 70% share is sent directly to your M-Pesa account.',
  },
]

// Platform features
const platformFeatures = [
  {
    icon: '📺',
    title: 'HD Live Streaming',
    desc: 'Stream from OBS, your phone, or any encoder. We handle HLS conversion and global delivery automatically.',
  },
  {
    icon: '📱',
    title: 'M-Pesa Instant Payouts',
    desc: 'Viewers pay via M-Pesa STK Push. Your earnings land in your M-Pesa within 24 hours — no bank account needed.',
  },
  {
    icon: '🛡️',
    title: 'Device Anti-Piracy',
    desc: 'Unique stream tokens, device fingerprinting, and IP locking ensure only paying viewers can watch.',
  },
  {
    icon: '📊',
    title: 'Real-time Analytics',
    desc: 'Live viewer counts, peak attendance, ticket sales, and revenue — all visible in your host dashboard.',
  },
  {
    icon: '🎯',
    title: 'Custom Ticket Pricing',
    desc: 'Set any price from KES 0 (free events) to KES 100,000. Full pricing control is in your hands.',
  },
  {
    icon: '💬',
    title: '24/7 Support',
    desc: 'WhatsApp, email, and in-app support. Our team is online during your events to make sure everything runs smoothly.',
  },
]

// Earnings calculator
const calc = ref({ price: 500, viewers: 500 })
const grossRevenue = computed(() => calc.value.price * calc.value.viewers)
const platformFee = computed(() => Math.round(grossRevenue.value * 0.3))
const yourEarnings = computed(() => grossRevenue.value - platformFee.value)
</script>
