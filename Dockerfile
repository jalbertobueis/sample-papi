FROM golang:alpine

ARG pkg=sample-papi

RUN apk add --no-cache ca-certificates

COPY . $GOPATH/src/$pkg

RUN set -ex \
      && apk add --no-cache --virtual .build-deps \
              git \
      && go get -v $pkg/... \
      && apk del .build-deps

RUN go install $pkg/...

# Needed for templates for the front-end app.
#WORKDIR $GOPATH/src/$pkg/app

# Users of the image should invoke either of the commands.
CMD echo "OK!"; exit 1