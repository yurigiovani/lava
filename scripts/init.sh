KEY="client_1"
CHAINID="lava_123-1"
MONIKER="monikertest"
KEYRING="os"
LOGLEVEL="info"

# Reinstall daemon
rm -rf ~/.simapp*
#make install

output=$(echo "yes" | simd keys add $KEY --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_2" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_3" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_4" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_5" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_6" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_7" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_8" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_9" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_10" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_11" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_12" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_13" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_14" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_15" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_16" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_17" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_18" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_19" --keyring-backend $KEYRING)
output=$(echo "yes" | simd keys add "client_20" --keyring-backend $KEYRING)

MY_VALIDATOR_ADDRESS=$(simd keys show $KEY -a --keyring-backend $KEYRING)

simd init $MONIKER --chain-id $CHAINID

simd tendermint show-validator

simd add-genesis-account $MY_VALIDATOR_ADDRESS 500000000stake --keyring-backend $KEYRING

simd add-genesis-account $(simd keys show "client_2" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_3" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_4" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_5" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_6" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_7" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_8" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_9" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_10" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_11" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_12" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_13" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_14" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_15" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_16" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_17" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_18" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_19" -a) 500000000stake  --keyring-backend $KEYRING
simd add-genesis-account $(simd keys show "client_20" -a) 500000000stake  --keyring-backend $KEYRING

# setting timeout_commit to 5 minutes
sed -i 's/timeout_commit = "5s"/timeout_commit = "300s"/' $HOME/.simapp/config/config.toml

simd gentx $KEY 100000000stake --keyring-backend $KEYRING --chain-id $CHAINID --gas=0stake

simd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
simd validate-genesis

# start node
simd start
