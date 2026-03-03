FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go mod tidy
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8082
CMD ["./main"]
