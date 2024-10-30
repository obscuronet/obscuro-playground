import React from "react";
import { downloadMetaMask, ethereum } from "@repo/ui/lib/utils";
import { ConnectWalletButtonProps } from "@repo/ui/lib/interfaces/ui";
import {
  ExclamationTriangleIcon,
  Link2Icon,
  LinkBreak2Icon,
} from "@repo/ui/components/shared/react-icons";
import ConnectWalletButton from "@repo/ui/components/common/connect-wallet";
import useWalletStore from "@/stores/wallet-store";
import { isValidTokenFormat } from "@/lib/utils";
import useGatewayService from "@/services/useGatewayService";

const GatewayConnectButton = (
  props: Omit<ConnectWalletButtonProps, "onConnect" | "renderContent">
) => {
  const {
    walletConnected,
    disconnectWallet,
    accounts,
    token,
    isWrongNetwork,
    switchNetwork,
  } = useWalletStore();

  const { connectToTenTestnet } = useGatewayService();

  const handleConnect = () => {
    if (!ethereum) {
      downloadMetaMask();
      return;
    }

    if (!walletConnected) {
      connectToTenTestnet();
      return;
    }

    if (isWrongNetwork) {
      if (token && isValidTokenFormat(token)) {
        switchNetwork();
      } else {
        disconnectWallet();
      }
      return;
    }

    disconnectWallet();
  };

  const renderGatewayContent = () => {
    if (!ethereum) {
      return (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          Install MetaMask
        </>
      );
    }

    if (isWrongNetwork) {
      return (
        <>
          <ExclamationTriangleIcon className="h-4 w-4 mr-1 text-yellow-500" />
          Unsupported network
        </>
      );
    }

    if (walletConnected) {
      const connectedAccounts = accounts?.filter((acc) => acc.connected) || [];
      const totalAccounts = accounts?.length || 0;
      return (
        <>
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          {`${connectedAccounts.length}/${totalAccounts} Connected`}
        </>
      );
    } else {
      return (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          Connect to TEN Testnet
        </>
      );
    }
  };

  return (
    <ConnectWalletButton
      {...props}
      onConnect={handleConnect}
      renderContent={renderGatewayContent}
    />
  );
};

export default GatewayConnectButton;
