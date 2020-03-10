package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Describe your actions, these will implment the interface of `sdk.Msg`
/*
verify interface at compile time
var _ sdk.Msg = &Msg<Action>{}

Msg<Action> - struct for unjailing jailed validator
type Msg<Action> struct {
	ValidatorAddr sdk.ValAddress `json:"address" yaml:"address"` // address of the validator operator
}

NewMsg<Action> creates a new Msg<Action> instance
func NewMsg<Action>(validatorAddr sdk.ValAddress) Msg<Action> {
	return Msg<Action>{
		ValidatorAddr: validatorAddr,
	}
}

const <action>Const = "<action>"

// nolint
func (msg Msg<Action>) Route() string { return RouterKey }
func (msg Msg<Action>) Type() string  { return <action>Const }
func (msg Msg<Action>) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

GetSignBytes gets the bytes for the message signer to sign on
func (msg Msg<Action>) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

ValidateBasic validity check for the AnteHandler
func (msg Msg<Action>) ValidateBasic() error {
	if msg.ValidatorAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address"
	}
	return nil
}
*/


// MsgSetName

type MsgSetName struct {
	Name string				`json:"name"`
	Value string			`json:"value"`	
	Owner sdk.AccAddress	`json:"owner"`	
}

// Constructor for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName{
	return MsgSetName {

	}
}

func (msg MsgSetName) Route() string { return RouterKey }

func (msg MsgSetName) Type() string { return "set_name" }

// Stateless Checks
func (msg MsgSetName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, msg.Owner.String())
	}

	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name and/or Value cannot be empty")
	}

	return nil
}

func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBuyName

type MsgBuyName struct {
	Name string				`json:"name"`
	Bid sdk.Coins			`json:"bid"`
	Buyer sdk.AccAddress	`json:"buyer"`
}

// Constructor For MsgBuyName

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) {
	return MsgBuyName {
		Name: name,
		Bid: bid,
		Buyer: buyer
	}
}

func (msg MsgBuyName) Route() string { return RouterKey }
func (msg MsgBuyName) Type() string {return "buy_name"}

func (msg MsgBuyName) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}

	if !msg.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}

	return nil
}

func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// MsgDeleteName

type MsgDeleteName struct {
	Name string				`json:"name"`
	Owner sdk.AccAddress	`json:"owner"`
}

// Constructor for MsgDeleteName

func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName {
		Name: name,
		Owner: owner
	}
}

func (msg MsgDeleteName) Route() string { return RouterKey }
func (msg MsgDeleteName) Type() string { return "delete_name" }

func (msg MsgDeleteName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}

	return nil
}

func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)) 
}

func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

