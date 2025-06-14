package openrtb

import "encoding/json"

// A quantifiable often historical data point about an impression.
type Metric struct {
	// Type of metric being presented using exchange curated string names which should be
	// published to bidders a priori.
	//
	// Required.
	Type string `json:"type" bson:"type"`

	// Number representing the value of the metric. Probabilities must be in the range 0.0 – 1.0.
	//
	// Required.
	Value float64 `json:"value" bson:"value"`

	// Source of the value using exchange curated string names which should be published
	// to bidders a priori.
	//
	// If the exchange itself is the source versus a third party, “EXCHANGE” is recommended.
	Vendor string `json:"vendor,omitempty" bson:"vendor"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (m *Metric) Copy() Metric {

	var ext []byte
	if len(m.Ext) != 0 {
		ext = make([]byte, len(m.Ext))
		copy(ext, m.Ext)
	}

	return Metric{
		Type:   m.Type,
		Value:  m.Value,
		Vendor: m.Vendor,
		Ext:    ext,
	}
}
