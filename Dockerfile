FROM alpine:3.1

RUN apk update
RUN apk add ca-certificates

COPY khaos-monkey /

CMD ["/khaos-monkey"]