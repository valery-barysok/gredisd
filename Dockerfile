FROM golang:1.7.4

MAINTAINER Valery Barysok <valery.barysok@gmail.com>

COPY . /go/src/github.com/valery-barysok/gredisd
WORKDIR /go/src/github.com/valery-barysok/gredisd

RUN go get -t ./...
RUN CGO_ENABLED=0 go install -v -a

EXPOSE 16379
ENTRYPOINT ["gredisd"]
