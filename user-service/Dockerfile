# 1단계: 빌드 스테이지
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# 빌드 경로를 ./cmd로 변경!
RUN go build -o /main ./cmd

# 2단계: 최종 실행 스테이지
FROM alpine:latest

WORKDIR /app
COPY --from=builder /main .
CMD ["./main"]