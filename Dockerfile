FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod init && go build -o app ./cmd/server

EXPOSE 8080

CMD ./app
