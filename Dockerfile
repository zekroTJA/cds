FROM golang:1.12.6-stretch

LABEL maintainer="zekro <conatct@zekro.de>"

ENV PATH="$GOPATH/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

RUN apt-get update -y &&\
    apt-get install -y \
        git

WORKDIR $GOPATH/src/github.com/zekroTJA/cds

ADD . .

RUN dep ensure -v

RUN go build -v -o ./bin/cds -ldflags "\
		-X github.com/zekroTJA/cds/internal/static.AppVersion=$(git describe --tags) \
        -X github.com/zekroTJA/cds/internal/static.Release=TRUE" \
        ./cmd/cds/*.go

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/data &&\
    mkdir -p /etc/pages

EXPOSE 8080

CMD ./bin/cds \
        -c /etc/config/config.yml \
        -addr ":8080"