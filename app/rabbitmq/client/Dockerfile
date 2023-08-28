FROM golang:latest

WORKDIR /app

COPY . .
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

RUN go mod init client
RUN go get -d
RUN go build -o /go/bin/client

CMD [ "/wait-for-it.sh", "rabbitmq:5672", "--", "/go/bin/client" ]
