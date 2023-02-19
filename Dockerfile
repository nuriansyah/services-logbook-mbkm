FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -o services-logbook-mbkm

EXPOSE 8080

CMD ./api-gateway