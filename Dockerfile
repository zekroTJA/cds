FROM golang:1.21-alpine AS build
WORKDIR /build
COPY cmd cmd
COPY pkg pkg
COPY go.mod .
COPY go.sum .
RUN go build -o dist/cds cmd/cds/main.go


FROM alpine
WORKDIR /app
COPY --from=build /build/dist/cds /app/cds
EXPOSE 80
ENV CDS_ADDRESS="0.0.0.0:80"
ENTRYPOINT [ "/app/cds" ]