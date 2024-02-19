FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy && go mod verify
RUN CGO_ENABLED=0 go build -o worker-service cmd/worker-service/main.go

FROM gcr.io/distroless/static-debian12 as runner

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/worker-service .
COPY --from=builder --chown=nonroot:nonroot /app/.env .

ENTRYPOINT ["./worker-service"]