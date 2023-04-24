# Start from the latest golang base image
FROM golang:1.17 as builder

# Add Maintainer Info
LABEL maintainer="Selly <dev@selly.vn>"

# Set the Current Working Directory inside the container
RUN mkdir -p /go/src/git.selly.red/Selly-Server/warehouse
WORKDIR /go/src/git.selly.red/Selly-Server/warehouse

COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy data to working dir
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o ./cmd/app/main ./cmd/app/main.go

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /go/src/git.selly.red/Selly-Server/warehouse

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/git.selly.red/Selly-Server/warehouse .

# Command to run the executable
CMD ["./cmd/app/main"]