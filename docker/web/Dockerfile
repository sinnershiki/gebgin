FROM golang:1.16

WORKDIR /go/app/

ENV GO111MODULE=on

ADD . /go/app/.

RUN go mod download

CMD ["go", "run", "cmd/gebgin/main.go"]

