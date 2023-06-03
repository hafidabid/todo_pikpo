# Start from the official Go image
FROM golang:1.17 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN sudo apt install -y protobuf-compiler
RUN protoc --version
RUN chmod +x generate_proto.sh
RUN ./generate_proto.sh

# Copy the entire project to the container
COPY . .

# Build the Go application
RUN go build -o /go-app

# Start from a clean image
FROM golang:1.17

# Copy the built executable from the previous stage
COPY --from=builder /go-app /app

# Set the working directory inside the container
WORKDIR /app

# Expose the port that the application listens on
EXPOSE 8080

# Run the Go application
CMD ["./go-app"]
