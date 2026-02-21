package hmc

type inputJson struct {
	Error string `json:"error"`
	Value string `json:"value"`
}

type inputJsonDisabled struct {
	Disabled bool `json:"disabled"`
}
