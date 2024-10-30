import { showToast, toast } from "@repo/ui/components/shared/use-toast";
import { joinTestnet } from "../api/gateway";
import {
  SWITCHED_CODE,
  tenGatewayAddress,
  tenGatewayVersion,
} from "../lib/constants";
import { isTenChain, isValidTokenFormat } from "../lib/utils";
import { fetchTestnetStatus } from "@/api/general";
import { ToastType } from "@repo/ui/lib/enums/toast";
import { ethers } from "ethers";
import useWalletStore from "@/stores/wallet-store";
import ethService from "./ethService";

const useGatewayService = () => {
  const { token, provider, fetchUserAccounts, loading, setLoading } =
    useWalletStore();

  const isMetamaskConnected = async () => {
    if (!provider) {
      return false;
    }
    try {
      const accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to get accounts");
      throw error;
    }
  };

  const connectToTenTestnet = async () => {
    showToast(ToastType.INFO, "Connecting to TEN Testnet...");
    setLoading(true);
    try {
      if (await isTenChain()) {
        if (!token || !isValidTokenFormat(token)) {
          toast({
            title: "Invalid Encrypted Token",
            variant: ToastType.DESTRUCTIVE,
            description:
              "Please restart the process to get a new encryption token by removing TEN Testnet from your wallet and reconnecting.",
          });
          return;
        }
      }
      showToast(ToastType.INFO, "Switching to TEN Testnet...");
      const switched = await ethService.switchToTenNetwork();
      showToast(ToastType.SUCCESS, `Switched to TEN Testnet`);
      // SWITCHED_CODE=4902; error 4902 means that the chain does not exist
      if (
        switched === SWITCHED_CODE ||
        !isValidTokenFormat(
          await ethService.getToken(provider as ethers.providers.Web3Provider)
        )
      ) {
        showToast(ToastType.INFO, "Adding TEN Testnet...");
        const user = await joinTestnet();
        const rpcUrls = [
          `${tenGatewayAddress}/${tenGatewayVersion}/?token=${user}`,
        ];
        await ethService.addNetworkToMetaMask(rpcUrls);
        showToast(ToastType.SUCCESS, "Added TEN Testnet");
      }

      if (!(await isMetamaskConnected())) {
        showToast(ToastType.INFO, "No accounts found, connecting...");
        await ethService.connectAccounts();
        showToast(ToastType.SUCCESS, "Connected to TEN Testnet");
      }
      await fetchUserAccounts();
    } catch (error: Error | any) {
      toast({
        title: "Invalid Encrypted Token",
        variant: ToastType.DESTRUCTIVE,
        description:
          error instanceof Error
            ? error.message
            : error?.data?.message?.includes("not found")
              ? "Please restart the process to get a new encryption token by removing TEN Testnet from your wallet and reconnecting."
              : "An error occurred. Please try again.}",
      });
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const getTestnetStatus = async () => {
    try {
      return await fetchTestnetStatus();
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to connect to TEN Testnet");
      throw error;
    }
  };

  return {
    connectToTenTestnet,
    getTestnetStatus,
    isMetamaskConnected,
    loading,
  };
};

export default useGatewayService;
