FROM golang:1.15

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app

# RUN go get -v
# RUN go build . -o service
# ENTRYPOINT ["/go/src/app/service"]