FROM docker.io/library/golang:1.23 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o /app/main .

FROM scratch

COPY --from=build /app/main /app/main

CMD ["/app/main"]