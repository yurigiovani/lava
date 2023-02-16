package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (k Keeper) EnterLottery(ctx sdk.Context, msg *types.MsgEnterLottery) error {
	ctx.Logger().Info("entering on lottery", msg)

	//lotteryAccount := k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
	lotteryStore := ctx.KVStore(k.storeKey)

	k.setLottery(lotteryStore, 1, msg)

	return nil
}

func (k Keeper) setLottery(store sdk.KVStore, id int64, msg *types.MsgEnterLottery) {
	fmt.Println(store.Has(types.GetLotteryEntriesKey(id)))

	if store.Has(types.GetLotteryEntriesKey(id)) == true {
		fmt.Println("already in lottery")
	}
}
