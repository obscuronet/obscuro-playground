import Link from "next/link";
import { socialLinks } from "../../lib/constants";
import {
  GitHubLogoIcon,
  TwitterLogoIcon,
  DiscordLogoIcon,
} from "@radix-ui/react-icons";

export default function Footer() {
  return (
    <div className="border-t px-2">
      <div className="flex h-16 items-center px-4">
        <div className="flex-1 flex items-center space-x-4">
          <a
            href={socialLinks.github}
            aria-label="GitHub"
            className="text-muted-foreground hover:text-primary transition-colors"
          >
            <GitHubLogoIcon />
          </a>
          <a
            href={socialLinks.twitter}
            aria-label="Twitter"
            className="text-muted-foreground hover:text-primary transition-colors"
          >
            <TwitterLogoIcon />
          </a>
          <a
            href={socialLinks.discord}
            aria-label="Discord"
            className="text-muted-foreground hover:text-primary transition-colors"
          >
            <DiscordLogoIcon />
          </a>
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