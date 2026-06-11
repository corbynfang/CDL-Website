import { Helmet } from "react-helmet-async";

const SITE_NAME = "CDLytics";
const BASE_URL = "https://cdlytics.com";
const DEFAULT_IMAGE = `${BASE_URL}/assets/logos/OpticTexasLogo.png`;
const DEFAULT_DESC =
  "CDLytics — independent Call of Duty League statistics. Player K/D ratios, tournament results, team rosters, and transfer history.";

interface PageMetaProps {
  title?: string;
  description?: string;
  canonical?: string;
  image?: string;
  type?: "website" | "profile" | "article";
  noIndex?: boolean;
}

const PageMeta = ({
  title,
  description = DEFAULT_DESC,
  canonical,
  image = DEFAULT_IMAGE,
  type = "website",
  noIndex = false,
}: PageMetaProps) => {
  const fullTitle = title ? `${title} | ${SITE_NAME}` : `${SITE_NAME} — CDL Stats & Analytics`;
  const canonicalUrl = canonical ? `${BASE_URL}${canonical}` : undefined;

  return (
    <Helmet>
      <title>{fullTitle}</title>
      <meta name="description" content={description} />
      {noIndex && <meta name="robots" content="noindex, nofollow" />}
      {canonicalUrl && <link rel="canonical" href={canonicalUrl} />}

      <meta property="og:site_name" content={SITE_NAME} />
      <meta property="og:type" content={type} />
      <meta property="og:title" content={fullTitle} />
      <meta property="og:description" content={description} />
      <meta property="og:image" content={image} />
      {canonicalUrl && <meta property="og:url" content={canonicalUrl} />}

      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={fullTitle} />
      <meta name="twitter:description" content={description} />
      <meta name="twitter:image" content={image} />
    </Helmet>
  );
};

export default PageMeta;
