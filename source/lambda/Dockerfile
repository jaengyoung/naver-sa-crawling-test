# # -------- Stage 1: Builder (경량 이미지 사용) --------
# FROM golang:1.24-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download && go mod verify

# COPY main.go ./

# COPY internal/ ./internal/

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
#     -ldflags="-s -w -extldflags '-static'" \
#     -a -installsuffix cgo \
#     -trimpath \
#     -o main main.go

# # -------- Stage 2: Lambda Go Runtime --------
# FROM public.ecr.aws/lambda/go:1

# COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main

# CMD ["main"]



# # -------- Stage 1: Builder (경량 이미지 사용) --------
# FROM golang:1.24-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download && go mod verify

# COPY main.go ./

# COPY internal/ ./internal/

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
#     -ldflags="-s -w" \
#     -o main main.go

# # -------- Stage 2: Lambda Go Runtime --------
# FROM public.ecr.aws/lambda/go:1

# COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main

# CMD ["main"]


# -------- Stage 1: Builder (경량 이미지) --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 종속성 파일만 먼저 복사 (Docker 레이어 캐싱 최적화)
COPY go.mod go.sum ./

# Go 모듈 병렬 다운로드
RUN GOPROXY=https://proxy.golang.org,direct go mod download && go mod verify

# 소스 코드 복사 (필요한 파일만)
COPY main.go ./
COPY internal/ ./internal/

# Lambda용 최적화 빌드 (성능 + 크기 최적화)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -ldflags="-s -w -extldflags '-static'" \
    -tags netgo \
    -installsuffix netgo \
    -p 4 \
    -o main main.go

# -------- Stage 2: Lambda Go Runtime --------
FROM public.ecr.aws/lambda/go:1

# 빌드된 실행파일 복사
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main

# Lambda 핸들러 (main 함수 실행)
CMD ["main"]