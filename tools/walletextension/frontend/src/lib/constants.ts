export const tenGatewayAddress =
  process.env.NEXT_PUBLIC_API_GATEWAY_URL || "http://127.0.0.1:1443";

export const tenscanLink = "https://testnet.tenscan.io";

export const socialLinks = {
  faucet: "https://faucet.ten.xyz",
  github: "https://github.com/ten-protocol",
  discord: "https://discord.gg/tenprotocol",
  twitter: "https://twitter.com/tenprotocol",
  twitterHandle: "@tenprotocol",
};

export const GOOGLE_ANALYTICS_ID = process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID;

export const testnetUrls = {
  sepolia: {
    name: "Ten Testnet",
    url: "https://sepolia-testnet.ten.xyz",
    rpc: "https://rpc.sepolia-testnet.ten.xyz",
  },
  uat: {
    name: "Ten UAT-Testnet",
    url: "https://uat-testnet.ten.xyz",
    rpc: "https://rpc.uat-testnet.ten.xyz",
  },
  dev: {
    name: "Ten Dev-Testnet",
    url: "https://dev-testnet.ten.xyz",
    rpc: "https://rpc.dev-testnet.ten.xyz",
  },
  default: {
    name: "Ten Testnet",
    url: tenGatewayAddress,
  },
};

export const SWITCHED_CODE = 4902;
export const tokenHexLength = 42;

export const tenGatewayVersion = "v1";
export const tenChainIDDecimal = 443;

export const tenChainIDHex = "0x" + tenChainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'
export const METAMASK_CONNECTION_TIMEOUT = 3000;

export const userStorageAddress = "0x0000000000000000000000000000000000000001";

export const nativeCurrency = {
  name: "Sepolia Ether",
  symbol: "ETH",
  decimals: 18,
};
