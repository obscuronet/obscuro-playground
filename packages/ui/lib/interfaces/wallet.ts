import { ethers } from "ethers";

export interface IChain {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
  chainId?: string;
}

export interface IToken {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
  address: string;
}

export interface INetworkConfig {
  ManagementContractAddress: string;
  L1StartHash: string;
  MessageBusAddress: string;
  L2MessageBusAddress: string;
  ImportantContracts: {
    L1Bridge: string;
    L1CrossChainMessenger: string;
    L2Bridge: string;
    L2CrossChainMessenger: string;
  };
}

export interface IContractState {
  bridgeContract: ethers.Contract | null;
  managementContract: ethers.Contract | null;
  messageBusContract: ethers.Contract | null;
  messageBusAddress: string;
  bridgeAddress: string;
  setContractState: (state: Partial<IContractState>) => void;
}

export interface IWalletState {
  provider: ethers.providers.Web3Provider | null;
  signer: ethers.Signer | null;
  address: string;
  chainId: number | null;
  walletConnected: boolean;
  isWrongNetwork: boolean;
  loading: boolean;
  initializeProvider: () => void;
  connectWallet?: () => void;
  disconnectWallet: () => void;
  switchNetwork: () => void;
  restoreWalletState?: () => void;
}
