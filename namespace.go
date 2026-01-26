package hyprctl

type Namespace struct {
	HcXmlns string `xml:"xmlns:c,attr" json:"-"`
}

func SetNamespace() Namespace {
	return Namespace{
		HcXmlns: "https://github.com/Teajey/hyprctl/blob/v0.4.0/README.md",
	}
}
