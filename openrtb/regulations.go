package openrtb

import "pkg/pointer"

// Regulatory conditions in effect for all impressions in this bid request.
type Regulations struct {
	// Flag indicating if this request is subject to the COPPA regulations established by
	// the USA FTC, where:
	//   0 = no;
	//   1 = yes.
	COPPA int `json:"coppa,omitempty" bson:"coppa"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext *RegulationExt `json:"ext,omitempty" bson:"ext"`
}

type RegulationExt struct{}

func (r *RegulationExt) copy() RegulationExt {
	return RegulationExt{}
}

func (r *Regulations) Copy() Regulations {

	var ext *RegulationExt
	if r.Ext != nil {
		ext = pointer.Pointer(r.Ext.copy())
	}

	return Regulations{
		COPPA: r.COPPA,
		Ext:   ext,
	}
}
