package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Name Value
var minNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

type whoIs struct {
	Value string			`json:"value"`
	Owner sdk.AccAddress 	`json:"owner"`
	Price sdk.Coins			`json:"price"`
}

func newWhoIs() whoIs {
	return whoIs {
		Price: minNamePrice
	}
}

// Display whoIs
func (w whoIs) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s\n Value: %s\n Price: %s`, w.Owner, w.Value, w.Price))
}
