# Use the official Golang base image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app with static compilation flags to ensure it runs on Alpine which uses musl libc
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pinger ./src

# Start a new stage from scratch
FROM alpine:3.20

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/pinger .

# Expose port 8080 to the outside world
EXPOSE 8080

# Ensure the binary is executable
RUN chmod +x ./pinger

# Command to run the executable
CMD ["./pinger"]