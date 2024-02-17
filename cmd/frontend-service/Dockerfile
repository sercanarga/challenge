FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go mod tidy && go mod verify
RUN CGO_ENABLED=0 go build -o frontend-service cmd/frontend-service/main.go

FROM gcr.io/distroless/static-debian12 as runner

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/frontend-service .
COPY --from=builder --chown=nonroot:nonroot /app/.env .

#EXPOSE 8080
ENTRYPOINT ["./frontend-service"]