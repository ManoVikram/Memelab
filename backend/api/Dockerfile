# =================================
# STAGE 1 - BUILDER
# =================================

# Step 1 - Start with the official Go base iamge
FROM golang:1.24-alpine AS builder

# Step 2 - Set working directory inside the container
WORKDIR /app

# Step 3 - Copy go.mod and go.sum first to handle dependencies
COPY go.mod go.sum ./

# Step 4 - Download dependencies
RUN go mod download

# Step 5 - Copy the rest of the code
COPY . .

# Step 6 - Compile the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# =================================
# STAGE 2 - RUNNER
# =================================

# Step 1 - Start with the official light weight Alpine linux
FROM alpine:latest

# Step 2 - Set working directory inside the container
WORKDIR /app

# Step 3 - Copy only the compiled binary from Stage 1
COPY --from=builder /app/server .

# Step 4 - Expose the Go/Gin server port
EXPOSE 8080

# Step 5 - Run the Go/Gin server
CMD ["./server"]