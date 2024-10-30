/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ["@repo/ui"],
  images: {
    unoptimized: true,
  },
};

module.exports = nextConfig;
