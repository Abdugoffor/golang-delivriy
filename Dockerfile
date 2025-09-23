# Dockerfile

# Base image
FROM golang:1.24.4-alpine

# Workdir
WORKDIR /app

# Copy go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build app
RUN go build -o main .

# Run
CMD ["./main"]
