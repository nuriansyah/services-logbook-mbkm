FROM golang:1.19

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

RUN go build -o /main

EXPOSE 8080

CMD [ "/main" ]