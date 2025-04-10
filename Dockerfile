FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./auth ./cmd/main.go

EXPOSE 8080

RUN chmod +x auth

CMD ["./auth"]

