# -------- Stage 1: Builder --------
FROM golang:1.24 AS builder

WORKDIR /app

# go.mod, go.sum 복사
COPY go.mod go.sum ./

# 의존성 설치
RUN go mod download

# 전체 소스 파일 복사
COPY .. .

# Lambda용 빌드 (Linux, amd64 환경 및 CGO 비활성화)
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

# -------- Stage 2: Lambda Runner --------
FROM public.ecr.aws/lambda/go:1

# 빌드된 실행 파일 복사
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}

# Lambda 진입점
CMD ["main"]