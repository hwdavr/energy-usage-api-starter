# 1) build
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./cmd/server

# 2) runtime
FROM gcr.io/distroless/base-debian12
ENV ADDR=:8080
EXPOSE 8080
COPY --from=build /bin/server /server
ENTRYPOINT ["/server"]
