import { ethers } from "ethers";
import { currentNetwork } from "../lib/utils";
import {
  initializeSigner,
  setupEventListeners,
} from "../lib/utils/walletEvents";
import { getEthereumProvider, handleError } from "../lib/utils/walletUtils";
import { toast } from "../components/shared/use-toast";
import { requestMethods } from "../routes";
import { ToastType } from "../lib/enums/toast";
import { StoreSet, StoreGet } from "../lib/types/common";

export const walletService = {
  initializeProvider: async (set: StoreSet, get: StoreGet) => {
    try {
      let detectedProvider = await getEthereumProvider();
      if (!detectedProvider) {
        let attempts = 0;
        const maxAttempts = 5;
        const intervalId = setInterval(async () => {
          attempts++;
          try {
            detectedProvider = await getEthereumProvider();
            if (detectedProvider) {
              clearInterval(intervalId);
              set({ provider: detectedProvider, loading: false });
              await walletService.proceedWithProviderInitialization(set, get);
            } else if (attempts >= maxAttempts) {
              clearInterval(intervalId);
              toast({
                title: "Unable to connect",
                description:
                  "Please check if your web3 provider is installed and unlocked.",
                variant: ToastType.INFO,
              });
            }
          } catch (error) {
            console.error("Error detecting provider:", error);
          }
        }, 3000);
        return;
      }
      set({ provider: detectedProvider, loading: false });
      await walletService.proceedWithProviderInitialization(set, get);
    } catch (error) {
      console.error("Error initializing provider:", error);
      toast({
        title: "Provider Error",
        description: "Unable to initialize wallet provider.",
        variant: ToastType.DESTRUCTIVE,
      });
    }
  },

  proceedWithProviderInitialization: async (set: StoreSet, get: StoreGet) => {
    try {
      const { provider: detectedProvider } = get();

      const signer = await initializeSigner(detectedProvider!);
      const accounts = await walletService.getAccounts(detectedProvider!);
      const network = await detectedProvider?.getNetwork();
      const chainId = network?.chainId;
      const expectedChainId = currentNetwork.l2;

      const isWrongNetwork = chainId !== expectedChainId;

      // if (isWrongNetwork) {
      //   toast({
      //     description: "You are on the wrong network. Please switch to TEN.",
      //     variant: ToastType.INFO,
      //   });
      // }

      set({
        address: accounts[0],
        signer: signer,
        isWrongNetwork,
      });

      const cleanup = setupEventListeners(
        (address: string, chainId: number, isWrongNetwork: boolean) => {
          set({ address, chainId, isWrongNetwork });
        }
      );
      return cleanup;
    } catch (error) {
      console.error("Error during provider initialization:", error);
      toast({
        title: "Connection Error",
        description: "Failed to initialize provider.",
        variant: ToastType.DESTRUCTIVE,
      });
    }
  },

  connectWallet: async (set: StoreSet, get: StoreGet) => {
    try {
      await walletService.initializeProvider(set, get);

      set({
        walletConnected: true,
      });

      toast({
        description: "Connected to wallet!",
        variant: ToastType.INFO,
      });
    } catch (error) {
      console.error("Error connecting to wallet:", error);
      toast({
        title: "Error",
        description: "Failed to connect to wallet.",
        variant: ToastType.DESTRUCTIVE,
      });
    }
  },

  disconnectWallet: async (set: StoreSet, get: StoreGet) => {
    const { provider } = get();
    try {
      if (provider) {
        provider.removeAllListeners();
        set({
          provider: null,
          signer: null,
          address: "",
          walletConnected: false,
        });
        // window.location.reload();

        toast({
          title: "Disconnected",
          description: "Disconnected from wallet",
          variant: ToastType.INFO,
        });
      }
    } catch (error) {
      console.error("Error disconnecting from wallet:", error);
      toast({
        title: "Error",
        description: "Failed to disconnect wallet.",
        variant: ToastType.DESTRUCTIVE,
      });
    }
  },

  switchNetwork: async (set: StoreSet, get: StoreGet) => {
    const { provider } = get();
    if (!provider) {
      toast({
        title: "Error",
        description: "Please connect to wallet first.",
        variant: ToastType.DESTRUCTIVE,
      });
      return;
    }

    set({ loading: true });

    const desiredNetwork = ethers.utils.hexValue(currentNetwork.l2);

    try {
      await provider?.send(requestMethods.switchNetwork, [
        { chainId: desiredNetwork },
      ]);

      toast({
        title: "Network Switched",
        description: `Switched to TEN Testnet`,
        variant: ToastType.SUCCESS,
      });
    } catch (error: any) {
      if (error.code === 4902) {
        toast({
          title: "Network Not Found",
          description: "Redirecting to TEN Gateway...",
          variant: ToastType.INFO,
        });
        return window.open(currentNetwork.l2Gateway, "_blank");
      } else {
        toast({
          title: "Error",
          description: error.message || "Error switching network",
          variant: ToastType.DESTRUCTIVE,
        });
      }
      handleError(error, "Error switching network");
    } finally {
      set({ loading: false });
    }
  },

  getAccounts: async (detectedProvider: ethers.providers.Web3Provider) => {
    try {
      const accounts = await detectedProvider?.send("eth_requestAccounts", []);
      if (accounts.length === 0) {
        toast({
          description: "No accounts found.",
          variant: ToastType.DESTRUCTIVE,
        });
        return [];
      }
      return accounts;
    } catch (error) {
      console.error("Error getting accounts:", error);
      throw error;
    }
  },
};
