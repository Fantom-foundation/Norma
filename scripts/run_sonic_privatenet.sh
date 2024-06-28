#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Sonic is going to export its services on ${external_ip}"

### Initialize datadir

# if validator, 
genesis_flag=""
if [[ -n $VALIDATOR_PUBKEY ]]
then
	genesis_flag="--mode validator"
fi

mkdir /datadir
./sonictool --datadir /datadir genesis json --experimental ${genesis_flag} genesis.json

# If validator, initialize here
pubkey_flag=""
if [[ -n $VALIDATOR_PUBKEY ]] 
then
	echo "Sonic is now running as validator, pubkey=${VALIDATOR_PUBKEY}"
	pubkey_flag="--validator.pubkey ${VALIDATOR_PUBKEY} --validator.password password --mode validator"
else
	echo "Sonic is now running as an observer"
fi

# Start sonic as part of a fake net with RPC service.
./sonicd \
    --datadir=/datadir \
    --http --http.addr 0.0.0.0 --http.port 18545 --http.api admin,eth,ftm \
    --ws --ws.addr 0.0.0.0 --ws.port 18546 --ws.api admin,eth,ftm \
    --pprof --pprof.addr 0.0.0.0 \
    --nat=extip:${external_ip} \
    --metrics \
    --metrics.expensive
