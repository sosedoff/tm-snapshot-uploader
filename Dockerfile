# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.17 AS build

ARG PACKAGE
ARG GIT_COMMIT

WORKDIR /go/src/${PACKAGE}

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.14 AS release

WORKDIR /app

ARG PACKAGE

COPY --from=build \
  /go/src/${PACKAGE}/build/tm-snapshot-uploader \
  /app/tm-snapshot-uploader

RUN addgroup --gid 1234 app && \
    adduser --system --uid 1234 app && \
    chown -R app:app /app

USER 1234

ENTRYPOINT ["/app/tm-snapshot-uploader"]
