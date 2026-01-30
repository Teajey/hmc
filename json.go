package hmc

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
