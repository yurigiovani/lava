package keeper

import (
	"context"
	"errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/lottery/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the lottery MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

//EnterLottery will receive a broadcasted types.MsgEnterLottery and execute the process
func (k msgServer) EnterLottery(goCtx context.Context, msg *types.MsgEnterLottery) (*types.MsgEnterLotteryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Info("entering into a lottery")

	defer func() {
		if msg.Bet.Amount.IsInt64() {
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", "send"},
				float32(msg.Bet.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Bet.Denom)},
			)
		}
	}()

	propAccAddr := sdk.AccAddress(ctx.BlockHeader().ProposerAddress)
	senderAddr, _ := sdk.AccAddressFromBech32(msg.Address)

	// checking if proposer is trying to enter on lottery
	if propAccAddr.Equals(senderAddr) {
		err := errors.New("proposer could not be a sender")
		k.Logger(ctx).Error(err.Error())
		return nil, err
	}

	if err := k.Keeper.EnterLottery(ctx, msg); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &types.MsgEnterLotteryResponse{}, nil
}
