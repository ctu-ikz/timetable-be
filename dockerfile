# Start with a base image containing the Go runtime
FROM golang:1.23.3-alpine3.20

# Create a directory for the application, and set it as the working directory
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o timetable-be server.go

# Use a smaller base image for the final container
FROM alpine:latest
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=0 /app/timetable-be .

# Copy the .env file if it’s needed in the runtime environment
COPY .env .

# Expose the port the application runs on
EXPOSE 8080

# Run the application by default when the container starts
CMD ["./timetable-be"]