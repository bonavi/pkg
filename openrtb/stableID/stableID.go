package stableID

import "pkg/pointer"

type StableID struct {
	Version string   `json:"version"`
	Value   string   `json:"value"`
	Type    string   `json:"type"`
	Source  string   `json:"source"`
	Privacy *Privacy `json:"privacy"`
}

type Privacy struct {
	ConsentScope []string `json:"consent_scope"`
}

func (o *StableID) Validate() error {
	return nil
}

func (o *StableID) Copy() StableID {

	var privacy *Privacy
	if o.Privacy != nil {
		privacy = pointer.Pointer(o.Privacy.copy())
	}

	return StableID{
		Version: o.Version,
		Value:   o.Value,
		Type:    o.Type,
		Source:  o.Source,
		Privacy: privacy,
	}
}

func (o *Privacy) copy() Privacy {

	var consentScopes []string
	if len(o.ConsentScope) != 0 {
		consentScopes = make([]string, len(o.ConsentScope))
		copy(consentScopes, o.ConsentScope)
	}

	return Privacy{
		ConsentScope: consentScopes,
	}
}
