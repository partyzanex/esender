FROM golang:1.12.6-alpine as builder

RUN apk add --no-cache openssh-client git

COPY . /go/src/github.com/partyzanex/esender
WORKDIR /go/src/github.com/partyzanex/esender

ENV CGO_ENABLED 0
ENV GO111MODULE on

RUN go build -v -o /esender ./http/server/main.go


FROM alpine:3.9

EXPOSE 9000

WORKDIR /service

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /esender /service/esender

RUN apk --update add tzdata ca-certificates && \
  cp /usr/share/zoneinfo/UTC /etc/localtime && \
  echo "UTC" > /etc/timezone && \
  date && \
  apk del tzdata


ENTRYPOINT ["./esender"]
