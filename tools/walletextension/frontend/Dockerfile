# Use an official Node.js LTS version as a base image
FROM node:18-buster AS base

# Install pnpm
RUN npm install -g pnpm

# Set the working directory
WORKDIR /home/obscuro/go-obscuro/

# Copy the entire project
COPY . .

# Install dependencies for the entire project
RUN pnpm install

# Set the working directory for the frontend
WORKDIR /home/obscuro/go-obscuro/tools/walletextension/frontend

# Install dependencies specific to the frontend
RUN pnpm install --filter ./tools/walletextension/frontend...

# Build the Next.js app (if needed)
# RUN pnpm run build

# Set environment variables
ARG NEXT_PUBLIC_NETWORK_NAME
ARG NEXT_PUBLIC_TENSCAN_URL
ARG NEXT_PUBLIC_GATEWAY_URL
ARG NEXT_PUBLIC_API_HOST_ENVIRONMENT
ENV NEXT_PUBLIC_NETWORK_NAME=$NEXT_PUBLIC_NETWORK_NAME
ENV NEXT_PUBLIC_TENSCAN_URL=$NEXT_PUBLIC_TENSCAN_URL
ENV NEXT_PUBLIC_GATEWAY_URL=$NEXT_PUBLIC_GATEWAY_URL
ENV NEXT_PUBLIC_API_HOST_ENVIRONMENT=${NEXT_PUBLIC_API_HOST_ENVIRONMENT}

# Expose the port
EXPOSE 80

# Start the application
CMD ["pnpm", "run", "start"]