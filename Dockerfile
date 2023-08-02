FROM golang:1.20 as builder

WORKDIR /workspace

# Copy what is needed for builds
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY Makefile Makefile
COPY .git/ .git/

# Run the build
RUN make build

# Final image
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.2

# Update the image to get the latest CVE updates
RUN microdnf update -y && \
    microdnf install -y openssl && \
    microdnf install -y util-linux && \
    microdnf clean all

WORKDIR /
COPY --from=builder /workspace/bin/lvm_driver .
EXPOSE 23532
USER 65532:65532

ENTRYPOINT [ "/lvm_driver" ]