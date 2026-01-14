package hyprctl

import "encoding/json"

type inputJson struct {
	Label     string
	Type      string `json:",omitempty"`
	Name      string
	Required  bool     `json:",omitempty"`
	Value     *string  `json:",omitempty"`
	MinLength uint     `json:",omitempty"`
	MaxLength uint     `json:",omitempty"`
	Step      float32  `json:",omitempty"`
	Min       string   `json:",omitempty"`
	Max       string   `json:",omitempty"`
	Error     string   `json:",omitempty"`
	Multiple  bool     `json:",omitempty"`
	Options   []Option `json:",omitempty"`
}

func (i Input) MarshalJSON() ([]byte, error) {
	vals := i.Values()
	j := inputJson{
		Label:     i.Label,
		Type:      i.Type,
		Name:      i.Name,
		Required:  i.Required,
		MinLength: i.MinLength,
		MaxLength: i.MaxLength,
		Step:      i.Step,
		Min:       i.Min,
		Max:       i.Max,
		Error:     i.Error,
		Multiple:  i.Multiple,
		Options:   i.Options,
	}
	if len(vals) < 2 {
		j.Value = &vals[0]
	}
	return json.Marshal(j)
}
