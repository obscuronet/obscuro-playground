import {
  authenticateAccountWithTenGatewayEIP712,
  getUserID,
} from "@/api/ethRequests";
import { accountIsAuthenticated } from "@/api/gateway";
import { showToast } from "@/components/ui/use-toast";
import { METAMASK_CONNECTION_TIMEOUT } from "@/lib/constants";
import { isTenChain, isValidUserIDFormat } from "@/lib/utils";
import { ToastType } from "@/types/interfaces";
import { Account } from "@/types/interfaces/WalletInterfaces";
import { ethers } from "ethers";

const { ethereum } = typeof window !== "undefined" ? window : ({} as any);

const ethService = {
  checkIfMetamaskIsLoaded: async (provider: ethers.providers.Web3Provider) => {
    if (ethereum) {
      await ethService.handleEthereum(provider);
    } else {
      showToast(ToastType.INFO, "Connecting to MetaMask...");

      const handleEthereumOnce = () => {
        ethService.handleEthereum(provider);
      };

      window.addEventListener("ethereum#initialized", handleEthereumOnce, {
        once: true,
      });

      setTimeout(() => {
        handleEthereumOnce(); // Call the handler function after the timeout
      }, METAMASK_CONNECTION_TIMEOUT);
    }
  },

  handleEthereum: async (provider: ethers.providers.Web3Provider) => {
    if (ethereum && ethereum.isMetaMask) {
      const fetchedUserID = await getUserID(provider);
      if (fetchedUserID && isValidUserIDFormat(fetchedUserID)) {
        showToast(ToastType.SUCCESS, "MetaMask connected!");
      } else {
        showToast(
          ToastType.WARNING,
          "Please connect to the Ten chain to use Ten Gateway."
        );
      }
    } else {
      showToast(
        ToastType.WARNING,
        "Please install MetaMask to use Ten Gateway."
      );
    }
  },

  fetchUserID: async (provider: ethers.providers.Web3Provider) => {
    try {
      return await getUserID(provider);
    } catch (e: any) {
      showToast(
        ToastType.DESTRUCTIVE,
        `${e.message} ${e.data?.message}` ||
          "Error: Could not fetch your user ID. Please try again later."
      );
      return null;
    }
  },

  getCorrectScreenBasedOnMetamaskAndUserID: async (userID: string) => {
    if (await isTenChain()) {
      if (userID && isValidUserIDFormat(userID)) {
        return true;
      } else {
        return false;
      }
    } else {
      return false;
    }
  },

  getAccounts: async (provider: ethers.providers.Web3Provider) => {
    const id = await getUserID(provider);
    if (!id || !isValidUserIDFormat(id)) {
      return;
    }

    try {
      if (!provider) {
        showToast(
          ToastType.DESTRUCTIVE,
          "No provider found. Please try again later."
        );
        return;
      }

      showToast(ToastType.INFO, "Getting accounts...");

      if (!(await isTenChain())) {
        showToast(ToastType.DESTRUCTIVE, "Please connect to the Ten chain.");
        return;
      }

      const accounts = await provider.listAccounts();

      if (accounts.length === 0) {
        showToast(ToastType.DESTRUCTIVE, "No MetaMask accounts found.");
        return [];
      }

      let updatedAccounts: Account[] = [];

      for (let i = 0; i < accounts.length; i++) {
        const account = accounts[i];
        await authenticateAccountWithTenGatewayEIP712(id, account);
        const { status } = await accountIsAuthenticated(id, account);
        updatedAccounts.push({
          name: account,
          connected: status,
        });
      }
      showToast(ToastType.SUCCESS, "Accounts fetched successfully.");
      return updatedAccounts;
    } catch (error) {
      console.error(error);
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
    }
  },
};

export default ethService;
