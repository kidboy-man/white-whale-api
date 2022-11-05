FROM golang:1.18-alpine3.16

RUN apk update
RUN apk add git

ADD . /go/src/fetch-app
WORKDIR /go/src/fetch-app
COPY go.mod . 
COPY go.sum .
COPY .env .

RUN go mod tidy -v
RUN go build -v


EXPOSE 8080

CMD [ "./fetch-app"]
