FROM golang:1.24rc2-alpine3.21

WORKDIR /app

COPY . .

RUN go get -d -v ./...

EXPOSE 8080

CMD ["go", "run", "./cmd/main.go"]