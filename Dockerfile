# Begin with a builder image from alpine
FROM golang:alpine as builder

LABEL maintainer="Felipe Galdino dos Santos <felipegaldino16@gmail.com>"


# Install git
RUN apk update && apk add --no-cache git


# Copy go.mod and go.sum in order to download the dependencies
WORKDIR /app
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the project to the working directory
COPY . .

# Build Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



# ============================================================================================================

# Now start the creation of the image to run the application
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binaries from the builder container
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080
EXPOSE 8080

# Run the app
CMD ["./main"]