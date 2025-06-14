package openrtb

import "encoding/json"

// Details for a banner impression (incl. in-banner video) or video companion ad.
type Banner struct {
	// Array of format objects representing the banner sizes permitted.
	//
	// If none are specified, then use of the h and w attributes is highly recommended.
	Formats []Format `json:"format,omitempty" bson:"format"`

	// Exact width in device independent pixels (DIPS).
	//
	// Recommended if no format objects are specified.
	Width int `json:"w,omitempty" bson:"w"`

	// Exact height in device independent pixels (DIPS).
	//
	// Recommended if no format objects are specified
	Height int `json:"h,omitempty" bson:"h"`

	// Maximum width in device independent pixels (DIPS).
	// Deprecated: deprecated in favor of the format array
	WidthMax int `json:"wmax,omitempty" bson:"wmax"`

	// Maximum height in device independent pixels (DIPS).
	//
	// Deprecated: deprecated in favor of the format array.
	HeightMax int `json:"hmax,omitempty" bson:"hmax"`

	// Minimum width in device independent pixels (DIPS).
	//
	// Deprecated: deprecated in favor of the format array.
	WidthMin int `json:"wmin,omitempty" bson:"wmin"`

	// Minimum height in device independent pixels (DIPS).
	//
	// Deprecated: deprecated in favor of the format array.
	HeightMin int `json:"hmin,omitempty" bson:"hmin"`

	// Blocked banner ad types.
	BlockedTypes []BannerType `json:"btype,omitempty" bson:"btype"`

	// Blocked creative attributes.
	BlockedAttrs []CreativeAttribute `json:"battr,omitempty" bson:"battr"`

	// Ad position on screen.
	Position AdPosition `json:"pos,omitempty" bson:"pos"`

	// Content MIME types supported. Popular MIME types may include
	// “application/x-shockwave-flash”, “image/jpg”, and “image/gif”.
	MIMEs []string `json:"mimes,omitempty" bson:"mimes"`

	// Indicates if the banner is in the top frame as opposed to an iframe, where:
	//    0 = no;
	//    1 = yes.
	TopFrame int `json:"topframe,omitempty" bson:"topframe"`

	// Directions in which the banner may expand.
	ExpDirs []ExpDir `json:"expdir,omitempty" bson:"expdir"`

	// List of supported API frameworks for this impression.
	//
	// If an API is not explicitly listed, it is assumed not to be supported.
	APIs []APIFramework `json:"api,omitempty" bson:"api"`

	// Unique identifier for this banner object.
	//
	// Recommended when Banner objects are used with a Video object to represent an array
	// of companion ads. Values usually start at 1 and increase with each object;
	// should be unique within an impression.
	ID string `json:"id,omitempty" bson:"id"`

	// Relevant only for Banner objects used with a Video object in an array of companion ads.
	// Indicates the companion banner rendering mode relative to the associated video, where:
	//   0 = concurrent;
	//   1 = end-card.
	VCM int `json:"vcm,omitempty" bson:"vcm"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (b *Banner) Copy() Banner {

	var formats []Format
	if len(b.Formats) != 0 {
		formats = make([]Format, len(b.Formats))
		for i := range b.Formats {
			formats[i] = b.Formats[i].Copy()
		}
	}

	var blockedTypes []BannerType
	if len(b.BlockedTypes) != 0 {
		blockedTypes = make([]BannerType, len(b.BlockedTypes))
		copy(blockedTypes, b.BlockedTypes)
	}

	var blockedAttrs []CreativeAttribute
	if len(b.BlockedAttrs) != 0 {
		blockedAttrs = make([]CreativeAttribute, len(b.BlockedAttrs))
		copy(blockedAttrs, b.BlockedAttrs)
	}

	var expDirs []ExpDir
	if len(b.ExpDirs) != 0 {
		expDirs = make([]ExpDir, len(b.ExpDirs))
		copy(expDirs, b.ExpDirs)
	}

	var mimes []string
	if len(b.MIMEs) != 0 {
		mimes = make([]string, len(b.MIMEs))
		copy(mimes, b.MIMEs)
	}

	var apis []APIFramework
	if len(b.APIs) != 0 {
		apis = make([]APIFramework, len(b.APIs))
		copy(apis, b.APIs)
	}

	var ext []byte
	if len(b.Ext) != 0 {
		ext = make([]byte, len(b.Ext))
		copy(ext, b.Ext)
	}

	return Banner{
		Formats:      formats,
		Width:        b.Width,
		Height:       b.Height,
		WidthMax:     b.WidthMax,
		HeightMax:    b.HeightMax,
		WidthMin:     b.WidthMin,
		HeightMin:    b.HeightMin,
		BlockedTypes: blockedTypes,
		BlockedAttrs: blockedAttrs,
		Position:     b.Position,
		MIMEs:        mimes,
		TopFrame:     b.TopFrame,
		ExpDirs:      expDirs,
		APIs:         apis,
		ID:           b.ID,
		VCM:          b.VCM,
		Ext:          ext,
	}
}
