FROM golang:1.19-alpine

# Set the working directory
WORKDIR /app

# Copy the necessary files
RUN go mod download && go mod verify

COPY . .

# Add Git
RUN apk add --no-cache git

# Build the application
RUN go build -o main cmd/server/main.go


# Expose the port
EXPOSE 8080

# Set the entrypoint command
CMD ["./main"]
