export default function DisclaimerPage() {
  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <h1 className="text-2xl font-bold tracking-tight text-white mb-2">Disclaimer</h1>
      <p className="text-xs text-[#737373] mb-10">Last updated: May 2026</p>

      <div className="space-y-8 text-sm text-[#a3a3a3] leading-relaxed">

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Independent Project</h2>
          <p>
            CDLytics is an independent analytics project and is not affiliated with,
            endorsed by, or sponsored by Activision, Activision Blizzard, the Call of Duty
            League, Esports World Cup, or any team, organization, or player featured on
            the site.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Trademarks &amp; Ownership</h2>
          <p>
            All trademarks, team names, logos, game titles, and related branding belong to
            their respective owners. CDLytics makes no claim of ownership over any
            third-party intellectual property. Logos and names are used for identification
            purposes only in an informational, non-commercial context.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Data Accuracy</h2>
          <p>
            Statistics, results, and other data on CDLytics are compiled from publicly
            available sources and provided for informational and educational purposes.
            Data may be incomplete, delayed, or contain errors. CDLytics makes no
            representations about the accuracy or completeness of any information on
            the site.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Not for Betting or Financial Use</h2>
          <p>
            Information on CDLytics must not be used as the basis for betting, wagering,
            financial decisions, or any official league or competitive purposes. Always
            consult official sources for such decisions.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Portfolio &amp; Educational Purpose</h2>
          <p>
            CDLytics exists as a personal analytics and software portfolio project. It
            demonstrates data aggregation, visualization, and web development techniques
            using publicly available esports data.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Limitation of Liability</h2>
          <p>
            CDLytics is provided as-is. To the extent permitted by applicable law, its
            operators accept no liability for any loss, damage, or inconvenience arising
            from use of or reliance on information presented on this site.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">Contact &amp; Takedown</h2>
          <p>
            For questions, content concerns, or takedown requests, contact{' '}
            <a href="mailto:hello@cdlytics.com" className="text-white hover:underline">
              hello@cdlytics.com
            </a>
            . Rights-holder requests will be addressed promptly.
          </p>
        </section>

      </div>
    </div>
  )
}
