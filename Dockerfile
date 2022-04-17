FROM golang:1.16
ENV GOPROXY="https://goproxy.cn" \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on \
    CONFIG="prod"

WORKDIR /
COPY ./src /src
COPY ./config /config
WORKDIR /src
RUN go build .
EXPOSE 3456
ENTRYPOINT ["./bug-carrot"]
