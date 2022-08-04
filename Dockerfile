ARG BASE_IMAGE_FULL
ARG BASE_IMAGE_MINIMAL

# Build node feature discovery
FROM golang:1.17.2-buster as builder

ENV GOPROXY="https://goproxy.io"
# Download the grpc_health_probe bin
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Get (cache) deps in a separate layer
COPY go.mod go.sum /go/node-feature-discovery/

WORKDIR /go/node-feature-discovery

RUN go mod download

# Do actual build
COPY . /go/node-feature-discovery

ARG VERSION
ARG HOSTMOUNT_PREFIX

RUN make install VERSION=$VERSION HOSTMOUNT_PREFIX=$HOSTMOUNT_PREFIX

RUN make test

# Create full variant of the production image
FROM ${BASE_IMAGE_FULL} as full

# Install pciutils and update pci.ids to the current version
RUN apt-get update && apt-get install pciutils wget curl -y && update-pciids

# Run as unprivileged user
USER 65534:65534

# Use more verbose logging of gRPC
ENV GRPC_GO_LOG_SEVERITY_LEVEL="INFO"

COPY --from=builder /go/node-feature-discovery/deployment/components/worker-config/nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf
COPY --from=builder /go/bin/* /usr/bin/
COPY --from=builder /bin/grpc_health_probe /usr/bin/grpc_health_probe

# Create minimal variant of the production image
FROM ${BASE_IMAGE_MINIMAL} as minimal

# Run as unprivileged user
USER 65534:65534

# Use more verbose logging of gRPC
ENV GRPC_GO_LOG_SEVERITY_LEVEL="INFO"

COPY --from=builder /go/node-feature-discovery/deployment/components/worker-config/nfd-worker.conf.example /etc/kubernetes/node-feature-discovery/nfd-worker.conf
COPY --from=builder /go/bin/* /usr/bin/
COPY --from=builder /bin/grpc_health_probe /usr/bin/grpc_health_probe
