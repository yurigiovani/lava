package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gomath "math"
)

// MsgEnterLotteryList collection of MsgEnterLottery
type MsgEnterLotteryList []*MsgEnterLottery

// IsHighestBet method to check if a msg have the highest ben on the list
func (m MsgEnterLotteryList) IsHighestBet(msg MsgEnterLottery) bool {
	highestBet := math.NewInt(0)

	for _, row := range m {
		if row.Bet.Amount.LTE(highestBet) {
			continue
		}

		highestBet = row.Bet.Amount
	}

	return msg.Bet.Amount.Equal(highestBet)
}

// IsLowestBet
func (m MsgEnterLotteryList) IsLowestBet(msg MsgEnterLottery) bool {
	lowestBet := math.NewInt(gomath.MaxInt64)

	for _, row := range m {
		if row.Bet.Amount.GTE(lowestBet) {
			continue
		}

		lowestBet = row.Bet.Amount
	}

	return msg.Bet.Amount.Equal(lowestBet)
}

func NewMsgEnterLottery(address string, bet sdk.Coin, fee sdk.Coin) MsgEnterLottery {
	return MsgEnterLottery{
		Address: address,
		Bet:     bet,
		Fee:     fee,
	}
}

func (m MsgEnterLottery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if !m.Bet.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid bet")
	}

	minBet := sdk.NewCoin("stake", math.NewInt(GetMinBetLottery()))

	if m.Bet.IsLT(minBet) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("bet must be greater than %d", GetMinBetLottery()))
	}

	return nil
}

func (m MsgEnterLottery) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(m.Address)
	return []sdk.AccAddress{fromAddress}
}
