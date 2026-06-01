ARG GO_BASE_IMAGE=swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/golang:1.22.12-alpine3.21
FROM ${GO_BASE_IMAGE} AS builder

WORKDIR /src

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/billserver ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=builder /out/billserver /app/bin/billserver
COPY --from=builder /src/configs /app/configs

EXPOSE 8080

CMD ["/app/bin/billserver"]
