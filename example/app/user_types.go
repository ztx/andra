// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/ztx/entp/designsvc
// --out=$(GOPATH)/src/github.com/ztx/entp
// --version=v1.1.0-dirty
//
// API "entp": Application User Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package app

// prLinePayload user type.
type prLinePayload struct {
	// ID of the item
	ID *int `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Item Code
	ItemCode *string `form:"itemCode,omitempty" json:"itemCode,omitempty" xml:"itemCode,omitempty"`
	// price
	Price *int `form:"price,omitempty" json:"price,omitempty" xml:"price,omitempty"`
	// Quantity
	Qty *int `form:"qty,omitempty" json:"qty,omitempty" xml:"qty,omitempty"`
	// serial number
	Sl *int `form:"sl,omitempty" json:"sl,omitempty" xml:"sl,omitempty"`
}

// Publicize creates PrLinePayload from prLinePayload
func (ut *prLinePayload) Publicize() *PrLinePayload {
	var pub PrLinePayload
	if ut.ID != nil {
		pub.ID = ut.ID
	}
	if ut.ItemCode != nil {
		pub.ItemCode = ut.ItemCode
	}
	if ut.Price != nil {
		pub.Price = ut.Price
	}
	if ut.Qty != nil {
		pub.Qty = ut.Qty
	}
	if ut.Sl != nil {
		pub.Sl = ut.Sl
	}
	return &pub
}

// PrLinePayload user type.
type PrLinePayload struct {
	// ID of the item
	ID *int `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Item Code
	ItemCode *string `form:"itemCode,omitempty" json:"itemCode,omitempty" xml:"itemCode,omitempty"`
	// price
	Price *int `form:"price,omitempty" json:"price,omitempty" xml:"price,omitempty"`
	// Quantity
	Qty *int `form:"qty,omitempty" json:"qty,omitempty" xml:"qty,omitempty"`
	// serial number
	Sl *int `form:"sl,omitempty" json:"sl,omitempty" xml:"sl,omitempty"`
}
