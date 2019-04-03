FROM golang:1.12.0-alpine3.9 as builder

RUN apk add --no-cache openssh-client git

COPY . /go/src/github.com/partyzanex/esender
WORKDIR /go/src/github.com/partyzanex/esender

ENV CGO_ENABLED 0
ENV GO111MODULE on

RUN go get -u github.com/pressly/goose/cmd/goose

RUN go build -v -o /esender ./http/server/main.go



FROM alpine:3.9

EXPOSE 9000

WORKDIR /service

ARG tz
ARG MYSQL_ROOT_PASSWORD
ARG MYSQL_DATABASE

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /esender /service/esender
COPY --from=builder /go/bin/goose /service/goose
COPY --from=builder /go/src/github.com/partyzanex/esender/boiler/migrations/mysql /service/mysql

RUN apk --update add tzdata ca-certificates && \
  cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
  echo "$tz" > /etc/timezone && \
  date && \
  apk del tzdata


ENTRYPOINT ["./esender"]
