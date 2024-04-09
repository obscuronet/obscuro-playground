# Build Stages:
# build-base = downloads modules and prepares the directory for compilation. Based on the ego-dev image
# build-enclave = copies over the actual source code of the project and builds it using a compiler cache
# deploy = copies over only the enclave executable without the source
#          in a lightweight base image specialized for deployment and prepares the /data/ folder.

# Final container folder structure:
#   /home/ten/data                          contains working files for the enclave
#   /home/ten/go-ten/go/enclave/main    contains the executable for the enclave
#

# Defaults to restricted flag mode
ARG TESTMODE=false

FROM ghcr.io/edgelesssys/ego-dev:v1.3.0 AS build-base

# setup container data structure
RUN mkdir -p /home/ten/go-ten

# Ensures container layer caching when dependencies are not changed
WORKDIR /home/ten/go-ten
COPY go.mod .
COPY go.sum .
RUN ego-go mod download


# Trigger new build stage for compiling the enclave
FROM build-base as build-enclave
COPY . .

WORKDIR /home/ten/go-ten/go/enclave/main

# Build the enclave using the cross image build cache.
RUN --mount=type=cache,target=/root/.cache/go-build \
    ego-go build

# New build stage for compiling the enclave with restricted flags mode
FROM build-enclave as build-enclave-testmode-false
# Sign the enclave executable
RUN ego sign enclave.json

# New build stage for compiling the enclave without restricted flags mode
FROM build-enclave as build-enclave-testmode-true
# Sign the enclave executable
RUN ego sign enclave-test.json

# Tag the restricted mode as the current build
FROM build-enclave-testmode-${TESTMODE} as build-enclave

# Trigger a new build stage and use the smaller ego version:
FROM ghcr.io/edgelesssys/ego-deploy:v1.3.0

# Copy just the binary for the enclave into this build stage
COPY --from=build-enclave \
    /home/ten/go-ten/go/enclave/main /home/ten/go-ten/go/enclave/main
    
WORKDIR /home/ten/go-ten/go/enclave/main

# simulation mode is ACTIVE by default
ENV OE_SIMULATION=1
EXPOSE 11000