FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o /app ./cmd/server

EXPOSE 8080

CMD ["/app/server"]
