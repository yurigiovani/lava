package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgEnterLottery(address string, coin sdk.Coin) MsgEnterLottery {
	return MsgEnterLottery{
		Address: address,
		Amount:  coin,
	}
}

func (msg MsgEnterLottery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}

	if msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount cannot be zero")
	}

	return nil
}

func (msg MsgEnterLottery) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Address)
	return []sdk.AccAddress{fromAddress}
}
