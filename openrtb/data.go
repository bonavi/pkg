package openrtb

import "encoding/json"

// Collection of additional user targeting data from a specific data source.
type Data struct {
	// Exchange-specific ID for the data provider.
	ID string `json:"id,omitempty" db:"id"`

	// Exchange-specific name for the data provider.
	Name string `json:"name,omitempty" db:"name"`

	// Array of Segment objects that contain the actual data values.
	Segment []Segment `json:"segment,omitempty" db:"segment"`

	// Placeholder for exchange-specific extensions to OpenRTB
	Ext json.RawMessage `json:"ext,omitempty" db:"ext"`
}

func (d *Data) Copy() Data {

	var segments []Segment
	if len(d.Segment) != 0 {
		segments = make([]Segment, len(d.Segment))
		for i := range d.Segment {
			segments[i] = d.Segment[i].Copy()
		}
	}

	var ext []byte
	if len(d.Ext) != 0 {
		ext = make([]byte, len(d.Ext))
		copy(ext, d.Ext)
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
	Ext json.RawMessage `json:"ext,omitempty" db:"ext"`
}

func (s *Segment) Copy() Segment {

	var ext []byte
	if len(s.Ext) != 0 {
		ext = make([]byte, len(s.Ext))
		copy(ext, s.Ext)
	}

	return Segment{
		ID:    s.ID,
		Name:  s.Name,
		Value: s.Value,
		Ext:   ext,
	}
}
