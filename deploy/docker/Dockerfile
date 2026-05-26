FROM eclipse-temurin:17-jdk-alpine AS builder

WORKDIR /workspace

COPY pom.xml .
COPY src src

RUN apk add --no-cache maven protobuf \
    && mvn clean package -DskipTests -Pprod \
    && java -Djarmode=layertools -jar target/user-system-*.jar extract

FROM eclipse-temurin:17-jre-alpine

RUN apk add --no-cache curl tzdata bind-tools \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && addgroup -g 1000 appgroup \
    && adduser -u 1000 -G appgroup -D appuser

WORKDIR /app

COPY --from=builder /workspace/dependencies/ ./
COPY --from=builder /workspace/spring-boot-loader/ ./
COPY --from=builder /workspace/snapshot-dependencies/ ./
COPY --from=builder /workspace/application/ ./

RUN mkdir -p /app/logs /app/data \
    && chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080 9999

HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/actuator/health/liveness || exit 1

ENV JAVA_OPTS="-Xms512m -Xmx2048m \
    -XX:+UseG1GC \
    -XX:MaxGCPauseMillis=200 \
    -XX:+HeapDumpOnOutOfMemoryError \
    -XX:HeapDumpPath=/app/logs/heapdump.hprof \
    -XX:+ExitOnOutOfMemoryError \
    -Djava.security.egd=file:/dev/./urandom \
    -Dfile.encoding=UTF-8 \
    -Duser.timezone=Asia/Shanghai"

ENTRYPOINT ["sh", "-c", "java ${JAVA_OPTS} org.springframework.boot.loader.launch.JarLauncher"]