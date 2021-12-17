FROM golang:1.17.5-alpine3.15

WORKDIR /go/src/app
COPY ../.. .
#RUN apk add --no-cache bash

RUN go get -d -v ./...
RUN go install -v ./...
#RUN go build cmd/api
#RUN go build cmd/migrate
#RUN go build



CMD ["api"]