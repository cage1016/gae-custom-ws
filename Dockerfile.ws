FROM gcr.io/gcp-runtimes/go1-builder:1.12.7-2019-07-10-211636 AS builder

RUN apt-get update && apt-get -y install git
WORKDIR /go/src/app

COPY . .

RUN /usr/local/go/bin/go build -o app cmd/ws/main.go

# Application image.
FROM gcr.io/distroless/base:latest
COPY --from=builder /go/src/app/app /usr/local/bin/app

CMD ["/usr/local/bin/app"]