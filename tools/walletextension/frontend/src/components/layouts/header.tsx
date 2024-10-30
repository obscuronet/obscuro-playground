import { MainNav } from "../main-nav";
import { ModeToggle } from "../mode-toggle";
import Link from "next/link";
import { HamburgerMenuIcon } from "@repo/ui/components/shared/react-icons";
import { useState } from "react";
import { Button } from "@repo/ui/components/shared/button";
import HealthIndicator from "../health-indicator";
import Image from "next/image";
import useWalletStore from "@/stores/wallet-store";
import GatewayConnectButton from "../modules/connect-wallet/gateway-connect-button";

export default function Header() {
  const { walletConnected } = useWalletStore();
  return (
    <div className="border-b">
      <div className="flex h-16 justify-between items-center px-4">
        <Link href="/">
          <Image
            src="/assets/images/black_logotype.png"
            alt="Logo"
            width={150}
            height={40}
            className="cursor-pointer dark:hidden"
          />
          <Image
            src="/assets/images/white_logotype.png"
            alt="Logo"
            width={150}
            height={40}
            className="cursor-pointer hidden dark:block"
          />
        </Link>
        <div className="hidden md:flex items-center space-x-4">
          <MainNav className="mx-6" />
          <div className="flex items-center space-x-4">
            <HealthIndicator />
            <ModeToggle />
            <GatewayConnectButton />
          </div>
        </div>
        <div className="flex items-center space-x-4 md:hidden">
          <MobileMenu />
        </div>
      </div>
    </div>
  );
}

const MobileMenu = () => {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="relative">
      <ModeToggle />
      <Button
        variant={"clear"}
        className="text-muted-foreground hover:text-primary transition-colors"
        onClick={() => setMenuOpen(!menuOpen)}
      >
        <HamburgerMenuIcon />
      </Button>
      {menuOpen && (
        <div className="absolute z-10 top-0 right-0 mt-12">
          <div className="bg-background border rounded-lg shadow-lg">
            <div className="flex flex-col p-4 space-y-2">
              <MainNav className="flex flex-col" />
              <GatewayConnectButton />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
