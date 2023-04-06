FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod init github.com/nuriansyah/services-logbook-mbkm && go get -d -v ./... && go build -o /app ./cmd/server/main.go

EXPOSE 8080

CMD ./app
