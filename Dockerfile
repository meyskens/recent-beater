FROM --platform=$BUILDPLATFORM node:16 as frontend

COPY ./frontend /opt/recent-beater
WORKDIR /opt/recent-beater

RUN npm install
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.16-alpine as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN apk add --no-cache git

COPY ./ /go/src/github.com/meyskens/recent-beater

WORKDIR /go/src/github.com/meyskens/recent-beater

RUN export GOARM=6 && \
    export GOARCH=amd64 && \
    if [ "$TARGETPLATFORM" == "linux/arm64" ]; then export GOARCH=arm64; fi && \
    if [ "$TARGETPLATFORM" == "linux/arm" ]; then export GOARCH=arm; fi && \
    go build -ldflags "-X main.revision=$(git rev-parse --short HEAD)" ./cmd/recent-beater/

FROM alpine:3.13

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/meyskens/recent-beater/recent-beater /usr/local/bin/

COPY --from=frontend /opt/recent-beater/build/ /opt/recent-beater/static/
WORKDIR /opt/recent-beater

CMD [ "/usr/local/bin/recent-beater", "serve" ]
