import { formatDistanceToNow } from "date-fns";
import { tenChainIDHex, tokenHexLength } from "./constants";
import { ethereum } from "@repo/ui/lib/utils";

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

export async function isTenChain() {
  if (!ethereum) {
    return false;
  }
  let currentChain = await ethereum.request({
    method: "eth_chainId",
  });
  return currentChain === tenChainIDHex;
}
