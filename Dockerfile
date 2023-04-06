FROM golang:1.19

WORKDIR /app

COPY . .

RUN go build -o app ./cmd/server


EXPOSE 8080

CMD ./app
