# 同一份 Dockerfile 建兩個服務：--build-arg SERVICE=api 或 SERVICE=guide
FROM golang:1.26 AS build
ARG SERVICE=api
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/${SERVICE}

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/server /server
ENV GIN_MODE=release
USER nonroot
ENTRYPOINT ["/server"]
