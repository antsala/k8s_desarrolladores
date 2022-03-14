FROM  golang:1.11-alpine
ADD . /go/src
RUN go install helloContainer

FROM alpine:latest
COPY --from=0 /go/bin/helloContainer .
CMD ["./helloContainer"]