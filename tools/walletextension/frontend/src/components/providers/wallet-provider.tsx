import { createContext, useContext, useEffect } from "react";
import { ethereum } from "@repo/ui/lib/utils";
import { showToast } from "@repo/ui/components/shared/use-toast";
import { ToastType } from "@repo/ui/lib/enums/toast";
import useWalletStore from "@/stores/wallet-store";

const WalletConnectionContext = createContext<void | ReturnType<
  typeof useWalletStore
>>(undefined);

export const useWalletConnection = () => {
  const context = useContext(WalletConnectionContext);
  if (!context) {
    throw new Error(
      "useWalletConnection must be used within a WalletConnectionProvider"
    );
  }
  return context;
};

export const WalletConnectionProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const walletStore = useWalletStore();
  const { initializeGateway, initializeProvider } = walletStore;

  useEffect(() => {
    const initializeWallet = async () => {
      if (ethereum && ethereum.isMetaMask) {
        try {
          initializeProvider();
          initializeGateway();
        } catch (error) {
          console.error("Failed to initialize wallet:", error);
          showToast(
            ToastType.DESTRUCTIVE,
            "Failed to initialize wallet. Please refresh and try again."
          );
        }
      } else {
        showToast(
          ToastType.WARNING,
          "MetaMask not detected. Some features may be unavailable."
        );
      }
    };

    initializeWallet();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ethereum]);

  return (
    <WalletConnectionContext.Provider value={walletStore}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
