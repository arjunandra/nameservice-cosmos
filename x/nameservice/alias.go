package nameservice

import (
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/keeper"
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewMsgSetName 		= types.NewMsgSetName
	NewMsgBuyName 		= types.NewMsgBuyName
	NewMsgDeleteName 	= types.NewMsgDeleteName
	NewWhoIs			= types.NewWhoIs
	RegisterCodec       = types.RegisterCodec
	// TODO: Fill out function aliases

	// variable aliases
	ModuleCdc     = types.ModuleCdc
)

type (
	Keeper       	= keeper.Keeper
	Params       	= types.Params
	MsgSetName	 	= types.MsgSetName
	MsgBuyName 	 	= types.MsgBuyName
	MsgDeleteName	= types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames	= types.QueryResNames
	whoIs			= types.WhoIs
)
