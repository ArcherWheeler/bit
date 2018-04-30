FROM golang:alpine
WORKDIR /go/src/github.com/ArcherWheeler/bit
RUN apk update && apk --no-cache add git
RUN go get -u github.com/Masterminds/glide
COPY bit/glide.yaml bit/glide.yaml
COPY bit/glide.lock bit/glide.lock
COPY test/glide.yaml test/glide.yaml
COPY test/glide.lock test/glide.lock
COPY ./ ./

WORKDIR /go/src/github.com/ArcherWheeler/bit/bit
RUN glide install
RUN go install

WORKDIR /go/src/github.com/ArcherWheeler/bit/test
ENV TEST_FAILSAFE="off"
RUN go test cli_test.go
