FROM golang:1.22.2 AS build

WORKDIR /app
COPY . /app
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -s /bin/sh apiuser

WORKDIR /app
COPY --from=build /app/api ./
COPY --from=build /app/.env* ./

USER apiuser

EXPOSE 8080
CMD [ "./api" ]
