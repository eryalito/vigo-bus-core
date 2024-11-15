FROM docker.io/library/golang:1.23 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/main .

FROM ubuntu:24.04

RUN apt-get update && apt install -y python3 curl && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/main /app/main

CMD ["/app/main"]
