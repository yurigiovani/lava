package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	authKeeper stakingtypes.AccountKeeper
	bankKeeper stakingtypes.BankKeeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ak stakingtypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		authKeeper: ak,
		bankKeeper: bk,
	}
}

func (k Keeper) EnterLottery(ctx context.Context, msg *types.MsgEnterLottery) error {
	fmt.Println(msg.String())

	return nil
}
