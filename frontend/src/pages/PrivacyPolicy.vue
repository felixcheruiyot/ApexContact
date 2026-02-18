<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
    <!-- Header -->
    <div class="mb-12">
      <p class="text-accent-red text-sm font-semibold uppercase tracking-widest mb-3">Legal</p>
      <h1 class="font-display text-5xl uppercase tracking-wide text-white mb-4">Privacy Policy</h1>
      <p class="text-text-muted">Last updated: February 18, 2026</p>
    </div>

    <div class="prose-apex space-y-10">

      <section>
        <h2>1. Who We Are</h2>
        <p>
          ApexContact operates a live sports streaming platform for boxing and car racing events.
          We are the data controller for all personal information collected through our website
          and services. By using ApexContact, you agree to the practices described in this policy.
        </p>
      </section>

      <section>
        <h2>2. Information We Collect</h2>

        <h3>Account Information</h3>
        <p>When you register, we collect your full name, email address, phone number (required for M-Pesa payments), and a hashed version of your password. We never store your password in plain text.</p>

        <h3>Payment Information</h3>
        <p>
          Payments are processed through <strong>IntaSend</strong> using M-Pesa STK Push. We collect
          your M-Pesa phone number and store payment reference numbers (IntaSend transaction IDs).
          We do not store M-Pesa PINs, card numbers, or any full financial credentials on our servers.
        </p>

        <h3>Device & Technical Information</h3>
        <p>To protect our content from piracy and enforce our one-device-per-ticket policy, we collect and process:</p>
        <ul>
          <li><strong>Device fingerprint</strong> — a hash derived from your browser's canvas rendering, screen resolution, colour depth, timezone, and language settings. This hash cannot be reversed to identify you as an individual.</li>
          <li><strong>IP address</strong> — your public IP address is recorded at the time of stream access and locked to your subscription token.</li>
          <li><strong>VPN/proxy detection signals</strong> — we check whether your connection routes through a known VPN, proxy, or data centre to enforce our content distribution agreements.</li>
          <li><strong>Stream session data</strong> — we track whether your unique stream token is active on more than one device simultaneously.</li>
        </ul>

        <h3>Usage Data</h3>
        <p>We collect event viewing history, stream quality metrics, and general platform analytics (pages visited, time spent) to improve our service.</p>
      </section>

      <section>
        <h2>3. How We Use Your Information</h2>
        <table>
          <thead>
            <tr>
              <th>Purpose</th>
              <th>Legal Basis</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Delivering the live stream you paid for</td>
              <td>Contract performance</td>
            </tr>
            <tr>
              <td>Processing M-Pesa payments via IntaSend</td>
              <td>Contract performance</td>
            </tr>
            <tr>
              <td>Preventing piracy and token sharing</td>
              <td>Legitimate interest</td>
            </tr>
            <tr>
              <td>Fraud detection and account security</td>
              <td>Legitimate interest</td>
            </tr>
            <tr>
              <td>Sending event reminders and receipts</td>
              <td>Consent / contract</td>
            </tr>
            <tr>
              <td>Improving platform performance</td>
              <td>Legitimate interest</td>
            </tr>
            <tr>
              <td>Complying with legal obligations</td>
              <td>Legal obligation</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section>
        <h2>4. Anti-Piracy Measures</h2>
        <p>
          ApexContact uses technical measures to protect the intellectual property of event promoters and rights holders.
          These include but are not limited to:
        </p>
        <ul>
          <li>Unique, single-use stream tokens issued per subscription</li>
          <li>IP address locking — your stream token is bound to the IP address used at first playback</li>
          <li>Device fingerprint binding — tokens are associated with the browser/device used at payment</li>
          <li>Concurrent session detection — opening your stream on a second device invalidates both sessions</li>
          <li>VPN and proxy detection — access via anonymising services may be restricted</li>
        </ul>
        <p>
          Attempting to circumvent these measures or share your stream token constitutes a breach of our
          Terms of Service and may result in permanent account suspension without a refund.
        </p>
      </section>

      <section>
        <h2>5. Data Sharing</h2>
        <p>We do not sell your personal data. We share data only with:</p>
        <ul>
          <li><strong>IntaSend</strong> — to process M-Pesa payments. Their privacy policy governs how they handle your phone number and transaction data.</li>
          <li><strong>Event Promoters</strong> — aggregated, non-identifiable analytics (ticket count, viewer peaks, revenue totals). We do not share your personal details with promoters.</li>
          <li><strong>Law Enforcement</strong> — when required by a valid legal order issued by Kenyan or applicable authorities.</li>
          <li><strong>Infrastructure Providers</strong> — our hosting and database providers process data solely on our instructions.</li>
        </ul>
      </section>

      <section>
        <h2>6. Data Retention</h2>
        <ul>
          <li><strong>Account data</strong> — retained for as long as your account is active, plus 3 years after deletion for legal compliance.</li>
          <li><strong>Payment records</strong> — retained for 7 years as required by Kenyan financial regulations.</li>
          <li><strong>Stream tokens</strong> — expire when the event ends. Session data in Redis is purged within 30 minutes of inactivity.</li>
          <li><strong>Device fingerprints</strong> — retained for 12 months after the last associated subscription.</li>
          <li><strong>Fraud flags</strong> — retained for 2 years after resolution.</li>
        </ul>
      </section>

      <section>
        <h2>7. Your Rights</h2>
        <p>Under applicable data protection law, you have the right to:</p>
        <ul>
          <li><strong>Access</strong> — request a copy of the personal data we hold about you</li>
          <li><strong>Correction</strong> — request that inaccurate data be corrected</li>
          <li><strong>Deletion</strong> — request erasure of your data (subject to legal retention obligations)</li>
          <li><strong>Restriction</strong> — request that we limit how we process your data</li>
          <li><strong>Portability</strong> — receive your data in a structured, machine-readable format</li>
          <li><strong>Objection</strong> — object to processing based on legitimate interests</li>
        </ul>
        <p>
          To exercise any of these rights, contact us at
          <a href="mailto:privacy@apexcontact.com" class="text-accent-red hover:underline">privacy@apexcontact.com</a>.
          We will respond within 30 days.
        </p>
      </section>

      <section>
        <h2>8. Cookies & Local Storage</h2>
        <p>We use browser <strong>localStorage</strong> to store your authentication token and user profile so you remain signed in. No third-party advertising cookies are used. We may use analytics cookies (first-party) to measure platform performance.</p>
      </section>

      <section>
        <h2>9. Security</h2>
        <p>
          We use industry-standard security measures: HTTPS/TLS for all data in transit, bcrypt password hashing,
          JWT token authentication, Redis session management with automatic expiry, and rate limiting on all API endpoints.
          Despite these measures, no system is completely secure and we cannot guarantee absolute security.
        </p>
      </section>

      <section>
        <h2>10. Children's Privacy</h2>
        <p>
          ApexContact is intended for users aged 18 and above. We do not knowingly collect personal information
          from minors. If you believe a child has registered, contact us immediately and we will delete the account.
        </p>
      </section>

      <section>
        <h2>11. Changes to This Policy</h2>
        <p>
          We may update this policy from time to time. When we make material changes, we will notify you via
          email or a prominent notice on the platform. Your continued use of ApexContact after the effective
          date constitutes acceptance of the updated policy.
        </p>
      </section>

      <section>
        <h2>12. Contact</h2>
        <p>For any privacy-related questions or requests:</p>
        <div class="card p-6 not-prose">
          <dl class="space-y-2 text-sm">
            <div class="flex gap-3"><dt class="text-text-muted w-24 shrink-0">Email</dt><dd class="text-white"><a href="mailto:privacy@apexcontact.com" class="text-accent-red hover:underline">privacy@apexcontact.com</a></dd></div>
            <div class="flex gap-3"><dt class="text-text-muted w-24 shrink-0">Platform</dt><dd class="text-white">ApexContact</dd></div>
          </dl>
        </div>
      </section>

    </div>
  </div>
</template>

<style scoped>
.prose-apex section {
  @apply border-b border-white/5 pb-10;
}
.prose-apex section:last-child {
  @apply border-b-0;
}
.prose-apex h2 {
  @apply font-display text-2xl uppercase tracking-wide text-white mb-4;
}
.prose-apex h3 {
  @apply text-white font-semibold text-base mb-2 mt-5;
}
.prose-apex p {
  @apply text-text-muted leading-relaxed mb-4;
}
.prose-apex ul {
  @apply list-disc list-inside space-y-2 text-text-muted mb-4 pl-2;
}
.prose-apex ul li strong {
  @apply text-white;
}
.prose-apex table {
  @apply w-full text-sm mb-4;
}
.prose-apex table th {
  @apply text-left text-text-muted text-xs uppercase tracking-wider pb-3 border-b border-white/10;
}
.prose-apex table td {
  @apply text-text-muted py-3 border-b border-white/5 align-top;
}
.prose-apex table tr td:first-child {
  @apply text-white pr-8;
}
</style>
