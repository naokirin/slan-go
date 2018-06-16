FROM golang

RUN go get -u github.com/golang/dep/...

RUN mkdir -p /go/src/github.com/naokirin/slan-go
RUN mkdir -p /go/src/github.com/naokirin/slan-go/db
RUN mkdir -p /go/src/github.com/naokirin/slan-go/lgtm
WORKDIR /go/src/github.com/naokirin/slan-go

ADD Gopkg.lock /go/src/github.com/naokirin/slan-go/
ADD Gopkg.toml /go/src/github.com/naokirin/slan-go/
RUN dep ensure -vendor-only=true

ADD ./app /go/src/github.com/naokirin/slan-go/app
RUN go build -o slan-go app/main.go

ADD ./img /go/src/github.com/naokirin/slan-go/img

CMD ["./slan-go"]

