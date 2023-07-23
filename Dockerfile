FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags='go_tarantool_ssl_disable' -o ./dialogs ./cmd/app/main.go

FROM scratch
COPY --from=builder /app/dialogs /usr/bin/dialogs
COPY --from=builder /app/internal/migrations /usr/bin/migrations
ENTRYPOINT [ "/usr/bin/dialogs" ]