#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Opera is going to export its services on ${external_ip}"

# Starting opera as part of a fake net with RPC service.
./opera --fakenet ${VALIDATOR_NUMBER}/${VALIDATORS_COUNT} --statedb.impl=${STATE_DB_IMPL} --http --http.addr 0.0.0.0 --http.api admin,eth,ftm --nat=extip:${external_ip} --pprof --pprof.addr 0.0.0.0
