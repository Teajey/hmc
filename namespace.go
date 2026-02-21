package hmc

import (
	"encoding/xml"
)

// Namespace provides the XML namespace declaration for hmc elements.
// It should be embedded in your top-level struct and initialized using NS().
//
// Without the namespace, hmc elements (c:Form, c:Input, etc.) will not be
// properly recognized by browsers.
//
// Example:
//
//	type MyPage struct {
//		Namespace hmc.Namespace
//		Title     string
//		Form      hmc.Form[MyFormData]
//	}
//
//	page := MyPage{
//		Namespace: hmc.NS(),
//		Title:     "My Page",
//	}
//
// This struct is omitted during JSON serialization.
type Namespace struct {
	HcXmlns string      `xml:"xmlns:c,attr" json:"-"`
	Docs    xml.Comment `xml:",comment" json:"-"`
}

// NS returns a Namespace with the correct namespace URI and documentation comment.
// Always use this function to initialize the Namespace field in your structs.
func NS() Namespace {
	return Namespace{
		HcXmlns: "https://github.com/Teajey/hmc",
		Docs:    xml.Comment("See an overview of what this XML means at https://github.com/Teajey/hmc/blob/main/README.md "),
	}
}
