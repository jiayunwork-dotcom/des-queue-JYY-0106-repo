FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /des-queue .

FROM alpine:3.19
COPY --from=builder /des-queue /usr/local/bin/des-queue
COPY web/ /app/web/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["des-queue"]
CMD ["serve", "-addr", ":8080"]
