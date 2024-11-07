FROM golang:alpine as builder
WORKDIR /build
COPY .  .
RUN ls
RUN go build -o kindplus .

FROM nginx:alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache curl bash inotify-tools supervisor
ADD docker/kindplus/kindplus.conf /etc/supervisord.conf

ADD docker/kindplus/reload.sh /app/reload.sh
COPY --from=builder /build/kindplus /app/kindplus
RUN chmod +x /app/reload.sh
ENTRYPOINT ["/app/reload.sh"]