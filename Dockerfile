FROM golang:latest as builder
WORKDIR /app
COPY . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o ./server ./cmd/main.go

FROM scratch
COPY --from=builder /app/server .
COPY --from=builder /app/configs/config.json ./configs/config.json
CMD ["./server"]