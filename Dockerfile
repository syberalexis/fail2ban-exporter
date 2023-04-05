FROM golang:1.19 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux make build


FROM alpine:3
ARG VERSION
COPY --from=builder /build/dist/fail2ban-exporter-${VERSION}-linux-amd64 fail2ban-exporter
RUN addgroup -S exporter && adduser -S exporter -G exporter
USER exporter
ENTRYPOINT [ "./linky-exporter" ]
