# STEP 1 - build executable binary
FROM golang:1.16-buster AS builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/api-server cmd/api-server/main.go

# STEP 2 - build small image
FROM scratch
EXPOSE 80/TCP
WORKDIR /
COPY --from=builder /app/bin/api-server /bin/api-server

ENV NEWS_APP_FEEDS_MGMT_WEBSERVER_HOST 0.0.0.0
ENV NEWS_APP_FEEDS_MGMT_WEBSERVER_PORT 80
CMD ["/bin/api-server"]