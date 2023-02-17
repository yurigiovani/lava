KEY="client_1"
CHAINID="lava_123-1"
#MONIKER="monikertest"
#KEYRING="os"
#LOGLEVEL="info"
#
## Reinstall daemon
#rm -rf ~/.simapp*
##make install
#
#output=$(echo "yes" | simd keys add $KEY --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_2" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_3" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_4" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_5" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_6" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_7" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_8" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_9" --keyring-backend $KEYRING)
#output=$(echo "yes" | simd keys add "client_10" --keyring-backend $KEYRING)
#
#MY_VALIDATOR_ADDRESS=$(simd keys show $KEY -a --keyring-backend $KEYRING)
#
#simd init $MONIKER --chain-id $CHAINID
#
#simd tendermint show-validator
#
#simd add-genesis-account $MY_VALIDATOR_ADDRESS 1000000000000000stake --keyring-backend $KEYRING
#
#simd gentx $KEY 100000000stake --keyring-backend $KEYRING --chain-id $CHAINID --gas=0stake
#
#simd collect-gentxs
#
## Run this to ensure everything worked and that the genesis file is setup correctly
#simd validate-genesis

simd tx bank send `simd keys show $KEY -a` `simd keys show client_2 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_3 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_4 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_5 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_6 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_7 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_8 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_9 -a` 500stake --chain-id $CHAINID
simd tx bank send `simd keys show $KEY -a` `simd keys show client_10 -a` 500stake --chain-id $CHAINID
