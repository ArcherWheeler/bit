FROM golang:alpine
WORKDIR /go/src/github.com/ArcherWheeler/bit
RUN apk update && apk --no-cache add git
RUN go get -u github.com/Masterminds/glide

# Install glide dependencies before copying source changes
# so small code changes don't invalide existing vendor tree
WORKDIR /go/src/github.com/ArcherWheeler/bit/bit
COPY bit/glide.yaml glide.yaml
COPY bit/glide.lock glide.lock
RUN glide install

WORKDIR /go/src/github.com/ArcherWheeler/bit/test
COPY test/glide.yaml glide.yaml
COPY test/glide.lock glide.lock
RUN glide install

WORKDIR /go/src/github.com/ArcherWheeler/bit/bit
COPY bit/ .
RUN go install

WORKDIR /go/src/github.com/ArcherWheeler/bit/test
COPY test/ .
ENV TEST_FAILSAFE="off"
CMD ["go", "test", "cli_test.go"]
