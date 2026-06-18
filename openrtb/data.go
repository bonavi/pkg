package openrtb

import (
	"pkg/pointer"
)

// Collection of additional user targeting data from a specific data source.
type Data struct {
	// Exchange-specific ID for the data provider.
	ID string `json:"id,omitempty" db:"id"`

	// Exchange-specific name for the data provider.
	Name string `json:"name,omitempty" db:"name"`

	// Array of Segment objects that contain the actual data values.
	Segment []Segment `json:"segment,omitempty" db:"segment"`

	// Placeholder for exchange-specific extensions to OpenRTB
	Ext *DataExt `json:"ext,omitempty" db:"ext"`
}

type DataExt struct{}

func (d *DataExt) copy() DataExt {
	return DataExt{}
}

func (d *Data) Copy() Data {

	var segments []Segment
	if len(d.Segment) != 0 {
		segments = make([]Segment, len(d.Segment))
		for i := range d.Segment {
			segments[i] = d.Segment[i].Copy()
		}
	}

	var ext *DataExt
	if d.Ext != nil {
		ext = pointer.Pointer(d.Ext.copy())
	}

	return Data{
		ID:      d.ID,
		Name:    d.Name,
		Segment: segments,
		Ext:     ext,
	}
}

// Specific data point about a user from a specific data source.
type Segment struct {
	// ID of the data segment specific to the data provider.
	ID string `json:"id,omitempty" db:"id"`

	// Name of the data segment specific to the data provider.
	Name string `json:"name,omitempty" db:"name"`

	// String representation of the data segment value.
	Value string `json:"value,omitempty" db:"value"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext *SegmentExt `json:"ext,omitempty" db:"ext"`
}

type SegmentExt struct{}

func (d *SegmentExt) copy() SegmentExt {
	return SegmentExt{}
}

func (s *Segment) Copy() Segment {

	var ext *SegmentExt
	if s.Ext != nil {
		ext = pointer.Pointer(s.Ext.copy())
	}

	return Segment{
		ID:    s.ID,
		Name:  s.Name,
		Value: s.Value,
		Ext:   ext,
	}
}
