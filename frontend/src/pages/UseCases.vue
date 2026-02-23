<template>
  <div>
    <!-- Hero -->
    <section class="relative overflow-hidden min-h-[60vh] flex items-center">
      <div class="absolute inset-0 bg-gradient-to-br from-bg via-bg to-accent-orange/10" />
      <div class="absolute inset-0"
        style="background-image: radial-gradient(circle at 30% 60%, rgba(255,107,0,0.07) 0%, transparent 55%);" />
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
        <p class="text-accent-orange font-semibold text-sm uppercase tracking-widest mb-4">One Platform. Infinite Possibilities.</p>
        <h1 class="font-display text-6xl md:text-8xl uppercase tracking-wide text-white leading-none mb-6">
          Use <span class="text-accent-red">Cases</span>
        </h1>
        <p class="text-text-muted text-lg md:text-xl leading-relaxed max-w-2xl mx-auto mb-10">
          Whether you're a sales trainer, a mentor, a lawyer, or a consultant — Live Streamify
          turns your expertise into a live revenue stream, accessible to anyone with a phone.
        </p>
        <div class="flex items-center justify-center gap-4 flex-wrap">
          <RouterLink to="/register" class="btn-primary text-base px-8 py-4">Start Hosting Free</RouterLink>
          <RouterLink to="/promoters" class="btn-ghost text-base px-8 py-4">See Pricing</RouterLink>
        </div>
      </div>
    </section>

    <!-- Use Cases Grid -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center mb-14">
        <p class="text-accent-red text-sm font-semibold uppercase tracking-widest mb-3">What You Can Host</p>
        <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
          Built for Every Expert
        </h2>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="uc in useCases" :key="uc.title"
          class="card p-7 group hover:border hover:border-accent-red/30 transition-all duration-200 flex flex-col gap-4">
          <div class="flex items-start gap-4">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center text-2xl shrink-0"
              :style="`background: ${uc.bg}`">
              {{ uc.icon }}
            </div>
            <div>
              <h3 class="text-white font-bold text-base mb-1">{{ uc.title }}</h3>
              <span class="inline-block text-[10px] font-bold uppercase tracking-widest px-2 py-0.5 rounded-full"
                :style="`background: ${uc.tagBg}; color: ${uc.tagColor}`">
                {{ uc.tag }}
              </span>
            </div>
          </div>
          <p class="text-text-muted text-sm leading-relaxed">{{ uc.desc }}</p>
          <div class="border-t border-white/5 pt-4">
            <p class="text-xs text-text-muted font-semibold uppercase tracking-wider mb-2">Example events</p>
            <ul class="space-y-1">
              <li v-for="ex in uc.examples" :key="ex" class="text-xs text-text-muted flex items-center gap-2">
                <span class="w-1 h-1 rounded-full bg-accent-red shrink-0" />
                {{ ex }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- Who It's For -->
    <section class="bg-bg-surface border-y border-white/5 py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <p class="text-accent-red text-sm font-semibold uppercase tracking-widest mb-3">Who We Serve</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Your Audience is Waiting
          </h2>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="persona in personas" :key="persona.title"
            class="card p-6 text-center hover:border hover:border-accent-orange/30 transition-all">
            <div class="text-4xl mb-4">{{ persona.icon }}</div>
            <h3 class="text-white font-bold mb-2">{{ persona.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed">{{ persona.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Platform stats -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div v-for="stat in platformStats" :key="stat.label" class="card p-6 text-center">
          <p class="font-display text-5xl text-accent-red mb-2">{{ stat.value }}</p>
          <p class="text-white font-semibold text-sm">{{ stat.label }}</p>
          <p class="text-text-muted text-xs mt-1">{{ stat.sub }}</p>
        </div>
      </div>
    </section>

    <!-- Testimonials -->
    <section class="bg-bg-surface border-y border-white/5 py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <p class="text-accent-orange text-sm font-semibold uppercase tracking-widest mb-3">Real Hosts. Real Results.</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            What Our Hosts Say
          </h2>
          <p class="text-text-muted text-base mt-3 max-w-xl mx-auto">
            From sales trainers to immigration consultants — here's how Live Streamify is changing how experts earn.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="t in testimonials" :key="t.name"
            class="card p-7 flex flex-col gap-5 relative overflow-hidden">
            <!-- Decorative quote -->
            <div class="absolute top-4 right-6 font-display text-7xl text-white/5 select-none leading-none">"</div>

            <div class="flex gap-1">
              <span v-for="n in 5" :key="n" class="text-accent-orange text-sm">★</span>
            </div>
            <p class="text-text-muted text-sm leading-relaxed flex-1 relative z-10">"{{ t.quote }}"</p>
            <div class="flex items-center gap-3 border-t border-white/5 pt-5">
              <div class="w-11 h-11 rounded-full flex items-center justify-center font-bold text-sm shrink-0"
                :style="`background: ${t.avatarBg}; color: ${t.avatarColor}`">
                {{ t.initials }}
              </div>
              <div>
                <p class="text-white text-sm font-semibold">{{ t.name }}</p>
                <p class="text-text-muted text-xs">{{ t.role }}</p>
              </div>
              <span class="ml-auto text-xl">{{ t.icon }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
      <h2 class="font-display text-5xl md:text-6xl uppercase tracking-wide text-white mb-4">
        Your Turn to <span class="text-accent-red">Go Live</span>
      </h2>
      <p class="text-text-muted text-lg mb-10 max-w-xl mx-auto">
        Join hundreds of experts already monetising their knowledge on Live Streamify.
        No setup costs — you only pay when tickets sell.
      </p>
      <div class="flex items-center justify-center gap-4 flex-wrap">
        <RouterLink to="/register" class="btn-primary text-base px-8 py-4">Create Free Account</RouterLink>
        <RouterLink to="/promoters" class="btn-ghost text-base px-8 py-4">Learn More</RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

const useCases = [
  {
    icon: '🎯',
    title: 'Sales & Negotiation',
    tag: 'High Demand',
    tagBg: 'rgba(232,0,45,0.15)',
    tagColor: '#E8002D',
    bg: 'rgba(232,0,45,0.1)',
    desc: 'Sales trainers, closers, and business coaches can host live masterclasses on prospecting, cold calling, objection handling, and deal closing — the skills every professional needs.',
    examples: ['Cold calling & outreach bootcamp', 'Closing high-ticket deals', 'Objection handling workshop', 'B2B sales masterclass'],
  },
  {
    icon: '🚀',
    title: 'Personal Development',
    tag: 'Growth',
    tagBg: 'rgba(255,107,0,0.15)',
    tagColor: '#FF6B00',
    bg: 'rgba(255,107,0,0.1)',
    desc: 'Life coaches, productivity experts, and motivational speakers can host live sessions on mindset, goal-setting, time management, and building a successful career.',
    examples: ['Morning mindset masterclass', 'Goal-setting workshop', 'Confidence & public speaking', 'Productivity system deep-dive'],
  },
  {
    icon: '🏫',
    title: 'Mentoring & Coaching',
    tag: 'Education',
    tagBg: 'rgba(255,107,0,0.15)',
    tagColor: '#FF6B00',
    bg: 'rgba(255,107,0,0.1)',
    desc: 'Run paid one-on-one or group coaching sessions. Mentors can offer career advice, business coaching, or personal development — live, interactive, and monetised.',
    examples: ['Career coaching masterclass', 'Business growth workshop', 'Leadership session', 'CV & interview prep'],
  },
  {
    icon: '🌍',
    title: 'Visa & Immigration',
    tag: 'Consulting',
    tagBg: 'rgba(0,200,83,0.15)',
    tagColor: '#00C853',
    bg: 'rgba(0,200,83,0.08)',
    desc: 'Immigration consultants and lawyers can host live visa application guides, country-specific briefings, and Q&A sessions that attendees pay to join.',
    examples: ['UK student visa walkthrough', 'Canada PR application guide', 'Schengen visa tips', 'US B1/B2 Q&A session'],
  },
  {
    icon: '📚',
    title: 'Education & Workshops',
    tag: 'Education',
    tagBg: 'rgba(255,107,0,0.15)',
    tagColor: '#FF6B00',
    bg: 'rgba(255,107,0,0.1)',
    desc: 'Teachers, trainers, and institutions can host live classes, skill bootcamps, and certification workshops — on demand and globally accessible.',
    examples: ['Python coding bootcamp', 'IELTS prep class', 'Financial literacy workshop', 'Graphic design masterclass'],
  },
  {
    icon: '🎵',
    title: 'Music & Performances',
    tag: 'Entertainment',
    tagBg: 'rgba(156,39,176,0.15)',
    tagColor: '#CE93D8',
    bg: 'rgba(156,39,176,0.08)',
    desc: 'Artists, DJs, and performers can sell virtual concert tickets and reach fans worldwide — no venue required, no travel costs.',
    examples: ['Live Afrobeats concert', 'DJ set stream', 'Spoken word night', 'Gospel music event'],
  },
  {
    icon: '💪',
    title: 'Fitness & Wellness',
    tag: 'Health',
    tagBg: 'rgba(0,200,83,0.15)',
    tagColor: '#00C853',
    bg: 'rgba(0,200,83,0.08)',
    desc: 'Personal trainers, yoga teachers, and wellness coaches can monetise live workout sessions and health talks for subscribers anywhere.',
    examples: ['Live HIIT class', 'Morning yoga session', 'Nutrition masterclass', 'Mental wellness talk'],
  },
  {
    icon: '⚖️',
    title: 'Legal Consultations',
    tag: 'Consulting',
    tagBg: 'rgba(0,200,83,0.15)',
    tagColor: '#00C853',
    bg: 'rgba(0,200,83,0.08)',
    desc: 'Lawyers and legal professionals can host live Q&A sessions, contract workshops, and know-your-rights seminars that people pay to attend.',
    examples: ['Tenant rights Q&A', 'Business contracts explained', 'Family law workshop', 'Employment law masterclass'],
  },
  {
    icon: '💼',
    title: 'Business & Finance',
    tag: 'Finance',
    tagBg: 'rgba(255,179,0,0.15)',
    tagColor: '#FFB300',
    bg: 'rgba(255,179,0,0.08)',
    desc: 'Financial advisors, investors, and entrepreneurs can stream investment talks, startup pitches, and financial planning sessions.',
    examples: ['Stock market masterclass', 'Forex trading workshop', 'Startup fundraising talk', 'SME growth seminar'],
  },
  {
    icon: '🍳',
    title: 'Cooking & Lifestyle',
    tag: 'Lifestyle',
    tagBg: 'rgba(255,107,0,0.15)',
    tagColor: '#FF6B00',
    bg: 'rgba(255,107,0,0.1)',
    desc: 'Chefs, food bloggers, and lifestyle creators can host live cooking shows and paid lifestyle workshops from any kitchen or studio.',
    examples: ['African cuisine masterclass', 'Baking & confectionery', 'Healthy meal prep session', 'Home barista class'],
  },
  {
    icon: '🎮',
    title: 'Gaming & Esports',
    tag: 'Entertainment',
    tagBg: 'rgba(156,39,176,0.15)',
    tagColor: '#CE93D8',
    bg: 'rgba(156,39,176,0.08)',
    desc: 'Esports organisers and streamers can monetise tournaments, gaming events, and interactive live streams with integrated ticketing.',
    examples: ['FIFA tournament final', 'Mobile legends championship', 'Gaming hackathon stream', 'Retro gaming night'],
  },
  {
    icon: '🙏',
    title: 'Faith & Community',
    tag: 'Community',
    tagBg: 'rgba(255,179,0,0.15)',
    tagColor: '#FFB300',
    bg: 'rgba(255,179,0,0.08)',
    desc: 'Churches, mosques, temples, and community groups can stream services, fundraisers, and community events with free or paid access.',
    examples: ['Sunday church service', 'Fundraiser gala stream', 'Youth conference', 'Community town hall'],
  },
]

const personas = [
  {
    icon: '🎯',
    title: 'Experts & Professionals',
    desc: 'Lawyers, consultants, coaches, and advisors who want to monetise their knowledge at scale.',
  },
  {
    icon: '📣',
    title: 'Event Organisers',
    desc: 'Conference producers, community leaders, and event hosts running ticketed live sessions online.',
  },
  {
    icon: '🎓',
    title: 'Educators & Trainers',
    desc: 'Teachers, trainers, and institutions delivering live courses and skill-building workshops.',
  },
  {
    icon: '🎤',
    title: 'Creators & Performers',
    desc: 'Musicians, comedians, and content creators who want to sell virtual audience access.',
  },
]

const platformStats = [
  { value: '70%', label: 'Revenue to You', sub: 'You keep 70 cents of every shilling' },
  { value: 'M-Pesa', label: 'Built-in Payments', sub: 'No card needed — STK Push in seconds' },
  { value: '2s', label: 'Stream Latency', sub: 'HLS delivery, near real-time' },
  { value: '0 KES', label: 'Setup Cost', sub: 'Pay only when tickets sell' },
]

const testimonials = [
  {
    name: 'Kevin Otieno',
    role: 'Sales Trainer & Business Coach, Nairobi',
    initials: 'KO',
    avatarBg: 'rgba(232,0,45,0.2)',
    avatarColor: '#E8002D',
    icon: '🎯',
    quote: 'I run weekly live sales masterclasses — closing techniques, cold calling, objection handling. 400+ professionals attend every session. Revenue lands in my M-Pesa before I even go live.',
  },
  {
    name: 'Dr. James Odhiambo',
    role: 'Career Mentor & Executive Coach',
    initials: 'JO',
    avatarBg: 'rgba(255,107,0,0.2)',
    avatarColor: '#FF6B00',
    icon: '🏫',
    quote: 'My mentorship sessions now reach students across East Africa. The M-Pesa integration means no one is excluded — everyone pays in seconds. My income tripled in three months.',
  },
  {
    name: 'Aisha Kamau',
    role: 'Immigration Consultant, Nairobi',
    initials: 'AK',
    avatarBg: 'rgba(0,200,83,0.2)',
    avatarColor: '#00C853',
    icon: '🌍',
    quote: 'I turned my visa application workshops into a recurring income stream. Clients attend live, ask real questions, and pay with M-Pesa in under 20 seconds. I now run three cohorts a month.',
  },
  {
    name: 'Grace Wanjiku',
    role: 'Business Strategist & Entrepreneur, Mombasa',
    initials: 'GW',
    avatarBg: 'rgba(255,179,0,0.2)',
    avatarColor: '#FFB300',
    icon: '💼',
    quote: 'My live "Start & Scale" business bootcamp sold out in 48 hours. Participants from Uganda, Tanzania, and Nigeria joined. The platform made it feel like we were all in the same room.',
  },
  {
    name: 'Adv. Beatrice Njoroge',
    role: 'Advocate, High Court of Kenya',
    initials: 'BN',
    avatarBg: 'rgba(0,200,83,0.2)',
    avatarColor: '#00C853',
    icon: '⚖️',
    quote: "Hosting my monthly \"Know Your Rights\" session on Live Streamify has brought legal awareness to hundreds of Kenyans who couldn't afford an in-person consultation.",
  },
  {
    name: 'Brian Mwangi',
    role: 'Personal Development Coach & Speaker',
    initials: 'BM',
    avatarBg: 'rgba(156,39,176,0.2)',
    avatarColor: '#CE93D8',
    icon: '🚀',
    quote: "My live mindset and productivity workshops reach 600+ people weekly. I replaced my 9-to-5 income within 90 days. This platform is the real deal for African knowledge creators.",
  },
]
</script>
