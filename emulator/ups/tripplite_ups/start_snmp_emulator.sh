#!/bin/bash

# Start the snmp emulator for the Tripplite UPS.

# Args are:
# data directory
# port
# log file name
# SNMP version, which is currently only V3 for this UPS.

DATA_DIRECTORY=$1
NETWORK_PORT=$2
LOG_DIRECTORY=$3
SNMP_VERSION=$4

# Enure SNMP_VERSION is set
if [ -z ${SNMP_VERSION+x} ]; then
  echo "SNMP_VERSION is unset";
else echo "SNMP_VERSION is set to '${SNMP_VERSION}'";
fi

# SNMP V3 only for this UPS.
if [[ ${SNMP_VERSION} -ne V3 ]] ; then
  echo "SNMP_VERSION is not V3"
  exit 1
fi

# The snmp emulator cannot run as root.
# Running a python file to load a configuration file will work,
# but then we don't have access to things like snmpsimd.py when we Popen.
# It is not a path issue.
#

python `which snmpsimd.py` \
    --data-dir=${DATA_DIRECTORY} \
    --agent-udpv4-endpoint=0.0.0.0:${NETWORK_PORT} \
    --v3-user=simulator \
    --v3-auth-key=auctoritas \
    --v3-auth-proto=MD5 \
    --v3-priv-key=privatus \
    --v3-priv-proto=DES \
    2>&1 | tee /logs/${LOG_DIRECTORY}
