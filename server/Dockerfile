# Server Dockerfile

FROM golang:1.19.1-buster

WORKDIR /app

COPY . .
RUN go build -o /server

EXPOSE 8080

CMD ["/server"]

