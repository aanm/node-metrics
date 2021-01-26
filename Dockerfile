FROM docker.io/library/golang:1.15.7 as builder
LABEL maintainer="maintainer@cilium.io"
ADD . /go/src/github.com/aanm/node-metrics
WORKDIR /go/src/github.com/aanm/node-metrics
RUN make node-metrics
RUN strip node-metrics

FROM docker.io/library/alpine:3.9.3 as certs
ARG CILIUM_SHA=""
LABEL cilium-sha=${CILIUM_SHA}
RUN apk --update add ca-certificates

FROM docker.io/library/busybox:1.31.1
ARG CILIUM_SHA=""
LABEL cilium-sha=${CILIUM_SHA}
LABEL maintainer="maintainer@cilium.io"
COPY --from=builder /go/src/github.com/aanm/node-metrics/node-metrics /usr/bin/node-metrics
COPY --from=builder /go/src/github.com/aanm/node-metrics/get_metric.sh /usr/bin/get_metric.sh
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /
CMD ["/usr/bin/node-metrics"]
