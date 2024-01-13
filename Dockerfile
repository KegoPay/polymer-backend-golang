FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main -ldflags="-s -w"
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
ADD /infrastructure/messaging/emails/templates /app/infrastructure/messaging/emails/templates
EXPOSE 8080
CMD ["./main"]
