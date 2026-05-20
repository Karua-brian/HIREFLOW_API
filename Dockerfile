# Use official Golang image as the base image
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the entire project into the working directory
COPY . .

# Build the go application
RUN go build -o hireflow_API ./cmd/api/main.go

# Use a smaller base image for the final stage
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=0 /app/hireflow_API /app/hireflow_API

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./hireflow_API"]