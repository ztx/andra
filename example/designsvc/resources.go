package designsvc

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("item", func() {
	DefaultMedia(Item)
	BasePath("/items")
	Action("list", func() {
		Routing(
			GET(""),
		)
		Description("Retrive all items")
		Response(OK, CollectionOf(Item))
	})
	Action("show", func() {
		Routing(
			GET("/:itemCode"),
		)
		Description("Retrive a task with given id")
		Params(func() {
			Param("itemCode", String, "Item Code")
		})
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Create new Item")
		Payload(func() {
			Member("code")
			Member("name")
			Member("uom")
			Required("code", "name", "uom")
		})
		Response(Created, "/items/[0-9A-Za-z]+")
		Response(BadRequest, ErrorMedia)
	})
})

var _ = Resource("pr", func() {
	DefaultMedia(Pr)
	BasePath("/prs")
	Action("list", func() {
		Routing(
			GET(""),
		)
		Description("Retrive all pr s")
		Response(OK, CollectionOf(Pr))
	})
	Action("show", func() {
		Routing(
			GET("/:prNum"),
		)
		Description("Retrive a pr by pr_num")
		Params(func() {
			Param("prNum", String, "pr num")
		})
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Create new Pr")
		Payload(func() {

			Member("prDate")
			Member("prNum")
			Required("prNum", "prDate")
		})
		Response(Created, "/prs/[0-9A-Za-z]+")
		Response(BadRequest, ErrorMedia)
	})
	Action("AddLine", func() {
		Routing(
			POST("/addline"),
		)
		Description("Add a Pr line")
		Payload(PrLinePayload, func() {

			Required("sl", "itemCode", "qty", "price")
		})
		Response(Created, "/prs/[0-9A-Za-z]+/prlines/[0-9A-Za-z]+")
		Response(BadRequest, ErrorMedia)
	})

})

var _ = Resource("PrLine", func() {
	DefaultMedia(Pr)
	BasePath("/prlines")
	Parent("pr")
	Action("list", func() {
		Routing(
			GET(""),
		)
		Description("Retrive all prlines")
		Response(OK, CollectionOf(PrDetails))
	})
	Action("show", func() {
		Routing(
			GET("/:prLineNum"),
		)
		Description("Retrive a pr line by pr_line_num")
		Params(func() {
			Param("prLineNum", String, "pr line num")
		})
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})

	Action("create", func() {
		Routing(
			POST(""),
		)
		Description("Create new Prline")
		Payload(func() {

			Member("sl", Integer)
			Member("itemCode", String)
			Member("qty", Integer)
			Member("price", Integer)
			Required("sl", "itemCode", "qty", "price")
		})
		Response(Created, "/prs/[0-9A-Za-z]+/prLines/[0-9A-Za-z]+")
		Response(BadRequest, ErrorMedia)
	})
	Action("aprove", func() {
		Routing(
			POST("/aprove/:qty"),
		)
		Description("Approve PR line with qty")
		Payload(func() {

			Member("sl", Integer)
			Member("itemCode", String)
			Member("qty", Integer)
			Member("price", Integer)
			Required("sl", "itemCode", "qty", "price")
		})
		Response(OK)
		Response(BadRequest, ErrorMedia)
	})
})
