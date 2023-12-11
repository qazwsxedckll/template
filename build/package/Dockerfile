FROM golang:1.21.1 AS builder

WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct

# modifiy other files will not invalidate cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o /out/core ./core

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /out/core .

ENTRYPOINT [ "./core" ]
CMD [ "run" ]
