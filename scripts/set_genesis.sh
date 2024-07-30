#!/bin/bash

GENESIS_PATH=$1
ADDRESS_LOAD_ACCOUNTS_COUNT=$2
VALIDATORS_COUNT=$3
MAX_BLOCK_GAS_ARG=$4
MAX_EPOCH_GAS_ARG=$5

# Set genesis network parameters/rules
DEFAULT_MAX_BLOCK_GAS=20500000000
DEFAULT_MAX_EPOCH_GAS=1500000000000

if [[ -n "$MAX_BLOCK_GAS_ARG" && "$MAX_BLOCK_GAS_ARG" -gt 0 ]]; then
  MaxBlockGas=$MAX_BLOCK_GAS_ARG
else
  MaxBlockGas=$DEFAULT_MAX_BLOCK_GAS
fi

if [[ -n "$MAX_EPOCH_GAS_ARG" && "$MAX_EPOCH_GAS_ARG" -gt 0 ]]; then
  MaxEpochGas=$MAX_EPOCH_GAS_ARG
else
  MaxEpochGas=$DEFAULT_MAX_EPOCH_GAS
fi

echo "MaxBlockGas=${MaxBlockGas}"
echo "MaxEpochGas=${MaxEpochGas}"

rules="{ \"networkName\": \"norma-privatenet\", \"networkId\": \"0xfa3\", \"MaxBlockGas\": $MaxBlockGas, \"MaxEventGas\": 10028000000, \"MaxEpochGas\": $MaxEpochGas, \"ShortGasAllocPerSec\": 5600000000000, \"LongGasAllocPerSec\": 2800000000000 }"
sed -i 's|GENESIS_RULES_PLACEHOLDER|'"$rules"'|g' "$GENESIS_PATH"

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
