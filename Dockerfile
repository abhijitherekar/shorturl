FROM herekar/abhi-go:1.0

RUN apk update

RUN apk -v --update \
    add git build-base && \
    rm -rf /var/cache/apk/* && \
    mkdir -p "$GOPATH/src/github.com/webpage"

ADD . "$GOPATH/src/github.com/webpage"

RUN cd "$GOPATH/src/github.com/webpage" && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o shorturl .

COPY ./shorturl /bin/shorturl

ENTRYPOINT ["/bin/shorturl"]

