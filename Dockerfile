FROM golang:1.23-alpine

# Install Git because Go mod needs it
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first (to leverage layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Now copy the full source
COPY . .

# Build the binary from main/main.go
RUN go build -o raft-node ./main

# Run the binary
CMD ["./raft-node"]
