# Get the local nodes IP 
list=`hostname -I`
array=($list)
echo "Opera is going to export its services on ${array[0]}"

# Starting opera as part of a fake net with RPC service
./opera --fakenet ${VALIDATOR_NUMBER}/${VALIDATORS_COUNT} --http --http.addr 0.0.0.0 --http.api admin,eth,ftm --nat=extip:${array[0]}