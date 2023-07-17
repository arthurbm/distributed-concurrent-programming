FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod init server
RUN go get -d
RUN go build -o /go/bin/server

EXPOSE 1313

CMD [ "/go/bin/server" ]
