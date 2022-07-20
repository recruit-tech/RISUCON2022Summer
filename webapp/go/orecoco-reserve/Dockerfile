FROM golang:1.17.6-alpine

RUN apk update && apk add git bash mysql-client

WORKDIR /webapp/go/orecoco-reserve

COPY go/orecoco-reserve/go.mod go/orecoco-reserve/go.sum ./

RUN go mod download

COPY go/orecoco-reserve/main.go ./
RUN go build -o main

EXPOSE 3003

COPY sql ../../sql

COPY scripts/wait-for-db.sh ./

COPY go/orecoco-reserve/meeting_room.txt ./

CMD ["./wait-for-db.sh", "./main"]
