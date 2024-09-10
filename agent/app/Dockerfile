FROM golang:1.21.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server app/cmd/main.go

EXPOSE 8080

CMD ["./server"]