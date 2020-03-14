package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Structure Declarations

type MsgSetName struct {
	Name string				`json:"name"`
	Value string			`json:"value"`	
	Owner sdk.AccAddress	`json:"owner"`	
}

type MsgBuyName struct {
	Name string				`json:"name"`
	Bid sdk.Coins			`json:"bid"`
	Buyer sdk.AccAddress	`json:"buyer"`
}

type MsgDeleteName struct {
	Name string				`json:"name"`
	Owner sdk.AccAddress	`json:"owner"`
}

// Message Constructors

func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName {
		Name: name,
		Value: value,
		Owner: owner,
	}
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName {
		Name: name,
		Bid: bid,
		Buyer: buyer,
	}
}

func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName {
		Name: name,
		Owner: owner,
	}
}

// Message Route Declarations

func (msg MsgSetName) Route() string { return RouterKey }
func (msg MsgBuyName) Route() string { return RouterKey }
func (msg MsgDeleteName) Route() string { return RouterKey }

// Message Type Declarations

func (msg MsgSetName) Type() string { return "set_name" }
func (msg MsgBuyName) Type() string {return "buy_name"}
func (msg MsgDeleteName) Type() string { return "delete_name" }

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

func (msg MsgDeleteName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}

	return nil
}

// Message Sign Bytes Getter

func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)) 
}

// Message Signers Getter

func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

