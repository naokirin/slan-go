FROM golang

RUN go get -u github.com/golang/dep/...

RUN mkdir -p /go/src/github.com/naokirin/slan-go
RUN mkdir -p /go/src/github.com/naokirin/slan-go/db
WORKDIR /go/src/github.com/naokirin/slan-go

ADD Gopkg.lock /go/src/github.com/naokirin/slan-go/
ADD Gopkg.toml /go/src/github.com/naokirin/slan-go/
RUN dep ensure -vendor-only=true

ADD ./app /go/src/github.com/naokirin/slan-go/app
RUN go build -o slan-go app/main.go

ADD ./secrets /go/src/github.com/naokirin/slan-go/secrets
ADD ./config /go/src/github.com/naokirin/slan-go/config

CMD ["./slan-go"]

