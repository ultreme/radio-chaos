# build
FROM            golang:1.14-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/ultre.me/radio-chaos/chaos-bot
COPY            go.* ./
RUN             go mod download
RUN             cp -rf /go/pkg/mod/moul.io/pipotron@*/dict /dict
COPY            . ./
RUN             go install

# minimalist runtime
FROM alpine:3.11
COPY            --from=builder /go/bin/chaos-bot /bin/
WORKDIR         /dict
COPY            --from=builder /dict/ .
COPY            replies.yml .
ENTRYPOINT      ["/bin/chaos-bot"]
