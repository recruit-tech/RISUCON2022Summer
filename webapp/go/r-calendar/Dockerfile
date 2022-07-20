FROM golang:1.17.6-alpine

RUN apk update && apk add git bash file mysql-client

WORKDIR /webapp/go/r-calendar

COPY go/r-calendar/go.mod go/r-calendar/go.sum ./

RUN go mod download

COPY go/r-calendar/main.go ./
RUN go build -o main

EXPOSE 3000

COPY sql ../../sql

COPY scripts/wait-for-db.sh ./

CMD ["./wait-for-db.sh", "./main"]
