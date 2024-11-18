#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Opera is going to export its services on ${external_ip}"

# Starting opera as part of a fake net with RPC service.
./opera --fakenet ${VALIDATOR_NUMBER}/${VALIDATORS_COUNT} \
    --statedb.impl=${STATE_DB_IMPL} \
    --vm.impl=${VM_IMPL} \
    --fakenetgaspower 1000 \
    --http --http.addr 0.0.0.0 --http.port 18545 --http.api admin,eth,ftm \
    --ws --ws.addr 0.0.0.0 --ws.port 18546 --ws.api admin,eth,ftm \
    --pprof --pprof.addr 0.0.0.0 \
    --nat=extip:${external_ip} \
    --metrics \
    --metrics.expensive

