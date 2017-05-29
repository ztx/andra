package designsvc

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// BottlePayload defines the data structure used in the create bottle request body.
// It is also the base type for the bottle media type used to render bottles.
var PrLinePayload = Type("PrLinePayload", func() {
	// Attribute("name", String, func() {
	// 	MinLength(2)
	// 	Example("Number 8")
	// })
	Attribute("id", Integer, "ID of the item")
	Attribute("itemCode", String, "Item Code")
	Attribute("sl", Integer, "serial number")
	Attribute("qty", Integer, "Quantity")
	Attribute("price", Integer, "price")
})
