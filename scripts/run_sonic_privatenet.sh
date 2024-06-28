#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Sonic is going to export its services on ${external_ip}"

datadir="/datadir"

##
## if $VALIDATOR_ID is set, it is a validator
##

# Initialize datadir

genesis_flag=""
if [[ -n $VALIDATOR_ID ]]
then
	genesis_flag="--mode validator"
fi

mkdir /datadir
./sonictool --datadir ${datadir} genesis json --experimental ${genesis_flag} genesis.json

# Create pubkey, secretfile

if [[ -n $VALIDATOR_ID ]]
then
	cmd=`normatool --datadir ${datadir} validator from -id ${VALIDATOR_ID}`
	res=($cmd)
	VALIDATOR_PUBKEY=res[0]
	VALIDATOR_SECRET=res[2]
fi

# If validator, initialize here
val_flag=""
if [[ -n $VALIDATOR_ID ]] 
then
	echo "Sonic is now running as validator, pubkey=${VALIDATOR_PUBKEY}"
	val_flag="--validator.id ${VALIDATOR_ID} --validator.pubkey ${VALIDATOR_PUBKEY} --validator.password ${VALIDATOR_SECRET} --mode validator"
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
