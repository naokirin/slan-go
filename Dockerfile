FROM golang

RUN go get -u github.com/golang/dep/...
RUN mkdir -p /go/src/github.com/naokirin/slan-go
ADD . /go/src/github.com/naokirin/slan-go
WORKDIR /go/src/github.com/naokirin/slan-go
RUN dep ensure -vendor-only=true

CMD ["go", "run", "main.go"]

