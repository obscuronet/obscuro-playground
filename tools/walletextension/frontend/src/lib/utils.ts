import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";
import {
  tenChainIDHex,
  tenGatewayAddress,
  testnetUrls,
  tokenHexLength,
} from "./constants";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatTimeAgo(unixTimestampSeconds: string) {
  const date = new Date(Number(unixTimestampSeconds) * 1000);
  return formatDistanceToNow(date, { addSuffix: true });
}

export function isValidTokenFormat(value: string) {
  return typeof value === "string" && value.length === tokenHexLength;
}

export function getRandomIntAsString(min: number, max: number) {
  min = Math.ceil(min);
  max = Math.floor(max);
  const randomInt = Math.floor(Math.random() * (max - min + 1)) + min;
  return randomInt.toString();
}

export function getNetworkName() {
  switch (tenGatewayAddress) {
    case testnetUrls.uat.url:
      return testnetUrls.uat.name;
    case testnetUrls.sepolia.url:
      return testnetUrls.sepolia.name;
    default:
      return testnetUrls.default.name;
  }
}

export function getRPCFromUrl() {
  // get the correct RPC endpoint for each network
  switch (testnetUrls.default.url) {
    // case 'https://testnet.obscu.ro':
    //     return 'https://rpc.sepolia-testnet.obscu.ro'
    case testnetUrls.sepolia.url:
      return testnetUrls.sepolia.rpc;
    case testnetUrls.uat.url:
      return testnetUrls.uat.rpc;
    case testnetUrls.dev.url:
      return testnetUrls.dev.rpc;
    default:
      return testnetUrls.default.url;
  }
}

export async function isTenChain() {
  if (!ethereum) {
    return false;
  }
  let currentChain = await ethereum.request({
    method: "eth_chainId",
  });
  return currentChain === tenChainIDHex;
}

export const { ethereum } =
  typeof window !== "undefined" ? window : ({} as any);

export const downloadMetaMask = () => {
  window ? window.open("https://metamask.io/download", "_blank") : null;
};
