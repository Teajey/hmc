package hmc

import (
	"encoding/xml"
	"fmt"
)

// Namespace should be embedded in the top-level element to provide context about what the hyper control elements are.
//
// SetNamespace should be used to populate the values to their default.
type Namespace struct {
	HcXmlns string      `xml:"xmlns:c,attr" json:"-"`
	Docs    xml.Comment `xml:",comment" json:"-"`
}

var docs xml.Comment

const repo string = "github.com/Teajey/hmc"

func init() {
	docs = xml.Comment(fmt.Sprintf("See an overview of what this XML means at https://%s/blob/main/README.md ", repo))
}

// SetNamespace provides a default setting for the Namespace struct.
func SetNamespace() Namespace {
	return Namespace{
		HcXmlns: "https://github.com/Teajey/hmc",
		Docs:    docs,
	}
}
