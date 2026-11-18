# Build context: the repository root (see compose.yaml), so libs/shared can be copied in.
FROM golang:1.27-alpine AS builder
WORKDIR /src
COPY libs/shared ./libs/shared
COPY services/api ./services/api
WORKDIR /src/services/api
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o /out/api .

FROM alpine:3.21
COPY --from=builder /out/api /usr/local/bin/api
EXPOSE 8080
USER nobody
ENTRYPOINT ["/usr/local/bin/api"]
