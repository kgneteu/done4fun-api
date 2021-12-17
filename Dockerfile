FROM golang:1.17.5-alpine3.15

WORKDIR /go/src/app
COPY ../.. .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["api"]