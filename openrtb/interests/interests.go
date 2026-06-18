package interests

type Interests struct {
	Version string `json:"version"`
	Value   string `json:"value"`
	Source  string `json:"source"`
}

func (o *Interests) Validate() error {
	return nil
}

func (o *Interests) Copy() Interests {
	return Interests{
		Version: o.Version,
		Value:   o.Value,
		Source:  o.Source,
	}
}
