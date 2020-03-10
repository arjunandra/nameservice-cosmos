package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"
)

// Keeper of the nameservice store
type Keeper struct {
	CoinKeeper	types.BankKeeper
	storeKey	sdk.StoreKey
	cdc 		*codec.Codec
}

// Keeper Constructor
func NewKeeper(coinkeeper types.BankKeeper, storekey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper {
		CoinKeeper: coinkeeper,
		storeKey: storekey,
		cdc: cdc,
	}
}

// whoIs Getter & Setter
func (k Keeper) SetWhoIs(ctx sdk.Context, name string, w types.WhoIs) {     

	// No Owner
	if w.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(w))
}

func (k Keeper) GetWhoIs(ctx sdk.Context, name string) types.WhoIs {
	store := ctx.KVStore(k.storeKey)

	// No whoIs
	if !k.IsNamePresent(ctx, name) {
		return types.NewWhoIs()
	}

	bz := store.Get([]byte(name))

	var whoIs types.WhoIs

	k.cdc.MustUnmarshalBinaryBare(bz, &whoIs)
	return whoIs
}

func (k Keeper) DeleteWhoIs(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))
}

// Name Getter & Setter & Bool & Iterator
func (k Keeper) GetName(ctx sdk.Context, name string) string {
	return k.GetWhoIs(ctx, name).Value
}

func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhoIs(ctx, name)
	whois.Value = value 
	k.SetWhoIs(ctx, name, whois)
}

func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// Owner Getter, Setter, & Bool

func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhoIs(ctx, name).Owner
}

func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhoIs(ctx, name)
	whois.Owner = owner 
	k.SetWhoIs(ctx, name, whois)
}

func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhoIs(ctx, name).Owner.Empty()
}

// Price Getter & Setter

func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhoIs(ctx, name).Price
}

func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhoIs(ctx, name)
	whois.Price = price
	k.SetWhoIs(ctx, name, whois)
}









