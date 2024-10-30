import Link from "next/link";
import { socialLinks } from "../../lib/constants";
import {
  GitHubLogoIcon,
  TwitterLogoIcon,
  DiscordLogoIcon,
} from "@repo/ui/components/shared/react-icons";
import useWalletStore from "@/stores/wallet-store";

const SOCIAL_LINKS = [
  {
    name: "GitHub",
    href: socialLinks.github,
    icon: GitHubLogoIcon,
  },
  {
    name: "Twitter",
    href: socialLinks.twitter,
    icon: TwitterLogoIcon,
  },
  {
    name: "Discord",
    href: socialLinks.discord,
    icon: DiscordLogoIcon,
  },
];

export default function Footer() {
  const { version } = useWalletStore();

  return (
    <div className="border-t p-2">
      <div className="flex h-16 items-center justify-between px-4 flex-wrap">
        <div className="flex items-center space-x-4 pr-2">
          {SOCIAL_LINKS.map((item, index) => (
            <a
              key={item.name}
              href={item.href}
              aria-label={item.name}
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-primary transition-colors"
            >
              <item.icon />
            </a>
          ))}
        </div>
        <div className="flex items-center justify-center space-x-4 pr-2">
          <h3 className="text-xs text-muted-foreground">
            Version: {version || "Unknown"}
          </h3>
        </div>
        <div className="flex items-center space-x-4">
          <Link
            href="/docs/privacy"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            Privacy
          </Link>
          <Link
            href="/docs/terms"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
          >
            Terms
          </Link>
        </div>
      </div>
    </div>
  );
}
