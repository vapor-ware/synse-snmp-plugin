#
# Builder Image
#
FROM docker.io/vaporio/foundation:bionic as builder

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

#
# Final Image
#
FROM docker.io/vaporio/scratch-ish:1.0.0

LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.name="vaporio/snmp-plugin" \
      org.label-schema.vcs-url="https://github.com/vapor-ware/synse-snmp-plugin" \
      org.label-schema.vendor="Vapor IO"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# NOTE
#  Plugin and device configurations are not built into the
#  image. They should be supplied on a per-deployment basis.

# Copy the executable.
COPY synse-snmp-plugin ./plugin

ENTRYPOINT ["./plugin"]
