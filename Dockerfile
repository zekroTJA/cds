FROM golang:1.13-buster AS build

WORKDIR /build

ADD . .

RUN go build -v -o ./bin/cds -ldflags "\
		-X github.com/zekroTJA/cds/internal/static.AppVersion=$(git describe --tags) \
        -X github.com/zekroTJA/cds/internal/static.Release=TRUE" \
        ./cmd/cds/*.go


FROM debian:buster-slim AS final

LABEL maintainer="zekro <conatct@zekro.de>"

WORKDIR /app

COPY --from=build /build/bin .

RUN mkdir -p /etc/config &&\
    mkdir -p /etc/data &&\
    mkdir -p /etc/pages

EXPOSE 8080

CMD /app/cds \
        -c /etc/config/config.yml \
        -addr ":8080"