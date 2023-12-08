import { NavLink } from "../types/interfaces";

export const apiRoutes = {
  // **** BATCHES ****
  getLatestBatch: "/items/batch/latest/",
  getBatches: "/items/batches/",
  getBatchByHash: "/items/batch/:hash",

  // **** BLOCKS ****
  getBlocks: "/items/blocks/",

  // **** CONTRACTS ****
  getContractCount: "/count/contracts/",
  getVerifiedContracts: "/info/obscuro/",

  // **** TRANSACTIONS ****
  getTransactions: "/items/transactions/",
  getTransactionCount: "/count/transactions/",

  getEtherPrice:
    "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd",

  // **** ROLLUPS ****
  getRollups: "/items/rollup/latest/",
  decryptEncryptedRollup: "/actions/decryptTxBlob/",
};

export const ethMethods = {
  getStorageAt: "eth_getStorageAt",
};

export const NavLinks: NavLink[] = [
  {
    href: "/",
    label: "Home",
    isExternal: false,
    isDropdown: false,
  },
  {
    href: "/personal",
    label: "Personal",
    isExternal: false,
    isDropdown: false,
  },
  {
    label: "Blockchain",
    isExternal: false,
    isDropdown: true,
    subNavLinks: [
      {
        href: "/transactions",
        label: "Transactions",
        isExternal: false,
      },
      {
        href: "/blocks",
        label: "Blocks",
        isExternal: false,
      },
      {
        href: "/batches",
        label: "Batches",
        isExternal: false,
      },
    ],
  },
  {
    label: "Resources",
    isExternal: false,
    isDropdown: true,
    subNavLinks: [
      {
        href: "/resources/decrypt",
        label: "Decrypt",
        isExternal: false,
      },
      {
        href: "/resources/verified-data",
        label: "Verified Data",
        isExternal: false,
      },
    ],
  },
];

export const externalLinks = {
  etherscanBlock: "https://etherscan.io/block/",
};
