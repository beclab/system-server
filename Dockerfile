# Build the manager binary
FROM golang:1.18 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod bytetrade.io/web3os/system-server/go.mod
COPY go.sum bytetrade.io/web3os/system-server/go.sum

# Copy the go source
COPY cmd/ bytetrade.io/web3os/system-server/cmd/
COPY pkg/ bytetrade.io/web3os/system-server/pkg/
COPY hack/ bytetrade.io/web3os/system-server/hack/

RUN git clone https://github.com/kubernetes/code-generator.git bytetrade.io/web3os/code-generator && \ 
    cd bytetrade.io/web3os/code-generator && git checkout -b release-1.27 && cd - 

# Build
RUN cd bytetrade.io/web3os/system-server && \
        go mod tidy 

RUN cd bytetrade.io/web3os/system-server && \ 
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o system-server cmd/server/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/base:nonroot
FROM gcr.io/distroless/base:debug
WORKDIR /
COPY --from=builder /workspace/bytetrade.io/web3os/system-server/system-server .

VOLUME [ "/data" ]

ENTRYPOINT ["/system-server"]
