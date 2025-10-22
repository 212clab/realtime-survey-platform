# 1단계: 빌드 스테이지
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Go 모듈 의존성 관리
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# 소스 코드 복사 및 빌드
COPY . .
RUN go build -o /main .

# 2단계: 최종 실행 스테이지
FROM alpine:latest

WORKDIR /app

# 빌드 스테이지에서 생성된 실행 파일만 복사
COPY --from=builder /main .

# 컨테이너 실행 시 이 명령어를 실행
CMD ["./main"]