# 1단계: 빌드 환경 (Go 컴파일러)
FROM golang:1.25-alpine AS builder
WORKDIR /app

# WebDAV 공식 확장 패키지 다운로드 및 초기화
RUN go mod init nas-webdav && \
    go get golang.org/x/net/webdav

COPY main.go .
# 어디서나 실행 가능하도록 정적 바이너리로 컴파일 (경량화 최적화)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o webdav-server main.go

# 2단계: 가벼운 실행 환경 (Alpine Linux)
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /

# 1단계에서 빌드된 결과물만 쏙 빼오기
COPY --from=builder /app/webdav-server /webdav-server

EXPOSE 80

CMD ["/webdav-server"]