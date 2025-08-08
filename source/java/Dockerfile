FROM amazoncorretto:21-alpine AS build

WORKDIR /app

COPY build.gradle .
COPY src/ ./src/

RUN apk add --no-cache gradle && \
    gradle clean shadowJar

FROM public.ecr.aws/lambda/java:21

COPY --from=build /app/build/libs/java-lambda.jar ${LAMBDA_TASK_ROOT}/lib/

CMD ["example.LambdaHandler::handleRequest"]