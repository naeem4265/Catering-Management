FROM golang:latest
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
