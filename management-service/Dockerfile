FROM eclipse-temurin:21-jdk-jammy as builder
WORKDIR /opt/app
COPY mvnw pom.xml ./
COPY .mvn/ .mvn
RUN chmod +x ./mvnw && ./mvnw dependency:go-offline
COPY ./src ./src
RUN ./mvnw clean install -DskipTests

FROM eclipse-temurin:21-jre-jammy
WORKDIR /opt/app
EXPOSE 8080
COPY --from=builder /opt/app/target/*.jar /opt/app/*.jar
ENTRYPOINT ["java", "-jar", "/opt/app/*.jar" ]
