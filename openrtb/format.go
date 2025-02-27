package openrtb

import "encoding/json"

// An allowed size of a banner.
type Format struct {
	// Width in device independent pixels (DIPS).
	Width int `json:"w" bson:"w"`

	// Height in device independent pixels (DIPS).
	Height int `json:"h" bson:"h"`

	// Relative width when expressing size as a ratio.
	WidthRatio int `json:"wratio" bson:"wratio"`

	// Relative height when expressing size as a ratio.
	HeightRatio int `json:"hratio" bson:"hratio"`

	// The minimum width in device independent pixels (DIPS) at which the ad will be
	// displayed the size is expressed as a ratio.
	WidthMin int `json:"wmin" bson:"wmin"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (f *Format) Copy() Format {

	var ext []byte
	if len(f.Ext) != 0 {
		ext = make([]byte, len(f.Ext))
		copy(ext, f.Ext)
	}

	return Format{
		Width:       f.Width,
		Height:      f.Height,
		WidthRatio:  f.WidthRatio,
		HeightRatio: f.HeightRatio,
		WidthMin:    f.WidthMin,
		Ext:         ext,
	}
}
