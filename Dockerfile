# Build the application from source
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /config-reloader

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /config-reloader /config-reloader

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/config-reloader"]
