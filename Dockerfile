#FROM golang:1.16
#ENV GOPROXY="https://goproxy.cn" \
#    CGO_ENABLED=1 \
#    GOOS=linux \
#    GOARCH=amd64 \
#    GO111MODULE=on \
#    CONFIG="prod"
#
#WORKDIR /
#COPY ./src /src
#COPY ./config /config
#RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone
#WORKDIR /src
#RUN go build .
#EXPOSE 3456
#ENTRYPOINT ["./bug-carrot"]

FROM golang:1.16-alpine3.15 as builder
ENV GOPROXY="https://goproxy.cn" \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

RUN apk add build-base
WORKDIR /
COPY ./src /src
WORKDIR /src
RUN go build -o /build/app .

FROM alpine:3.15
ENV CONFIG="prod"

RUN apk --no-cache add -U tzdata ca-certificates libc6-compat libgcc libstdc++ && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > "/etc/timezone"
COPY --from=builder /build/app /usr/bin/app
COPY --from=builder /go/pkg/mod/github.com/ttys3 /go/pkg/mod/github.com/ttys3
COPY ./config /config
WORKDIR /
EXPOSE 3456
ENTRYPOINT [ "app" ]
