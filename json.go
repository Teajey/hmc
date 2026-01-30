package hmc

type inputJson struct {
	Label     string  `json:"-"`
	Type      string  `json:"-"`
	Name      string  `json:"-"`
	Error     string  `json:"error"`
	Required  bool    `json:"-"`
	Value     string  `json:"value"`
	MinLength uint    `json:"-"`
	MaxLength uint    `json:"-"`
	Step      float32 `json:"-"`
	Min       string  `json:"-"`
	Max       string  `json:"-"`
}
