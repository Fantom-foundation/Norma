#!/bin/bash
# Get the local node's IP.
list=`hostname -I`
array=($list)
external_ip=${array[0]}
echo "Sonic is going to export its services on ${external_ip}"

datadir="/datadir"
echo "val id=${VALIDATOR_ID}"
echo "genesis validator count=${VALIDATORS_COUNT}"

# Set genesis validator balances
balances=""
for (( i=1; i<=${VALIDATORS_COUNT}; i++ )); do
  cmd=`./normatool validator from -id ${i}`
  res=($cmd)
  validator_address=${res[7]}

  # Remove prefix and reassign
  validator_address=${validator_address}
  echo "validator_address=${validator_address}"

  # Build the mask for this iteration
  mask=", { \"name\": \"validator${i}\", \"address\": \"${validator_address}\", \"balance\": 1000000000000000000000000000 }"
  balances+="$mask"
done
echo "balances=${balances}"
sed -i 's|GENESIS_VALIDATOR_BALANCES_PLACEHOLDER|'"$balances"'|g' "genesis.json"

# Set genesis validator delegations
validators=""
for (( i=1; i<=${VALIDATORS_COUNT}; i++ )); do
  # Convert ID to zero-padded hex string
  id_256int_hex=$(printf "%064x" $i)
  cmd=`./normatool validator from -id ${i}`
  res=($cmd)
  validator_public_key=${res[6]}
  validator_address=${res[7]}

  # Remove prefixes and reassign
  validator_public_key=${validator_public_key#0x}
  validator_address=${validator_address#0x}

  echo "id_256int_hex=${id_256int_hex}"
  echo "validator_address=${validator_address}"
  echo "validator_public_key=${validator_public_key}"

  # Build the mask for this iteration
  mask=",{ \"name\": \"SetGenesisValidator${i}\", \"to\": \"0xd100a01e00000000000000000000000000000000\", \"data\": \"0x4feb92f3000000000000000000000000${validator_address}${id_256int_hex}000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005fe149c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042${validator_public_key}000000000000000000000000000000000000000000000000000000000000\" }, { \"name\": \"SetGenesisDelegation1\", \"to\": \"0xd100a01e00000000000000000000000000000000\", \"data\": \"0x18f628d4000000000000000000000000${validator_address}${id_256int_hex}0000000000000000000000000000000000000000000422ca8b0a00a425000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\" }"
  validators+="$mask"
done
echo "validators=${validators}"
sed -i 's|GENESIS_VALIDATOR_DEPLOYMENTS_PLACEHOLDER|'"$validators"'|g' "genesis.json"

# Initialize datadir
mkdir /datadir
./sonictool --datadir ${datadir} genesis json --experimental genesis.json

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
