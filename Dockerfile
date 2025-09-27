FROM golang:1.23-alpine AS build

# Install git (needed for some Go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o api ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -s /bin/sh apiuser

WORKDIR /app
COPY --from=build /app/api ./
COPY --from=build /app/.env* ./

USER apiuser

EXPOSE 8080
CMD [ "./api" ]
