FROM golang:1.22-alpine3.20 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o ./previewer ./cmd/previewer/main.go

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/previewer .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["/app/previewer"]
