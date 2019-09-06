#
# Builder Image
#
FROM vaporio/golang:1.13 as builder

#
# Final Image
#
FROM scratch

LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.name="vaporio/snmp-plugin" \
      org.label-schema.vcs-url="https://github.com/vapor-ware/synse-snmp-plugin" \
      org.label-schema.vendor="Vapor IO"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# NOTE
#  Plugin and device configurations are not built into the
#  image. They should be supplied on a per-deployment basis.

# Copy the executable.
COPY synse-emulator-plugin ./plugin

ENTRYPOINT ["./plugin"]
