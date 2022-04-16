FROM alpine:3.15
RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
COPY ./app /bin/app
COPY ./config /config
WORKDIR /
EXPOSE 3456
ENTRYPOINT [ "/bin/app" ]