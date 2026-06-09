export default function TermsPage() {
  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <h1 className="text-2xl font-bold tracking-tight text-white mb-2">
        Terms of Use
      </h1>
      <p className="text-xs text-[#737373] mb-10">Last updated: May 2026</p>

      <div className="space-y-8 text-sm text-[#a3a3a3] leading-relaxed">
        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            Independent Project
          </h2>
          <p>
            CDLytics is an independent analytics and portfolio project. It is
            not affiliated with, endorsed by, or sponsored by Activision,
            Activision Blizzard, the Call of Duty League, Esports World Cup, or
            any of the teams, organizations, or players listed on the site.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            Intellectual Property
          </h2>
          <p>
            All trademarks, team names, logos, game titles, and related branding
            displayed on CDLytics belong to their respective owners. Their
            appearance on this site does not imply any affiliation or
            endorsement. If you are a rights holder and believe content should
            be removed, see the Contact &amp; Takedown section below.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            Acceptable Use
          </h2>
          <p>
            CDLytics is provided for informational, educational, and portfolio
            purposes. By using the site you agree not to scrape, reproduce, or
            redistribute its content at scale without permission. The site is
            intended for personal, non-commercial use.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            No Accounts
          </h2>
          <p>
            CDLytics does not currently offer user accounts. No registration is
            required to access the site.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            Accuracy of Data
          </h2>
          <p>
            Statistics and information on CDLytics are compiled from publicly
            available sources and may be incomplete, delayed, or inaccurate. Do
            not rely on this site for betting, financial decisions, or any
            official league purposes. Always verify information with official
            sources.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            No Warranty
          </h2>
          <p>
            CDLytics is provided as-is with no guarantees of accuracy,
            availability, or fitness for any particular purpose. Use of the site
            is at your own risk.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            Contact &amp; Takedown
          </h2>
          <p>
            If you are a rights holder and believe that content on CDLytics
            infringes your intellectual property, or if you have any other
            concern, please contact us at{" "}
            <a
              href="mailto:hello@cdlytics.com"
              className="text-white hover:underline"
            >
              hello@cdlytics.com
            </a>{" "}
            and we will respond promptly.
          </p>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-widest text-[#737373] mb-3">
            Changes
          </h2>
          <p>
            These terms may be updated from time to time. The date at the top of
            this page reflects the most recent revision. Continued use of the
            site after changes constitutes acceptance.
          </p>
        </section>
      </div>
    </div>
  );
}
