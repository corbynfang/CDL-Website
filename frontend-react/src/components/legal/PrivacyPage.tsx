export default function PrivacyPage() {
  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <h1 className="text-2xl font-bold tracking-tight text-white mb-2">Privacy Policy</h1>
      <p className="text-xs text-[#737373] mb-10">Last updated: May 2026</p>

      <div className="space-y-8 text-sm text-[#a3a3a3] leading-relaxed">

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Overview</h2>
          <p>
            CDLytics is an independent analytics project. We do not sell data, run
            advertising, or build user profiles. This page explains what limited technical
            data may be processed when you use the site.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">No User Accounts</h2>
          <p>
            CDLytics does not currently offer user accounts, registration, or login. No
            personal information is collected through account creation.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Technical Logs</h2>
          <p>
            Standard web server and infrastructure logs may record basic technical data such
            as IP addresses, browser type, pages visited, and request timestamps. This data
            is processed by our hosting, security, and logging providers for the purposes of
            security monitoring, debugging, and performance. It is not used to identify
            individual users and is not shared or sold.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Cookies</h2>
          <p>
            CDLytics does not use tracking cookies, advertising cookies, or third-party
            analytics services. Any browser storage used is limited to essential site
            functionality.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Data Sources &amp; Provenance</h2>
          <p>
            Match, player, and team statistics are compiled from publicly available sources.
            Source and provenance metadata is retained internally for data verification and
            quality auditing purposes. This data is not shared publicly.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Third-Party Infrastructure</h2>
          <p>
            CDLytics is hosted on third-party services including cloud hosting and database
            providers. These providers may process basic technical data as part of normal
            operation, subject to their own privacy policies.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Contact</h2>
          <p>
            For privacy-related questions or data removal requests, contact us at{' '}
            <a href="mailto:hello@cdlytics.com" className="text-white hover:underline">
              hello@cdlytics.com
            </a>.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Changes</h2>
          <p>
            This policy may be updated occasionally. The date at the top of this page
            reflects the most recent revision.
          </p>
        </section>

      </div>
    </div>
  )
}
