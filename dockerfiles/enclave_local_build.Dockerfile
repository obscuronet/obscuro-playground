FROM ghcr.io/edgelesssys/ego-dev:latest

# on the container:
#   /home/obscuro/data       contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
RUN mkdir /home/obscuro
RUN mkdir /home/obscuro/data
RUN mkdir /home/obscuro/go-obscuro

# build the enclave from the current branch
COPY . /home/obscuro/go-obscuro
RUN cd /home/obscuro/go-obscuro/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "/home/obscuro/go-obscuro/go/obscuronode/enclave/main/main"]
EXPOSE 11000
