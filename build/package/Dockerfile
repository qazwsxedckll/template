FROM golang AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
ENV GOTOOLCHAIN=auto

# modifiy other files will not invalidate cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o /out/main .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /out/main .

ENTRYPOINT [ "./main" ]
CMD [ "run" ]
