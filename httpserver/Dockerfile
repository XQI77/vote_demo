FROM golang:1.23-alpine AS builder

WORKDIR /src

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GONOSUMCHECK=*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/httpserver ./httpserver

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /out/httpserver .
COPY httpserver/static ./static

EXPOSE 8888

ENTRYPOINT ["./httpserver", "-f", "/app/etc/httpserver-api.yaml"]
