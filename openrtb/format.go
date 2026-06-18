package openrtb

import "pkg/pointer"

// An allowed size of a banner.
type Format struct {
	// Width in device independent pixels (DIPS).
	Width int `json:"w,omitempty" bson:"w"`

	// Height in device independent pixels (DIPS).
	Height int `json:"h,omitempty" bson:"h"`

	// Relative width when expressing size as a ratio.
	WidthRatio int `json:"wratio,omitempty" bson:"wratio"`

	// Relative height when expressing size as a ratio.
	HeightRatio int `json:"hratio,omitempty" bson:"hratio"`

	// The minimum width in device independent pixels (DIPS) at which the ad will be
	// displayed the size is expressed as a ratio.
	WidthMin int `json:"wmin,omitempty" bson:"wmin"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext *FormatExt `json:"ext,omitempty" bson:"ext"`
}

type FormatExt struct{}

func (d *FormatExt) copy() FormatExt {
	return FormatExt{}
}
func (f *Format) Copy() Format {

	var ext *FormatExt
	if f.Ext != nil {
		ext = pointer.Pointer(f.Ext.copy())
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
