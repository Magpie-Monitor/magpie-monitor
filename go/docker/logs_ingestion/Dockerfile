ARG GO_VERSION=1.22.5
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN go env -w GOCACHE=/go/pkg/mod/

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    go build -o /bin/server ./services/logs_ingestion/cmd/logs_ingestion

FROM build AS tests

# ARG LOGSDB_PORT
# ARG LOGSDB_HOST
# ARG LOGSDB_USER
# ARG LOGSDB_PASSWORD
# ARG LOGSDB_CERT_PATH
# ARG LOGS_INGESTION_QUEUE_HOST
# ARG LOGS_INGESTION_QUEUE_PORT
# ARG LOGS_INGESTION_QUEUE_USERNAME
# ARG LOGS_INGESTION_QUEUE_PASSWORD
#
# RUN --mount=type=cache,target=/go/pkg/mod/ \
#     go test ./...

# RUN --mount=type=cache,target=/go/pkg/mod/ \
#     --mount=type=bind,source=go.sum,target=go.sum \
#     --mount=type=bind,source=go.mod,target=go.mod \
#     go mod download -x

# RUN go env -w GOCACHE=/go/pkg/mod/

RUN echo "go env -w GOCACHE=/go/pkg/mod/ && go test ./..." > test.sh
RUN chmod +x test.sh

# ENTRYPOINT [ "go", "env" , "-w", "GOCACHE=/go/pkg/mod/", "&&" ,"go", "test", "./..." ]

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

EXPOSE 8090

ENTRYPOINT [ "/bin/server" ]



