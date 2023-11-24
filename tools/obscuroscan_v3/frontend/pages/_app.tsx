import { useState } from "react";
import { ThemeProvider } from "@/src/components/providers/theme-provider";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import {
  QueryClient,
  MutationCache,
  QueryClientProvider,
} from "@tanstack/react-query";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Toaster } from "@/src/components/ui/toaster";
import { useToast } from "@/src/components/ui/use-toast";
import { WalletConnectionProvider } from "@/src/components/providers/wallet-provider";
import { NetworkStatus } from "@/src/components/modules/common/network-status";

export default function App({ Component, pageProps }: AppProps) {
  const { toast } = useToast();

  const mutationCache = new MutationCache({
    onSuccess: (mutation: any) => {
      if (mutation?.message) {
        toast({
          description: mutation?.message,
        });
      }
    },
    onError: (error: any, mutation: any) => {
      if (error?.response?.data?.message) {
        toast({
          description: mutation?.message,
        });
      }
    },
  });

  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            refetchOnWindowFocus: false,
            staleTime: 300000,
          },
        },
        mutationCache,
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
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
          <ReactQueryDevtools initialIsOpen={false} />
        </WalletConnectionProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
}
