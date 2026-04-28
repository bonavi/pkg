package openrtb

import "pkg/pointer"

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
	Ext *MetricsExt `json:"ext,omitempty" bson:"ext"`
}

type MetricsExt struct{}

func (m *MetricsExt) copy() MetricsExt {
	return MetricsExt{}
}
func (m *Metric) Copy() Metric {

	var ext *MetricsExt
	if m.Ext != nil {
		ext = pointer.Pointer(m.Ext.copy())
	}

	return Metric{
		Type:   m.Type,
		Value:  m.Value,
		Vendor: m.Vendor,
		Ext:    ext,
	}
}
