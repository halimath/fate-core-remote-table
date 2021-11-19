FROM registry.access.redhat.com/ubi8/ubi-init as JDK

ARG openjdk_download_location=https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17%2B35/OpenJDK17-jdk_x64_linux_hotspot_17_35.tar.gz

RUN dnf install -y wget
RUN wget -O /tmp/openjdk.tar.gz ${openjdk_download_location}
RUN mkdir -p /usr/local/java
RUN tar xvfz /tmp/openjdk.tar.gz -C /usr/local/java

ENV JAVA_HOME /usr/local/java/jdk-17+35
ENV PATH "$PATH/:$JAVA_HOME/bin"

RUN java -version

FROM JDK as MVN

RUN mkdir -p /src
WORKDIR /src
COPY mvnw .
COPY .mvn ./.mvn
COPY service ./service
COPY app ./app
COPY pom.xml .
RUN ./mvnw package

FROM JDK

ARG commit

RUN mkdir -p /app
WORKDIR /app
COPY --from=MVN /src/service/target/quarkus-app .

ENV JAVA_OPTS ""
ENV APP_COMMIT=${commit}

ENTRYPOINT java ${JAVA_OPTS} -jar ./quarkus-run.jar