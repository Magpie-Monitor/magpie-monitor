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
    go build -o /bin/server ./services/cluster_metadata/cmd/cluster_metadata

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
COPY ./services/cluster_metadata/api /srv
COPY ./pkg/swagger/index.html /srv/swagger-ui.html

EXPOSE 9090

ENTRYPOINT [ "/bin/server" ]
