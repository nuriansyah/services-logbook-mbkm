FROM golang:1.19

WORKDIR /app

COPY . .

COPY go.mod .

RUN go mod download && go mod verify

RUN go build -o /main

EXPOSE 8080

CMD [ "/main" ]