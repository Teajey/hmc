package hyprctl

import (
	"encoding/xml"
	"fmt"
	"runtime/debug"
)

type Namespace struct {
	HcXmlns string      `xml:"xmlns:c,attr" json:"-"`
	Docs    xml.Comment `xml:",comment" json:"-"`
}

var docs xml.Comment

func init() {
	version := "main"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	docs = xml.Comment(fmt.Sprintf("See documentation for this version of hyprctl at https://github.com/Teajey/hyprctl/blob/%s/README.md", version))
}

func SetNamespace() Namespace {
	return Namespace{
		HcXmlns: "https://github.com/Teajey/hyprctl",
		Docs:    docs,
	}
}
