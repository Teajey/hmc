package hyprctl

import "encoding/json"

type inputJson struct {
	Label     string  `json:"label"`
	Type      string  `json:"type,omitempty"`
	Name      string  `json:"name"`
	Required  bool    `json:"required,omitempty"`
	Value     string  `json:"value"`
	Error     string  `json:"error,omitempty"`
	MinLength uint    `json:"minlength,omitempty"`
	MaxLength uint    `json:"maxlength,omitempty"`
	Step      float32 `json:"step,omitempty"`
	Min       string  `json:"min,omitempty"`
	Max       string  `json:"max,omitempty"`
}

func (i Input) MarshalJSON() ([]byte, error) {
	j := inputJson(i)
	if j.Type == "password" && j.Value != "" {
		j.Value = "********"
	}
	return json.Marshal(j)
}
