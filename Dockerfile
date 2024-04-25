FROM golang:1.20 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY a2w.go .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/a2w a2w.go

FROM scratch

ENV GIN_MODE=release
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/a2w .
COPY templates templates

EXPOSE 5001

ENTRYPOINT ["./a2w"]
