FROM golang:1.13-stretch as go-build

COPY . /go/src/github.com/rickbassham/example-go
WORKDIR /go/src/github.com/rickbassham/example-go/

ARG appname
ARG githash
ARG gittag

RUN echo "export APP_NAME=${appname}" >> start.sh
RUN echo "export BUILD_GIT_HASH=${githash}" >> start.sh
RUN echo "export BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" >> start.sh
RUN echo "export BUILD_GIT_TAG=${gittag}" >> start.sh
RUN echo "./app" >> start.sh

RUN chmod +x start.sh

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -v -o app ./${appname}

# Build the final image.
FROM alpine:3.10

USER root
RUN apk update && apk add tzdata ca-certificates && rm -rf /var/cache/apk/* && update-ca-certificates

RUN adduser -S app

USER app
WORKDIR /app
COPY --from=go-build /go/src/github.com/rickbassham/example-go/app .
COPY --from=go-build /go/src/github.com/rickbassham/example-go/start.sh .

CMD /bin/sh ./start.sh
