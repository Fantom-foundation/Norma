#!/bin/bash

GENESIS_PATH=$1
ADDRESS_LOAD_ACCOUNTS_COUNT=$2
VALIDATORS_COUNT=$3

# Set genesis validator accounts
accounts=""
for (( i=1; i<=${ADDRESS_LOAD_ACCOUNTS_COUNT}; i++ )); do
  cmd=`./normatool validator from -id ${i}`
  res=($cmd)
  validator_address=${res[7]}

  # Remove prefix and reassign
  validator_address=${validator_address}
  echo "validator_address=${validator_address}"

  # Build the mask for this iteration
  mask=", { \"name\": \"validator${i}\", \"address\": \"${validator_address}\", \"balance\": 1000000000000000000000000000 }"
  accounts+="$mask"
done
echo "accounts=${accounts}"
sed -i 's|GENESIS_VALIDATOR_ACCOUNTS_PLACEHOLDER|'"$accounts"'|g' "$GENESIS_PATH"

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
sed -i 's|GENESIS_VALIDATOR_DEPLOYMENTS_PLACEHOLDER|'"$validators"'|g' "$GENESIS_PATH"
