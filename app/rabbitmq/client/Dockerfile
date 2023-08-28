FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod init client
RUN go get -d
RUN go build -o /go/bin/client

CMD [ "/go/bin/client" ]
