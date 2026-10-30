FROM golang:1.26.2-alpine AS gobuilder

WORKDIR /app

COPY go.* .
RUN go mod download

COPY *.go .

RUN go build -o server .

FROM alpine:latest

WORKDIR /app
COPY --from=gobuilder /app/server .

EXPOSE 8080

CMD ["./server"]
