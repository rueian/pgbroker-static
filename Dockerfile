FROM golang:1.16 as build

WORKDIR /go/src/proxy
ADD . /go/src/proxy

RUN go build -o /go/bin/proxy cmd/proxy/main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/proxy /
CMD ["/proxy"]