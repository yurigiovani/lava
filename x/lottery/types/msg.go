package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgEnterLotteryList colletion of MsgEnterLottery
type MsgEnterLotteryList = []*MsgEnterLottery

func NewMsgEnterLottery(address string, bet sdk.Coin, fee sdk.Coin) MsgEnterLottery {
	return MsgEnterLottery{
		Address: address,
		Bet:     bet,
		Fee:     fee,
	}
}

func (msg MsgEnterLottery) GetAccAddress() (sdk.AccAddress, error) {
	acc, err := sdk.AccAddressFromBech32(msg.Address)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (msg MsgEnterLottery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if !msg.Bet.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid bet")
	}

	minBet := sdk.NewCoin("stake", math.NewInt(GetMinBetLottery()))

	if msg.Bet.IsLT(minBet) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("bet must be greater than %d", GetMinBetLottery()))
	}

	return nil
}

func (msg MsgEnterLottery) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Address)
	return []sdk.AccAddress{fromAddress}
}
