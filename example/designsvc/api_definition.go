package designsvc

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("entp", func() {
	Title("Demo: items and purchase indent")
	Description("add items, create purchase indent")
	Host("localhost:8899")
	Scheme("http")
	BasePath("/entp")
	ResponseTemplate(Created, func(pattern string) {
		Description("Resource created")
		Status(201)
		Headers(func() {
			Header("Location", String, "href to created resource", func() {
				Pattern(pattern)
			})
		})
	})
	Origin("*", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		MaxAge(600)
		Credentials()
	})
})
