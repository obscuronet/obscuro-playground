import { ethers } from "ethers";
import { ToastType } from "@repo/ui/lib/enums/toast";
import { showToast, toast } from "@repo/ui/components/shared/use-toast";
import { ethereum } from "@repo/ui/lib/utils";
import {
  nativeCurrency,
  tenChainIDDecimal,
  tenChainIDHex,
  tenNetworkName,
  tenscanAddress,
  userStorageAddress,
  METAMASK_CONNECTION_TIMEOUT,
  SWITCHED_CODE,
} from "@/lib/constants";
import {
  getRandomIntAsString,
  isTenChain,
  isValidTokenFormat,
} from "@/lib/utils";
import { requestMethods } from "@/routes";
import {
  accountIsAuthenticated,
  authenticateUser,
  fetchVersion,
  revokeAccountsApi,
} from "@/api/gateway";
import { Account } from "@/types/interfaces/WalletInterfaces";

const typedData = {
  types: {
    EIP712Domain: [
      { name: "name", type: "string" },
      { name: "version", type: "string" },
      { name: "chainId", type: "uint256" },
    ],
    Authentication: [{ name: "Encryption Token", type: "address" }],
  },
  primaryType: "Authentication",
  domain: {
    name: "Ten",
    version: "1.0",
    chainId: tenChainIDDecimal,
  },
  message: {
    "Encryption Token": "0x",
  },
};

const ethService = {
  initializeGateway: async (set: any, get: any) => {
    try {
      const providerInstance = new ethers.providers.Web3Provider(ethereum);
      set({ provider: providerInstance });

      await ethService.checkIfMetamaskIsLoaded(providerInstance);

      const fetchedToken = await ethService.getToken(providerInstance);
      set({ token: fetchedToken });

      const accounts = await ethService.getAccounts(providerInstance, set);
      set({ accounts: accounts || null });

      const status = await ethService.isUserConnectedToTenChain(fetchedToken);
      set({ walletConnected: status });

      const version = await fetchVersion();
      set({ version });
    } catch (error: any) {
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
      set({ walletConnected: false });
    } finally {
      set({ loading: false });
    }
  },

  checkIfMetamaskIsLoaded: async (provider: ethers.providers.Web3Provider) => {
    try {
      if (ethereum) {
        const conflictingWalletMap = {
          "Exodus Wallet": ethereum.isExodus,
          "Nest Wallet": ethereum.isNestWallet,
        };

        for (const [walletName, isWalletConnected] of Object.entries(
          conflictingWalletMap
        )) {
          if (isWalletConnected) {
            const message = `${walletName} is connected and is conflicting with MetaMask. Please disable ${walletName} and try again.`;
            showToast(ToastType.DESTRUCTIVE, message);
            throw new Error(message);
          }
        }

        return await ethService.handleEthereum(provider);
      } else {
        showToast(ToastType.INFO, "Connecting to MetaMask...");

        let timeoutId: ReturnType<typeof setTimeout>;

        const handleEthereumOnce = async () => {
          await ethService.handleEthereum(provider);
        };

        window.addEventListener(
          "ethereum#initialized",
          () => {
            clearTimeout(timeoutId);
            handleEthereumOnce();
          },
          { once: true }
        );

        timeoutId = setTimeout(() => {
          handleEthereumOnce();
        }, METAMASK_CONNECTION_TIMEOUT);
      }
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
      throw error;
    }
  },

  handleEthereum: async (provider: ethers.providers.Web3Provider) => {
    try {
      if (ethereum && ethereum.isMetaMask) {
        return;
      } else {
        return showToast(
          ToastType.WARNING,
          "Please install MetaMask to use TEN Gateway."
        );
      }
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
      throw error;
    }
  },

  isUserConnectedToTenChain: async (token: string) => {
    if (await isTenChain()) {
      return !!(token && isValidTokenFormat(token));
    }
    return false;
  },

  getAccounts: async (provider: ethers.providers.Web3Provider, set: any) => {
    if (!provider) {
      showToast(
        ToastType.DESTRUCTIVE,
        "No provider found. Please try again later."
      );
      return;
    }

    const token = await ethService.getToken(provider);

    if (!token || !isValidTokenFormat(token)) {
      set({ walletConnected: false });
      return;
    }

    try {
      showToast(ToastType.INFO, "Getting accounts...");

      if (!(await isTenChain())) {
        showToast(ToastType.DESTRUCTIVE, "Please connect to the TEN chain.");
        set({ walletConnected: false });
        return;
      }

      const accounts = await provider.listAccounts();

      if (accounts.length === 0) {
        showToast(ToastType.DESTRUCTIVE, "No MetaMask accounts found.");
        set({ walletConnected: false });
        return [];
      }
      showToast(ToastType.SUCCESS, "Accounts found!");

      set({ token, walletConnected: true });

      const authenticatedAccounts = await ethService.authenticateAccounts(
        accounts,
        provider,
        token
      );

      return authenticatedAccounts;
    } catch (error) {
      console.error(error);
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
      set({ walletConnected: false });
      throw error;
    }
  },

  authenticateAccounts: async (
    accounts: string[],
    provider: ethers.providers.Web3Provider,
    token: string
  ) => {
    if (!provider) {
      showToast(
        ToastType.DESTRUCTIVE,
        "No provider found. Please try again later."
      );
      return [];
    }
    showToast(ToastType.INFO, "Authenticating accounts...");
    const updatedAccounts = await Promise.all(
      accounts.map(async (account) => {
        const isAuthenticated = await accountIsAuthenticated(token, account);

        if (!isAuthenticated.status) {
          const authenticated = await ethService.authenticateWithGateway(
            token,
            account
          );

          return {
            name: account,
            connected: authenticated,
          };
        }

        return {
          name: account,
          connected: isAuthenticated.status,
        };
      })
    );

    showToast(ToastType.SUCCESS, "Accounts authenticated!");
    return updatedAccounts;
  },

  connectAccount: async (set: any, get: any, account: string) => {
    const { token } = get();

    try {
      if (!token) {
        showToast(
          ToastType.INFO,
          "Encryption token is required to connect an account."
        );
        return;
      }

      const authenticated = await ethService.authenticateWithGateway(
        token,
        account
      );

      if (authenticated) {
        showToast(ToastType.SUCCESS, "Account authenticated!");

        set((state: any) => ({
          accounts:
            state.accounts?.map((acc: Account) =>
              acc.name === account ? { ...acc, connected: true } : acc
            ) || null,
          walletConnected: true,
        }));
      } else {
        showToast(ToastType.DESTRUCTIVE, "Account authentication failed.");
      }
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, "Account authentication failed.");
    }
  },

  authenticateWithGateway: async (token: string, account: string) => {
    try {
      return await ethService.authenticateAccountWithTenGatewayEIP712(
        token,
        account
      );
    } catch (error) {
      showToast(
        ToastType.DESTRUCTIVE,
        `Error authenticating account: ${account}`
      );
      return false;
    }
  },

  switchToTenNetwork: async () => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      await ethereum.request({
        method: requestMethods.switchNetwork,
        params: [{ chainId: tenChainIDHex }],
      });
      return 0;
    } catch (switchError: any) {
      showToast(
        ToastType.DESTRUCTIVE,
        `switchToTenNetwork: ${switchError.code}`
      );
      return switchError.code;
    }
  },

  connectAccounts: async () => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      return await ethereum.request({
        method: requestMethods.connectAccounts,
      });
    } catch (error) {
      console.error(error);
      throw error;
    }
  },

  getSignature: async (account: string, data: any) => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      return await ethereum.request({
        method: requestMethods.signTypedData,
        params: [account, JSON.stringify(data)],
      });
    } catch (error) {
      console.error(error);
      throw "Failed to get signature";
    }
  },

  getToken: async (provider: ethers.providers.Web3Provider) => {
    if (!provider?.send) {
      return null;
    }
    try {
      if (await isTenChain()) {
        const token = await provider.send(requestMethods.getStorageAt, [
          userStorageAddress,
          getRandomIntAsString(0, 1000),
          null,
        ]);
        return token;
      } else {
        return null;
      }
    } catch (e: any) {
      toast({
        title: "Invalid Encrypted Token",
        variant: ToastType.DESTRUCTIVE,
        description:
          e instanceof Error
            ? e.message
            : e?.data?.message?.includes("not found")
              ? "Please restart the process to get a new encryption token by removing TEN Testnet from your wallet and reconnecting."
              : "An error occurred. Please try again.}",
      });
      console.error(e);
      throw e;
    }
  },

  addNetworkToMetaMask: async (rpcUrls: string[]) => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      await ethereum.request({
        method: requestMethods.addNetwork,
        params: [
          {
            chainId: tenChainIDHex,
            chainName: tenNetworkName,
            nativeCurrency,
            rpcUrls,
            blockExplorerUrls: [tenscanAddress],
          },
        ],
      });
      return true;
    } catch (error) {
      console.error(error);
      throw error;
    }
  },

  authenticateAccountWithTenGatewayEIP712: async (
    token: string,
    account: string
  ): Promise<any> => {
    if (!token) {
      return showToast(
        ToastType.INFO,
        "Encryption token not found. Please try again later."
      );
    }

    try {
      const isAuthenticated = await accountIsAuthenticated(token, account);
      if (isAuthenticated.status) {
        return {
          status: true,
          message: "Account already authenticated",
        };
      }
      const data = {
        ...typedData,
        message: {
          ...typedData.message,
          "Encryption Token": token,
        },
      };
      const signature = await ethService.getSignature(account, data);

      return await authenticateUser(token, {
        signature,
        address: account,
      });
    } catch (error: any) {
      throw error;
    }
  },

  revokeAccounts: async (set: any, get: any) => {
    const { token } = get();
    if (!token) {
      showToast(
        ToastType.INFO,
        "Encryption token is required to revoke accounts"
      );
      return;
    }
    const revokeResponse = await revokeAccountsApi(token);
    if (revokeResponse === ToastType.SUCCESS) {
      showToast(ToastType.DESTRUCTIVE, "Accounts revoked!");
      set({
        accounts: null,
        walletConnected: false,
        token: "",
        loading: false,
      });
    }
  },
};

export default ethService;
