FROM golang:1.16

WORKDIR /go/app/

ENV GO111MODULE=on

ADD . /go/app/.

RUN go build

CMD ["go", "run", "main.go"]

