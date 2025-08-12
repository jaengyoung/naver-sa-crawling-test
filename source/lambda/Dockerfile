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

# -------- Stage 1: Builder (경량 이미지 사용) --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY main.go ./

COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o main main.go

# -------- Stage 2: Lambda Go Runtime --------
FROM public.ecr.aws/lambda/go:1

COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}/main

CMD ["main"]