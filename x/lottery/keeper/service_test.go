package keeper_test

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/keeper"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func createTestApp(t *testing.T) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	return app, ctx
}

func createExpectedKeepers(ctx sdk.Context, app *simapp.SimApp) (types.AccountKeeper, types.BankKeeper) {
	maccPerms := simapp.GetMaccPerms()
	appCodec := simapp.MakeTestEncodingConfig().Codec
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec, app.GetKey(types.StoreKey), app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount, maccPerms, sdk.Bech32MainPrefix,
	)
	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, app.GetKey(types.StoreKey), accountKeeper, app.GetSubspace(banktypes.ModuleName), nil)

	return accountKeeper, bankKeeper
}

func TestCalculatePayoutLowestBet(t *testing.T) {
	app, ctx := createTestApp(t)
	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)

	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(0))
	msgs := types.MsgEnterLotteryList{}

	msgs = append(msgs, &types.MsgEnterLottery{
		Address: addrs[0].String(),
		Bet: sdk.Coin{
			Denom:  "stake",
			Amount: math.NewInt(10),
		},
		Fee: sdk.Coin{
			Denom:  "stake",
			Amount: math.NewInt(1),
		},
	}, &types.MsgEnterLottery{
		Address: addrs[1].String(),
		Bet: sdk.Coin{
			Denom:  "stake",
			Amount: math.NewInt(15),
		},
		Fee: sdk.Coin{
			Denom:  "stake",
			Amount: math.NewInt(1),
		},
	})

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)
	coinsMint := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(27)))

	if err := bankKeeper.(bankkeeper.Keeper).MintCoins(ctx, minttypes.ModuleName, coinsMint); err != nil {
		panic(fmt.Sprintf("some error while minting coins: %s", err))
	}

	if err := bankKeeper.(bankkeeper.Keeper).SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coinsMint); err != nil {
		panic(fmt.Sprintf("some error while sending coins from module lottery: %s", err))
	}

	err1 := k.EnterLottery(ctx, msgs[0])
	err2 := k.EnterLottery(ctx, msgs[1])

	require.Nil(t, err1)
	require.Nil(t, err2)

	isPayable, err := k.Payout(ctx, *msgs[0])

	require.Nil(t, err)
	require.False(t, isPayable, "this msg must not be paid")

	simapp.CheckBalance(t, app, addrs[0], sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(0))))
}

func TestCalculatePayoutLowestHighestBet(t *testing.T) {
	app, ctx := createTestApp(t)
	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)

	addrs := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(0))
	msgs := types.MsgEnterLotteryList{}

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))

	msgs = append(msgs, &msg1, &msg2)

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)
	coinsMint := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(27)))

	if err := bankKeeper.(bankkeeper.Keeper).MintCoins(ctx, minttypes.ModuleName, coinsMint); err != nil {
		panic(fmt.Sprintf("some error while minting coins: %s", err))
	}

	if err := bankKeeper.(bankkeeper.Keeper).SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coinsMint); err != nil {
		panic(fmt.Sprintf("some error while sending coins from module lottery: %s", err))
	}

	err1 := k.EnterLottery(ctx, msgs[0])
	err2 := k.EnterLottery(ctx, msgs[1])

	require.Nil(t, err1)
	require.Nil(t, err2)

	isPayable, err := k.Payout(ctx, *msgs[1])

	require.Nil(t, err)
	require.True(t, isPayable, "this msg must be paid")

	bal := bankKeeper.GetAllBalances(ctx, addrs[1])

	require.Equal(t, math.NewInt(27), bal[0].Amount)
}

func TestCalculatePayoutBetKeepingFeeOnLotteryPool(t *testing.T) {
	app, ctx := createTestApp(t)
	accountKeeper, bankKeeper := createExpectedKeepers(ctx, app)

	addrs := simapp.AddTestAddrs(app, ctx, 3, math.NewInt(0))
	msgs := types.MsgEnterLotteryList{}

	msg1 := types.NewMsgEnterLottery(addrs[0].String(), sdk.NewCoin("stake", math.NewInt(10)), sdk.NewCoin("stake", math.NewInt(1)))
	msg2 := types.NewMsgEnterLottery(addrs[1].String(), sdk.NewCoin("stake", math.NewInt(15)), sdk.NewCoin("stake", math.NewInt(1)))
	msg3 := types.NewMsgEnterLottery(addrs[2].String(), sdk.NewCoin("stake", math.NewInt(25)), sdk.NewCoin("stake", math.NewInt(1)))

	msgs = append(msgs, &msg1, &msg2, &msg3)

	k := keeper.NewKeeper(nil, app.GetKey(types.StoreKey), accountKeeper, bankKeeper)
	coinsMint := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(53)))

	if err := bankKeeper.(bankkeeper.Keeper).MintCoins(ctx, minttypes.ModuleName, coinsMint); err != nil {
		panic(fmt.Sprintf("some error while minting coins: %s", err))
	}

	if err := bankKeeper.(bankkeeper.Keeper).SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, types.ModuleName, coinsMint); err != nil {
		panic(fmt.Sprintf("some error while sending coins from module lottery: %s", err))
	}

	err1 := k.EnterLottery(ctx, msgs[0])
	err2 := k.EnterLottery(ctx, msgs[1])
	err3 := k.EnterLottery(ctx, msgs[2])

	require.Nil(t, err1)
	require.Nil(t, err2)
	require.Nil(t, err3)

	isPayable, err := k.Payout(ctx, *msgs[1])

	require.Nil(t, err)
	require.True(t, isPayable, "this msg must be paid")

	bal := bankKeeper.GetAllBalances(ctx, addrs[1])

	require.Equal(t, math.NewInt(50), bal[0].Amount)
}
