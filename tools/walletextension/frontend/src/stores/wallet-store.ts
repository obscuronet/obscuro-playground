import { create } from "zustand";
import ethService from "@/services/ethService";
import { walletService } from "@repo/ui/services/walletService";
import { IGatewayWalletState } from "@/types/interfaces";
import { IWalletState } from "@repo/ui/lib/interfaces/wallet";

interface WalletState extends IGatewayWalletState, IWalletState {}

const useWalletStore = create<WalletState>((set, get) => ({
  // gateway-specific state
  token: "",
  version: null,
  accounts: null,
  loading: true,

  // common state (from monorepo)
  provider: null,
  signer: null,
  address: "",
  chainId: null,
  walletConnected: false,
  isWrongNetwork: false,

  // gateway-specific methods
  initializeGateway: async () => {
    await ethService.initializeGateway(set, get);
  },

  connectAccount: async (account: string) => {
    await ethService.connectAccount(set, get, account);
  },

  revokeAccounts: async () => {
    await ethService.revokeAccounts(set, get);
  },

  setLoading: (loading: boolean) => {
    set({ loading });
  },

  fetchUserAccounts: async () => {
    const { provider } = get();
    if (!provider) {
      return;
    }
    const accounts = await ethService.getAccounts(provider, set);
    set({ accounts });
  },

  // monorepo methods
  //@ts-ignore
  initializeProvider: () => walletService.initializeProvider(set, get),

  disconnectWallet: async () => {
    //@ts-ignore
    await walletService.disconnectWallet(set, get);
    get().revokeAccounts();
  },

  //@ts-ignore
  switchNetwork: () => walletService.switchNetwork(set, get),
}));

export default useWalletStore;
