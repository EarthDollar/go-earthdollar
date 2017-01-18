FROM alpine:3.4

RUN \
  apk add --update go git make gcc musl-dev && \
  git clone --depth 1 https://github.com/ethereum/go-ethereum && \
  (cd go-ethereum && make ged) && \
  cp go-ethereum/build/bin/ged /ged && \
  apk del go git make gcc musl-dev && \
  rm -rf /go-ethereum && rm -rf /var/cache/apk/*

EXPOSE 8811
EXPOSE 20203

ENTRYPOINT ["/ged"]
