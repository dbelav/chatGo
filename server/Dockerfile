FROM golang:latest

WORKDIR /app/server

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

EXPOSE 8080

CMD ["air"]