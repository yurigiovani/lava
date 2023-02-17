package keeper

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	authKeeper stakingtypes.AccountKeeper
	bankKeeper stakingtypes.BankKeeper
}

// NewKeeper creates a new lottery Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// store returns a lottery store.
func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

// EnterLottery receive types.MsgEnterLottery and set a new entry for current lottery
func (k Keeper) EnterLottery(ctx sdk.Context, msg *types.MsgEnterLottery) error {
	ctx.Logger().With("address", msg.Address).Info("started to enter on lottery", msg)
	defer ctx.Logger().Info("finished to enter on lottery", msg)

	if k.hasEntryLottery(ctx, 1, msg.Address) {
		err := errors.New("address is already participating from lottery")
		k.Logger(ctx).Error(err.Error())

		return err
	}

	k.setLottery(ctx, 1, msg)

	return nil
}

func (k Keeper) setLottery(ctx sdk.Context, id int64, msg *types.MsgEnterLottery) {
	lotteryStore := k.store(ctx)
	entries, err := k.getLottery(ctx, id)

	if err != nil {
		return
	}

	entries = append(entries, msg)
	bentries, err := amino.MarshalBinaryLengthPrefixed(entries)

	if err != nil {
		fmt.Println("some error while MarshalBinary", err)
	}

	lotteryStore.Set(types.GetLotteryEntriesKey(id), bentries)
}

func (k Keeper) getLottery(ctx sdk.Context, id int64) (types.MsgEnterLotteryList, error) {
	var entries types.MsgEnterLotteryList
	lotteryStore := k.store(ctx)

	if err := amino.UnmarshalBinaryLengthPrefixed(lotteryStore.Get(types.GetLotteryEntriesKey(id)), &entries); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error to get lottery: %s", err))
		return nil, err
	}

	return entries, nil
}

func (k Keeper) hasEntryLottery(ctx sdk.Context, id int64, address string) bool {
	var msgs types.MsgEnterLotteryList
	lotteryStore := k.store(ctx)
	iterator := sdk.KVStorePrefixIterator(lotteryStore, types.GetLotteryEntriesKey(id))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		err := amino.UnmarshalBinaryLengthPrefixed(lotteryStore.Get(iterator.Key()), &msgs)

		if err != nil {
			fmt.Println("error to amino.UnmarshalBinaryLengthPrefixed", err)
			return false
		}
	}

	for _, msg := range msgs {
		if msg.Address != address {
			continue
		}

		return true
	}

	return false
}
