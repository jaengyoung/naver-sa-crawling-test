FROM amazoncorretto:21-alpine as build

WORKDIR /app

COPY pom.xml .
COPY src/ ./src/
COPY LambdaHandler.java ./src/main/java/example/

RUN apk add --no-cache maven && \
    mvn clean package -DskipTests

FROM public.ecr.aws/lambda/java:21

COPY --from=build /app/target/*.jar ${LAMBDA_TASK_ROOT}/lib/

CMD ["example.LambdaHandler::handleRequest"]