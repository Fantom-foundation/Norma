#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Sonic is going to export its services on ${external_ip}"

datadir="/datadir"
echo "val id=${VALIDATOR_ID}"
echo "genesis validator count=${VALIDATORS_COUNT}"

# Call set genesis script - add balance to all possible validators
# VALIDATOR_COUNT defines genesis validator count
# TODO change 100 funded validator addresses to specific number
./set_genesis.sh genesis.json 100 ${VALIDATORS_COUNT} ${MAX_BLOCK_GAS} ${MAX_EPOCH_GAS}

# Initialize datadir
mkdir /datadir /genesis
./sonictool --datadir ${datadir} genesis json --experimental genesis.json
cp genesis.json /genesis

##
## if $VALIDATOR_ID is set, it is a validator
##
if [[ $VALIDATOR_ID -ne 0 ]]
then
	cmd=`./normatool validator from -id ${VALIDATOR_ID} -d ${datadir}`
	res=($cmd)
	VALIDATOR_PUBKEY=${res[0]}
	VALIDATOR_ADDRESS=${res[1]}
fi

# Create password file - "password" is default normatool accounts password
echo password >> password.txt
VALIDATOR_PASSWORD="password.txt"

# If validator, initialize here
val_flag=""
if [[ $VALIDATOR_ID -ne 0 ]]
then
	echo "Sonic is now running as validator"
	echo "val.id=${VALIDATOR_ID}"
	echo "pubkey=${VALIDATOR_PUBKEY}"
	echo "address=${VALIDATOR_ADDRESS}"
	val_flag="--validator.id ${VALIDATOR_ID} --validator.pubkey ${VALIDATOR_PUBKEY} --validator.password ${VALIDATOR_PASSWORD} --mode rpc"
else
	echo "Sonic is now running as an observer"
fi

# Create config.toml
# when network starts with only one genesis validator, then he will not wait to start emitting
# if there are two or more validators at genesis they have to wait 5 seconds after connecting to the network
# if another validator connects to the network during run it will wait also 5 seconds to start emitting
echo [Emitter.EmitIntervals] >> config.toml
if [[ $VALIDATORS_COUNT == 1 && $VALIDATOR_ID == 1 ]]
then
  echo DoublesignProtection = 0 >> config.toml
else
#  5 seconds in golang time 5*10^9 nanoseconds
  echo DoublesignProtection = 5000000000 >> config.toml
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
    --metrics.expensive \
    --config config.toml
