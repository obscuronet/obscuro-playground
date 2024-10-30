import { IWalletState } from "@repo/ui/lib/interfaces/wallet";
import { Account } from "./WalletInterfaces";

export type NavLink = {
  label: string;
  href?: string;
  isDropdown?: boolean;
  isExternal?: boolean;
  subNavLinks?: NavLink[];
};

export type StoreSet = (
  partial:
    | IGatewayWalletState
    | Partial<IGatewayWalletState>
    | ((
        state: IGatewayWalletState
      ) => IGatewayWalletState | Partial<IGatewayWalletState>),
  replace?: boolean | undefined
) => void;

export type StoreGet = () => IGatewayWalletState;

export interface IGatewayWalletState extends IWalletState {
  token: string;
  version: string | null;
  accounts: Account[] | null;
  setLoading: (loading: boolean) => void;
  initializeGateway: () => Promise<void>;
  connectAccount: (account: string) => Promise<void>;
  revokeAccounts: () => Promise<void>;
  fetchUserAccounts: () => Promise<void>;
}
