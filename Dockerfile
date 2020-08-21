# Start from the latest golang base image
FROM golang:latest

LABEL maintainer="Anufriev Artem <anufriev.artem.mail@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 9000

CMD ["./main"]
