FROM golang:1.10.3
WORKDIR /go/src/github.com/ilya-korotya/solid
COPY . /go/src/github.com/ilya-korotya/solid
RUN go get ./...
CMD ["go", "run", "main.go"]
