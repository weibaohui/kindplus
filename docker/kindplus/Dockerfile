FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN go build -o kindplus .

FROM nginx
COPY --from=builder /build/kindplus /usr/local/bin/kindplus
CMD ["/usr/local/bin/kindplus"]
