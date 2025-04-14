FROM golang:1.23 AS builder

RUN addgroup gouser && adduser --ingroup gouser --disabled-password gouser && \
    mkdir -p /app && chown gouser:gouser /app
WORKDIR /app

COPY --chown=gouser:gouser go.mod go.sum ./
RUN go mod download

USER gouser

COPY --chown=gouser:gouser . .

RUN go build -o /app/main /app/cmd/main.go

FROM golang:1.23 AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY --from=builder /app ./

CMD ["air", "-c", ".air.toml"]

FROM golang:1.23 AS prod

WORKDIR /app

COPY --from=builder /app/main /app/main

CMD ["./main"]