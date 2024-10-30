import { ThemeProvider } from "../components/providers/theme-provider";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Toaster } from "@repo/ui/components/shared/toaster";
import { WalletConnectionProvider } from "../components/providers/wallet-provider";
import { NetworkStatus } from "@repo/ui/components/common/network-status";
import HeadSeo from "@/components/head-seo";
import Script from "next/script";
import { GOOGLE_ANALYTICS_ID } from "@/lib/constants";
import { siteMetadata } from "@/lib/siteMetadata";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Script
        strategy="lazyOnload"
        src={`https://www.googletagmanager.com/gtag/js?id=${GOOGLE_ANALYTICS_ID}`}
      />

      <Script strategy="lazyOnload" id="google-analytics">
        {`
        window.dataLayer = window.dataLayer || [];
        function gtag(){dataLayer.push(arguments);}
        gtag('js', new Date());
        gtag('config', '${GOOGLE_ANALYTICS_ID}');
        `}
      </Script>

      <HeadSeo
        title={`${siteMetadata.companyName} `}
        description={siteMetadata.description}
        canonicalUrl={`${siteMetadata.siteUrl}`}
        ogImageUrl={siteMetadata.siteLogo}
        ogTwitterImage={siteMetadata.siteLogo}
        ogType={"website"}
      ></HeadSeo>
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
      >
        <WalletConnectionProvider>
          <Component {...pageProps} />
          <Toaster />
          <NetworkStatus />
        </WalletConnectionProvider>
      </ThemeProvider>
    </>
  );
}
