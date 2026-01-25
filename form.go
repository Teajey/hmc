package hyprctl

import "encoding/xml"

// Form is analogous to HTML's <form> which represents a state transition that requires input from the client.
// It describes what data is needed, how it should be submitted, and where
// it should be sent.
//
// Elements should contain a struct representing
// and elements semantically related to the form,
// which would usually be [Input], [Link], [Submit], etc.
// but might also be something like `Error string` or `Warning string` fields.
type Form[T any] struct {
	Method   string `json:",omitempty"`
	Elements T
}

func (i Form[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "c:Form"}

	if i.Method != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "Method"}, Value: i.Method})
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	err = e.Encode(i.Elements)
	if err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
