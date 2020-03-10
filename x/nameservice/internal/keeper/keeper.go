package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/arjunandra/namservice-cosmos/x/nameservice/internal/types"
)

// Keeper of the nameservice store
type Keeper struct {
	CoinKeeper	types.BankKeeper
	storeKey	sdk.storeKey
	cdc 		*codec.Codec
}

// Keeper Constructor
func newKeeper(coinkeeper bank.Keeper, storekey sdk.storeKey, cdc *codec.Codec) Keeper {
	return Keeper {
		CoinKeeper: coinkeeper,
		storeKey: storekey,
		cdc: cdc
	}
}

// whoIs Getter & Setter
func (k Keeper) setWhoIs(ctx sdk.Context, name string, w types.whoIs) {     

	// No Owner
	if whoIs.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(w))
}

func (k Keeper) getWhoIs(ctx sdk.Context, name string) types.whoIs {
	store := ctx.KVStore(k.storeKey)

	// No whoIs
	if !k.IsNamePresent(ctx, name) {
		return types.newWhoIs()
	}

	bz := store.Get([]byte(name))

	var whoIs types.whoIs

	k.cdc.MustUnmarshalBinaryBare(bz, &whoIs)
	return whoIs
}

func (k Keeper) deleteWhoIs(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))
}

// Name Getter & Setter & Bool & Iterator
func (k Keeper) getName(ctx sdk.Context, name string) string {
	return k.getWhoIs(ctx, name).Value
}

func (k Keeper) setName(ctx sdk.Context, name string, value string) {
	whois := k.getWhoIs(ctx, name)
	whois.Value = value 
	k.setWhoIs(ctx, name, whoIs)
}

func (k Keeper) isNamePresent(ctx sdk.Context, name string) bool {
	store := k.getWhoIs(ctx, name)
	return store.Has([]byte(name))
}

func (k Keeper) getNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// Owner Getter, Setter, & Bool

func (k Keeper) getOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return whois := k.getWhoIs(ctx, name)
}

func (k Keeper) setName(ctx sdk.Context, name string, sdk.AccAddress) {
	whois := k.getWhoIs(ctx, name)
	whois.Owner = owner 
	k.setWhoIs(ctx, name, whoIs)
}

func (k Keeper) hasOwner(ctx sdk.Context, name string) bool {
	return !k.getWhoIs(ctx, name).Owner.Empty()
}

// Price Getter & Setter

func (k Keeper) getPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.getWhoIs(ctx, name).Price
}

func (k Keeper) setPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.getWhoIs(ctx, name).Price
	whois.Price = price
	k.setWhoIs(ctx, name, whois)
}









