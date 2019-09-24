# Builder Image
FROM vaporio/golang:1.13 as builder
WORKDIR /go/src/github.com/vapor-ware/synse-snmp-plugin
COPY . .

RUN make build CGO_ENABLED=0


# Plugin Image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/vapor-ware/synse-snmp-plugin/build/plugin ./plugin

# NOTE: The plugin config (config.yml) is not added in here, nor are
#  any device configs. Those configurations should be supplied on a
#  per-deployment basis.

# Image Metadata -- http://label-schema.org/rc1/
# This should be set after the dependency install so we can cache that layer
ARG BUILD_DATE
ARG BUILD_VERSION
ARG VCS_REF

LABEL maintainer="vapor@vapor.io" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="vaporio/sflow-plugin" \
      org.label-schema.vcs-url="https://github.com/vapor-ware/synse-sflow-plugin" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vendor="Vapor IO" \
      org.label-schema.version=$BUILD_VERSION

ENTRYPOINT ["./plugin"]
