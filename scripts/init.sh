KEY="my_key_test"
CHAINID="lava_123-1"
MONIKER="monikertest"
KEYRING="os"
LOGLEVEL="info"

# Reinstall daemon
rm -rf ~/.simapp*
#make install

simd keys add $KEY --keyring-backend $KEYRING

MY_VALIDATOR_ADDRESS=$(simd keys show $KEY -a --keyring-backend $KEYRING)

simd init $MONIKER --chain-id $CHAINID

simd tendermint show-validator

simd add-genesis-account $MY_VALIDATOR_ADDRESS 1000000000000000stake --keyring-backend $KEYRING

simd gentx $KEY 100000000stake --keyring-backend $KEYRING --chain-id $CHAINID --gas=0stake

simd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
simd validate-genesis

simd keys list

