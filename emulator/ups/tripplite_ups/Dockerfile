# Ths dockerfile starts the SNMP emulator for the tripplite ups in a container for testing.
# This emulator runs on port 1025.
FROM docker.io/vaporio/vapor-endpoint-base-x64:1.0

RUN pip install -I \
    snmpsim \
    pysnmp \
    pyasn1

# The emulator will not start as root, so we need to add a user.
# As root, the Error message is the following:
# snmp-emulator    | ERROR: cant drop priveleges: Must drop priveleges to a non-priveleged user&group (sic)
# Create the user and ${HOME}
RUN groupadd -r docker && useradd -r -g docker snmp
ADD . /home/snmp

# These SNMP emulator files are specfic to the device being emulated.
# Data are just places in /home/snmp/data on the emulator to keep it simple.
ADD data/public.snmpwalk /home/snmp/data/public.snmpwalk

# snmpsmi variation modules (like writecache) are getting installed to a location not in the search path,
# so copy where it will be found.
# snmp user owns /home/snmp and /logs.
RUN mkdir -p /home/snmp/.snmpsim/variation && \
    cp /usr/local/snmpsim/variation/* /home/snmp/.snmpsim/variation && \
    chown snmp:docker /home/snmp -R && \
    chown snmp:docker /logs -R
USER snmp
WORKDIR /home/snmp

# Default emulator port is 1024. Expose it.
EXPOSE 1024/udp

# data directory (typically /home/snmp/data)
# port (typically 1024 and up)
# log file name (trying to keep these names unique)
# SNMP version (V3)
CMD ["./start_snmp_emulator.sh","./data","1025","snmp-emulator-tripplite-ups.log", "V3"]
