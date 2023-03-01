package testutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
)

func GetModuleBalance(bankKeeper bankkeeper.Keeper, accountKeeper types.AccountKeeper, ctx sdk.Context, moduleName string, denom string) sdk.Coin {
	acc := accountKeeper.GetModuleAccount(ctx, moduleName)
	return bankKeeper.GetBalance(ctx, acc.GetAddress(), denom)
}

func GetAccountBalance(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return bankKeeper.GetBalance(ctx, addr, denom)
}
