FROM golang:1.22.12-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
# Set GO111MODULE to on to ensure module mode is used
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/main /app/time-service
ENV GO_ENV=development
EXPOSE 8084
CMD ["/app/time-service"]