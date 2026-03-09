FROM docker.io/maven:3-eclipse-temurin-21 AS build
WORKDIR /workspace/app
ARG OTEL_JAVA_AGENT_VERSION=2.8.0

COPY pom.xml .
COPY src src

RUN --mount=type=cache,target=/root/.m2 mvn install -DskipTests
RUN mkdir -p target/dependency && (cd target/dependency; jar -xf ../*.jar)
RUN curl -fSL --retry 5 --retry-delay 2 \
    -o /workspace/app/opentelemetry-javaagent.jar \
    https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v${OTEL_JAVA_AGENT_VERSION}/opentelemetry-javaagent.jar

FROM docker.io/eclipse-temurin:21-jre-alpine
RUN apk add --no-cache curl ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app \
    && mkdir -p /opt/todobackend/app
WORKDIR /opt/todobackend
VOLUME /tmp
ARG DEPENDENCY=/workspace/app/target/dependency
COPY --from=build --chown=app:app ${DEPENDENCY}/BOOT-INF/lib /opt/todobackend/app/lib
COPY --from=build --chown=app:app ${DEPENDENCY}/META-INF /opt/todobackend/app/META-INF
COPY --from=build --chown=app:app ${DEPENDENCY}/BOOT-INF/classes /opt/todobackend/app
COPY --from=build --chown=app:app /workspace/app/opentelemetry-javaagent.jar /opt/todobackend/opentelemetry-javaagent.jar

USER app

ENTRYPOINT ["java", "-cp", "/opt/todobackend/app:/opt/todobackend/app/lib/*",  "-javaagent:/opt/todobackend/opentelemetry-javaagent.jar","io.novatec.todobackend.TodobackendApplication"]
