import { ethers } from "ethers";
import { currentNetwork, ethereum, handleStorage } from ".";
import { toast } from "../../components/shared/use-toast";
import { ToastType } from "../enums/toast";

export const setupEventListeners = (
  updateState: (
    address: string,
    chainId: number,
    isWrongNetwork: boolean
  ) => void
) => {
  const expectedChainId = currentNetwork.l2;

  const handleAccountsChange = (accounts: string[]) => {
    const chainId = parseInt(ethereum.chainId, 16);
    const isWrongNetwork = chainId !== expectedChainId;

    if (accounts.length > 0) {
      const address = accounts[0] as string;
      updateState(address, chainId, isWrongNetwork);
      handleStorage.save("walletAddress", address);
    } else {
      updateState("", chainId, isWrongNetwork);
      handleStorage.remove("walletAddress");
      toast({
        title: "Disconnected",
        description: "No accounts found. Please reconnect your wallet.",
        variant: ToastType.INFO,
      });
    }
  };

  const handleChainChange = async (chainIdHex: string) => {
    const chainId = parseInt(chainIdHex, 16);
    const isWrongNetwork = chainId !== expectedChainId;
    const currentAddress = ethereum.selectedAddress || "";
    updateState(currentAddress, chainId, isWrongNetwork);

    if (isWrongNetwork) {
      toast({
        title: "Wrong Network",
        description: "Please switch to the correct network.",
        variant: ToastType.INFO,
      });
    } else {
      toast({
        title: "Network Changed",
        description: `Switched to the correct network (Chain ID: ${chainId}).`,
        variant: ToastType.SUCCESS,
      });
    }
  };

  ethereum.on("accountsChanged", handleAccountsChange);
  ethereum.on("chainChanged", handleChainChange);

  return () => {
    ethereum.removeListener("accountsChanged", handleAccountsChange);
    ethereum.removeListener("chainChanged", handleChainChange);
  };
};

export const initializeSigner = async (
  provider: ethers.providers.Web3Provider
) => {
  return provider.getSigner() as ethers.Signer;
};
