package hyprctl

import (
	"encoding/xml"
	"fmt"
)

type Namespace struct {
	HcXmlns string      `xml:"xmlns:c,attr" json:"-"`
	Docs    xml.Comment `xml:",comment" json:"-"`
}

var docs xml.Comment

const repo string = "github.com/Teajey/hyprctl"

func init() {
	docs = xml.Comment(fmt.Sprintf("See an overview of what this XML means at https://%s/blob/main/README.md ", repo))
}

func SetNamespace() Namespace {
	return Namespace{
		HcXmlns: "https://github.com/Teajey/hyprctl",
		Docs:    docs,
	}
}
