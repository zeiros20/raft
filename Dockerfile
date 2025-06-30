FROM golang:1.24-alpine

# Install git so Go can fetch modules (even if not needed now)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy mod files first to use Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the binary
RUN go build -o raft-node main.go

# Run the binary
ENTRYPOINT ["./raft-node"]

