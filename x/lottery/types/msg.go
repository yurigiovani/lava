package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgEnterLottery struct {
	User       sdk.AccAddress
	Amount     sdk.Coin
	TokenDenom string
}

func NewMsgEnterLottery(user sdk.AccAddress, amount sdk.Coin, tokenDenom string) MsgEnterLottery {
	return MsgEnterLottery{
		User:       user,
		Amount:     amount,
		TokenDenom: tokenDenom,
	}
}

func (msg MsgEnterLottery) Reset() {
	msg = MsgEnterLottery{}
}

func (msg MsgEnterLottery) String() string {
	return fmt.Sprintf("%v", msg)
}

func (msg MsgEnterLottery) ProtoMessage() {}

func (msg MsgEnterLottery) ValidateBasic() error {
	if msg.User.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing user address")
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}

	if msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	if msg.TokenDenom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "missing token denom")
	}

	return nil
}

func (msg MsgEnterLottery) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.User}
}
