import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatTimeAgo(unixTimestampSeconds: string) {
  const date = new Date(Number(unixTimestampSeconds) * 1000);
  return formatDistanceToNow(date, { addSuffix: true });
}

export const { ethereum } =
  typeof window !== "undefined" ? window : ({} as any);

export const downloadMetaMask = () => {
  window ? window.open("https://metamask.io/download", "_blank") : null;
};
