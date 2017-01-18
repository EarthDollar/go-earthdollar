FROM alpine:3.3

ADD . /go-ethereum
RUN \
  apk add --update git go make gcc musl-dev         && \
  (cd go-ethereum && make geth)                     && \
  cp go-ethereum/build/bin/geth /geth               && \
  apk del git go make gcc musl-dev                  && \
  rm -rf /go-ethereum && rm -rf /var/cache/apk/*

EXPOSE 8811
EXPOSE 20203

ENTRYPOINT ["/geth"]
