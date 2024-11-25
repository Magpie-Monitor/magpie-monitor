ARG GO_VERSION=1.23.3
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN go env -w GOCACHE=/go/pkg/mod/

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    go build -o /bin/server ./services/reports/cmd/reports

FROM build AS tests

RUN echo "go env -w GOCACHE=/go/pkg/mod/ && go test ./services/reports/..." > test.sh 
RUN chmod +x test.sh

ENTRYPOINT ["bash", "./test.sh"]

FROM debian:latest AS final

RUN apt-get update && apt-get install -y ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /bin/server /bin/

# Copy the openapi.yaml specification
COPY ./services/reports/api /srv
COPY ./pkg/swagger/index.html /srv/swagger-ui.html

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
