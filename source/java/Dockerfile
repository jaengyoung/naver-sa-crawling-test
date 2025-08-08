FROM amazoncorretto:21-alpine as build

WORKDIR /app

COPY build.gradle .
COPY gradle/ gradle/
COPY gradlew .
COPY src/ ./src/

RUN chmod +x gradlew && \
    ./gradlew clean shadowJar

FROM public.ecr.aws/lambda/java:21

COPY --from=build /app/build/libs/java-lambda.jar ${LAMBDA_TASK_ROOT}/lib/

CMD ["example.LambdaHandler::handleRequest"]