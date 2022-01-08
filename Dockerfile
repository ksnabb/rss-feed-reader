FROM golang:1.17

WORKDIR /go/src/rss-feed-reader
COPY . .

RUN go install -v ./...

EXPOSE 8080

CMD ["rss-feed-reader"]
