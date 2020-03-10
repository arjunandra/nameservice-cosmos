package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Name Value
var minNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

type WhoIs struct {
	Value string			`json:"value"`
	Owner sdk.AccAddress 	`json:"owner"`
	Price sdk.Coins			`json:"price"`
}

func NewWhoIs() WhoIs {
	return WhoIs {
		Price: minNamePrice,
	}
}

// Display whoIs
func (w WhoIs) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s\n Value: %s\n Price: %s`, w.Owner, w.Value, w.Price))
}
