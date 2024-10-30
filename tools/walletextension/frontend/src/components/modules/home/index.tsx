import React from "react";
import Connected from "./connected";
import Disconnected from "./disconnected";
import { Skeleton } from "@repo/ui/components/shared/skeleton";
import useWalletStore from "@/stores/wallet-store";

const Home = () => {
  const { walletConnected, loading } = useWalletStore();

  return (
    <div className="w-[800px] mx-auto">
      {loading ? (
        <Skeleton className="h-[400px]" />
      ) : walletConnected ? (
        <Connected />
      ) : (
        <Disconnected />
      )}
    </div>
  );
};

export default Home;
