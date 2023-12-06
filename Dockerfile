# Use the Go 1.21 base image for building the application.
FROM golang:1.21-alpine as builder

# Set the working directory inside the container for building.
WORKDIR /app

# Copy the Go application source code to the working directory.
COPY ./ /app

# Download dependencies using Go Modules.
RUN go mod download

# Build the Go application binary.
RUN go build -o /app/cloney

# Stage 2: Create the final lightweight image.
# Use Alpine Linux 3.17 as the base image.
FROM alpine:3.17 as cloney

# Update the package repository on the Alpine system.
RUN apk update

# Copy the compiled application binary from the builder stage.
# Copy to /usr/local/bin so that the application is available in the PATH.
COPY --from=builder /app/cloney /usr/local/bin

# Sleep forever to keep the container running.
CMD ["sleep", "infinity"]
