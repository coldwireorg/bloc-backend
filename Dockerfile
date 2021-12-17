# Building server
FROM golang:1.17-alpine AS builder
WORKDIR /build
RUN apk --no-cache add ca-certificates
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build main.go

# running app
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /build/main /app/main

ENV SERVER_DOMAIN=coldwire.org
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=3000
ENV SERVER_HTTPS=false
ENV DB_ADDRESS=127.0.0.1
ENV DB_PORT=27017
ENV DB_NAME=bloc
ENV STORAGE_DIR=/opt/bloc/files

CMD ["/app/main"]
