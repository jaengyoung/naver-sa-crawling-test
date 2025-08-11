FROM openjdk:21-jdk-slim AS build

WORKDIR /app

COPY build.gradle.kts .
COPY src/ ./src/

RUN apt-get update && apt-get install -y wget unzip && \
    wget -q https://services.gradle.org/distributions/gradle-8.5-bin.zip && \
    unzip gradle-8.5-bin.zip && \
    mv gradle-8.5 /opt/gradle && \
    ln -s /opt/gradle/bin/gradle /usr/bin/gradle && \
    gradle clean shadowJar

FROM public.ecr.aws/lambda/java:21

COPY --from=build /app/build/libs/kotlin-lambda.jar ${LAMBDA_TASK_ROOT}/lib/

CMD ["example.LambdaHandler::handleRequest"]