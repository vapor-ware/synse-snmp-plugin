FROM iron/go:dev as builder
WORKDIR /go/src/github.com/vapor-ware/synse-snmp-plugin
COPY . .
RUN make build


FROM iron/go
LABEL maintainer="vapor@vapor.io"

WORKDIR /plugin

COPY --from=builder /go/src/github.com/vapor-ware/synse-snmp-plugin/build/plugin ./plugin
# COPY config.yml .
# TODO: We need to remove this copy since it's the emulator config.yml.
# Copying here will cause deployments to fail since /plugin/config.yml gets
# picked up before /etc/synse/plugin/config/config.yml which means that a
# deploy to real hardware would use the emulator config.yml.
# TODO: We will need to bring this back somewhow, as well as to fix the tests.

CMD ["./plugin"]
