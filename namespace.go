package hyprctl

type Namespace struct {
	HcXmlns string `xml:"xmlns:c,attr" json:"-"`
}

func NewNamespace() Namespace {
	return Namespace{
		HcXmlns: "https://github.com/Teajey/hyprctl/blob/v0.3.0/README.md",
	}
}
