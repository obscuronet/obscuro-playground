import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";

import { Button } from "@repo/ui/components/shared/button";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
} from "@repo/ui/components/shared/dropdown-menu";

import { ChevronDownIcon } from "@repo/ui/components/shared/react-icons";
import { NavLink } from "../types/interfaces";
import { NavLinks } from "../routes";
import { cn } from "@repo/ui/lib/utils";

const NavItem = ({ navLink }: { navLink: NavLink }) => {
  const router = useRouter();

  const isDropdownActive = (navLink: NavLink) => {
    return navLink.subNavLinks?.some(
      (subNavLink: NavLink) =>
        subNavLink.href && router.pathname.includes(subNavLink.href)
    );
  };
  if (navLink.isDropdown) {
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="clear"
            size={"wrap"}
            className={cn(
              "text-sm font-medium text-muted-foreground transition-colors hover:text-primary flex justify-start",
              {
                "text-primary": isDropdownActive(navLink),
              }
            )}
          >
            {navLink.label} <ChevronDownIcon className="ml-1 h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end" forceMount>
          <DropdownMenuGroup>
            {navLink.subNavLinks &&
              navLink.subNavLinks.map((subNavLink: NavLink) => (
                <DropdownMenuItem key={subNavLink.label}>
                  <NavItem navLink={subNavLink} />
                </DropdownMenuItem>
              ))}
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  } else if (navLink.isExternal) {
    return (
      <a
        href={navLink.href}
        key={navLink.label}
        className="text-sm font-medium transition-colors hover:text-primary"
      >
        {navLink.label}
      </a>
    );
  } else {
    return (
      <Link
        href={navLink.href || ""}
        key={navLink.label}
        className={cn(
          "text-sm font-medium text-muted-foreground transition-colors hover:text-primary",
          {
            "text-primary": router.pathname === navLink.href,
          }
        )}
      >
        {navLink.label}
      </Link>
    );
  }
};

export const MainNav = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) => (
  <nav className={cn("flex items-center lg:space-x-6", className)} {...props}>
    {NavLinks.map((navLink) => (
      <div key={navLink.label} className="w-full flex items-center mr-4 p-2">
        <NavItem navLink={navLink} />
      </div>
    ))}
  </nav>
);
