# syntax = docker/dockerfile:experimental
FROM golang:1.24 as builder

ARG version=dev

WORKDIR /workspace

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY api api
COPY controllers controllers
COPY job-task-runner job-task-runner
COPY tools tools
COPY version version

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-X code.cloudfoundry.org/korifi/version.Version=${version}" -o manager job-task-runner/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
