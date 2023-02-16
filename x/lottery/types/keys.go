package types

const (
	// ModuleName is the name of the lottery module
	ModuleName = "lottery"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the lottery module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the lottery module
	RouterKey = ModuleName

	AttributeValueCategory = ModuleName

	// LotteryFee is the fee that will be charged to user
	LotteryFee = 5

	// MinBet is the minimal value for enter on Lottery
	MinBet = 1
)

// GetMinBetLottery is the function that will return the min bet to enter on Lottery
func GetMinBetLottery() int64 {
	return LotteryFee + MinBet
}

// Keys for distribution store
var (
	KeyLotteryCounter = []byte{0x01}
	KeyLotteryEntries = []byte{0x02}
)

// GetLotteryEntriesKey
func GetLotteryEntriesKey(id int64) []byte {
	return append(KeyLotteryEntries, byte(id))
}
