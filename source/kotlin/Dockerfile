FROM openjdk:21-jdk-slim as build

WORKDIR /app

COPY build.gradle.kts .
COPY gradle/ gradle/
COPY gradlew .
COPY LambdaHandler.kt ./src/main/kotlin/example/

RUN chmod +x gradlew && \
    ./gradlew clean shadowJar

FROM public.ecr.aws/lambda/java:21

COPY --from=build /app/build/libs/kotlin-lambda.jar ${LAMBDA_TASK_ROOT}/lib/

CMD ["example.LambdaHandler::handleRequest"]