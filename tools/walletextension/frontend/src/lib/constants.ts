export const tenGatewayAddress = "https://uat-testnet.obscu.ro";
export const tenscanLink = "https://testnet.tenscan.com";

export const socialLinks = {
  github: "https://github.com/obscuronet",
  discord: "https://discord.gg/2JQ2Z3r",
  twitter: "https://twitter.com/obscuronet",
};

export const testnetUrls = {
  sepolia: {
    name: "Ten Dev-Testnet",
    url: "https://sepolia-testnet.obscu.ro",
    rpc: "https://rpc.sepolia-testnet.obscu.ro",
  },
  uat: {
    name: "Ten UAT-Testnet",
    url: "https://uat-testnet.obscu.ro",
    rpc: "https://rpc.uat-testnet.obscu.ro",
  },
  dev: {
    name: "Ten Dev-Testnet",
    url: "https://dev-testnet.obscu.ro",
    rpc: "https://rpc.dev-testnet.obscu.ro",
  },
  default: {
    name: "Ten Testnet",
    url: tenGatewayAddress,
  },
};

export const SWITCHED_CODE = 4902;
export const userIDHexLength = 40;

export const tenGatewayVersion = "v1";
export const tenChainIDDecimal = 443;

export const tenChainIDHex = "0x" + tenChainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'
export const METAMASK_CONNECTION_TIMEOUT = 3000;

export const nativeCurrency = {
  name: "Sepolia Ether",
  symbol: "ETH",
  decimals: 18,
};
