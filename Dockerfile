FROM golang:1.17 as builder

WORKDIR /build

# Build flag
ENV CGO_ENABLED=0 

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o foodbot main.go

# Install CA certs
FROM alpine as certimage
RUN apk add --no-cache ca-certificates

# Minimal Prod Image
FROM scratch

# Copy binary & CA certs from builder.
COPY --from=builder /build/foodbot /foodbot
COPY --from=certimage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD  ["./foodbot"]