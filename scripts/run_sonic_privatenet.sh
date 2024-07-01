#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Sonic is going to export its services on ${external_ip}"

datadir="/datadir"
echo "val id=${VALIDATOR_ID}"

##
## if $VALIDATOR_ID is set, it is a validator
##

# Initialize datadir

mkdir /datadir
./sonictool --datadir ${datadir} genesis json --experimental genesis.json

# Create pubkey, secretfile

if [[ $VALIDATOR_ID -ne 0 ]]
then
	cmd=`./normatool validator from -id ${VALIDATOR_ID} -d ${datadir}`
	res=($cmd)
	VALIDATOR_PUBKEY=${res[0]}
	VALIDATOR_SECRET=${res[2]}
fi

# If validator, initialize here
val_flag=""
if [[ $VALIDATOR_ID -ne 0 ]]
then
	echo "Sonic is now running as validator"
	echo "val.id=${VALIDATOR_ID}"
	echo "pubkey=${VALIDATOR_PUBKEY}"
	echo "secret=${VALIDATOR_SECRET}"
	val_flag="--validator.id ${VALIDATOR_ID} --validator.pubkey ${VALIDATOR_PUBKEY} --validator.password ${VALIDATOR_SECRET} --mode rpc"
else
	echo "Sonic is now running as an observer"
fi

# Start sonic as part of a fake net with RPC service.
./sonicd \
    --datadir=${datadir} \
    ${val_flag} \
    --http --http.addr 0.0.0.0 --http.port 18545 --http.api admin,eth,ftm \
    --ws --ws.addr 0.0.0.0 --ws.port 18546 --ws.api admin,eth,ftm \
    --pprof --pprof.addr 0.0.0.0 \
    --nat=extip:${external_ip} \
    --metrics \
    --metrics.expensive
