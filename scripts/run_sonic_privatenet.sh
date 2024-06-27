#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Sonic is going to export its services on ${external_ip}"

# Initialize datadir
mkdir /datadir
./sonictool --datadir=/datadir genesis json --experimental genesis.json

# Start sonic as part of a fake net with RPC service.
./sonicd \
    --datadir=/datadir \
    --http --http.addr 0.0.0.0 --http.port 18545 --http.api admin,eth,ftm \
    --ws --ws.addr 0.0.0.0 --ws.port 18546 --ws.api admin,eth,ftm \
    --pprof --pprof.addr 0.0.0.0 \
    --nat=extip:${external_ip} \
    --metrics \
    --metrics.expensive
