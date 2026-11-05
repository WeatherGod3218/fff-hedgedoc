FROM node:22-alpine AS npmbuilder

WORKDIR /app

COPY web/package*.json ./

RUN npm ci

COPY web/ .

RUN npm run build

FROM golang:1.26.2-alpine AS gobuilder

WORKDIR /app

COPY go.* .
RUN go mod download

COPY *.go .
COPY internal internal

COPY --from=npmbuilder /app/dist/ ./web/dist/

RUN go build -o server .

FROM alpine:latest

WORKDIR /app
COPY --from=gobuilder /app/server .
COPY web web

EXPOSE 8080

CMD ["./server"]
