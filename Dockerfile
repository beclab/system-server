# Build the manager binary
FROM golang:1.22.6 as builder

WORKDIR /workspace

RUN git clone https://github.com/kubernetes/code-generator.git bytetrade.io/web3os/code-generator && \
    cd bytetrade.io/web3os/code-generator && git checkout -b release-1.27 && cd -

# Copy the Go Modules manifests
COPY go.mod bytetrade.io/web3os/system-server/go.mod
COPY go.sum bytetrade.io/web3os/system-server/go.sum

# Build
RUN cd bytetrade.io/web3os/system-server && \
        go mod download

# Copy the go source
COPY . bytetrade.io/web3os/system-server/

RUN cd bytetrade.io/web3os/system-server && \
    CGO_ENABLED=1 go build -ldflags="-s -w" -a -o system-server cmd/server/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/base:nonroot
FROM gcr.io/distroless/base:debug
WORKDIR /
COPY --from=builder /workspace/bytetrade.io/web3os/system-server/system-server .

VOLUME [ "/data" ]

ENTRYPOINT ["/system-server"]
