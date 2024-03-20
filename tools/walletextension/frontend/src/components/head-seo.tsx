import Head from "next/head";
import { siteMetadata } from "../lib/siteMetadata";
import { SeoProps } from "@/types/interfaces";

const HeadSeo = ({
  title,
  description,
  canonicalUrl,
  ogTwitterImage,
  ogImageUrl,
  ogType,
  children,
}: SeoProps) => {
  return (
    <Head>
      {/* Basic metadata */}
      <title>{title}</title>
      <meta charSet="utf-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
      <meta name={siteMetadata.metaTitle} content={siteMetadata.companyName} />
      {/* Beagle Security */}
      {/*<meta*/}
      {/*  name="_vgeujvlkxz15hyr8vbuvqxnfmzlkm059"*/}
      {/*  signature="_vd3udx2g2hfn9zclob5cat43b94q7fyk"*/}
      {/*></meta>*/}
      {/* twitter metadata */}
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:site" content={siteMetadata.twitterHandle} />
      <meta name="twitter:title" content={title} />
      <meta name="twitter:description" content={description} />
      <meta name="twitter:image" content={decodeURIComponent(ogTwitterImage)} />
      {/* canonical link */}
      <link rel="canonical" href={canonicalUrl} />
      {/* open graph metadata */}
      <meta property="og:locale" content="en_US" />
      <meta property="og:site_name" content={siteMetadata.companyName} />
      <meta property="og:type" content={ogType} key="og-type" />
      <meta property="og:title" content={title} key="og-title" />
      <meta
        property="og:description"
        content={description}
        key="og-description"
      />
      <meta
        property="og:image"
        content={decodeURIComponent(ogImageUrl)}
        key="og-image"
      />
      <meta property="og:url" content={canonicalUrl} key="og-url" />
      {children}
    </Head>
  );
};

export default HeadSeo;
