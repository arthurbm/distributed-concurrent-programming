FROM golang:latest

WORKDIR /app

COPY . .
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

RUN go mod init server
RUN go get -d
RUN go build -o /go/bin/server

EXPOSE 5672

CMD [ "/wait-for-it.sh", "rabbitmq:5672", "--", "/go/bin/server" ]
