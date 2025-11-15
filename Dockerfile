FROM golang:1.25.3-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

COPY --from=builder /app/main .


# Set env variables
ENV PORT=8080
ENV JWT_SECRET=mysecret
ENV DB_HOST=db
ENV DB_USER=root
ENV DB_PASS=root
ENV DB_NAME=mydb

EXPOSE 8080

# Jalankan binary + migration + seeder saat container start
CMD ["sh", "-c", "./main migrate && ./main seed && ./main"]
