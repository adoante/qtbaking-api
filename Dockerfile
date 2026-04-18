FROM golang:1.26.1-alpine AS build

# Set working directory
WORKDIR /backend

# Copy go module files first for optimal layer caching
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy application source code
COPY backend/ .

# Build static binary
RUN go build -o /backend/api .

FROM alpine:latest

WORKDIR /backend
COPY --from=build /backend/api .

EXPOSE 8080

CMD ["/backend/api"]
