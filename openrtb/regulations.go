package openrtb

import "encoding/json"

// Regulatory conditions in effect for all impressions in this bid request.
type Regulations struct {
	// Flag indicating if this request is subject to the COPPA regulations established by
	// the USA FTC, where:
	//   0 = no;
	//   1 = yes.
	COPPA int `json:"coppa,omitempty" bson:"coppa"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (r *Regulations) Copy() Regulations {

	var ext []byte
	if len(r.Ext) != 0 {
		ext = make([]byte, len(r.Ext))
		copy(ext, r.Ext)
	}

	return Regulations{
		COPPA: r.COPPA,
		Ext:   ext,
	}
}
