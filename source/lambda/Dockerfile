# -------- Stage 1: Builder (경량 이미지 사용) --------
# Alpine 기반으로 빌드 속도 향상
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 종속성 파일만 먼저 복사 (Docker 레이어 캐싱 최적화)
COPY go.mod go.sum ./

# 종속성 다운로드
RUN go mod download && go mod verify

# 소스 코드 복사 (필요한 파일만)
COPY main.go ./
COPY internal/ ./internal/

# Lambda 컨테이너용 정적 빌드
#   CGO_ENABLED=0: C 라이브러리 의존성 제거
#   GOOS=linux: Linux 환경용  
#   -ldflags="-s -w": 디버그 정보 제거로 바이너리 크기 축소
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -a -installsuffix cgo \
    -o main main.go

# -------- Stage 2: Lambda Go Runtime --------
# Go Lambda 전용 이미지로 최적화
FROM public.ecr.aws/lambda/go:1

# 빌드된 실행파일 복사
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main

# Lambda 핸들러 (main 함수 실행)
CMD ["main"]






# # -------- Stage 1: Builder (경량 이미지 사용) --------
# FROM golang:1.24-alpine AS builder

# WORKDIR /app

# # 종속성 파일만 먼저 복사 (Docker 레이어 캐싱 최적화)
# COPY go.mod go.sum ./

# # 종속성 다운로드 (Go 모듈 캐시 활용)
# # --mount=type=cache로 빌드 속도 향상
# RUN --mount=type=cache,target=/go/pkg/mod \
#     go mod download && go mod verify

# # 소스 코드 복사 (필요한 파일만)
# COPY main.go ./
# COPY internal/ ./internal/

# # Lambda 컨테이너용 정적 빌드
# # CGO_ENABLED=0: C 라이브러리 의존성 제거
# # GOOS=linux: Linux 환경용
# # -ldflags="-s -w": 디버그 정보 제거로 바이너리 크기 축소
# # -trimpath: 빌드 경로 정보 제거 (보안)
# RUN --mount=type=cache,target=/go/pkg/mod \
#     --mount=type=cache,target=/root/.cache/go-build \
#     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
#     -ldflags="-s -w -extldflags '-static'" \
#     -a -installsuffix cgo \
#     -trimpath \
#     -o main main.go

# # -------- Stage 2: Lambda Go Runtime --------
# FROM public.ecr.aws/lambda/go:1

# # 빌드된 실행파일 복사
# COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main

# # Lambda 핸들러 (main 함수 실행)
# CMD ["main"]