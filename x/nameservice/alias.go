package nameservice

import (
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/keeper"
	"github.com/arjunandra/nameservice-cosmos/x/nameservice/internal/types"
)

// Constants Exposed From Internal Package

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
)

// Functions Aliases

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewMsgSetName 		= types.NewMsgSetName
	NewMsgBuyName 		= types.NewMsgBuyName
	NewMsgDeleteName 	= types.NewMsgDeleteName
	NewWhoIs			= types.NewWhoIs
	RegisterCodec       = types.RegisterCodec
)

// Variable Aliases

var (
	ModuleCdc     = types.ModuleCdc
)

// Required Structures 

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
