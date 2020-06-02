FROM golang:1.14

WORKDIR $GOPATH/src/github.com/aryan9600/shorten-url

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8000

CMD ["shorten-url"]
