FROM golang:1.10.3 as builder
WORKDIR /go/src/github.com/ilya-korotya/solid
COPY . /go/src/github.com/ilya-korotya/solid
# TODO: use go dep
RUN go get ./...
RUN CGO_ENABLED=0 go build -o solid .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/ilya-korotya/solid .
# TODO: Do we need to build a container each time to change the code in the docker?
# Or you can skip everything through 'CMD'
CMD ["./solid"]  
