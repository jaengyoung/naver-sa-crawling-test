FROM golang:1.24-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM public.ecr.aws/lambda/provided:al2

COPY --from=build /app/main ${LAMBDA_RUNTIME_DIR}

CMD ["main"]