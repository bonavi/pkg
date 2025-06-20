package openrtb

import (
	"encoding/json"
	"errors"

	"pkg/decimal"
	"pkg/pointer"
)

// Validation errors.
var (
	ErrInvalidImpNoID     = errors.New("impression has no ID")
	ErrInvalidMultiAssets = errors.New("impression has multiple assets")
)

// Container for the description of a specific impression.
type Impression struct {
	// A unique identifier for this impression within the context of the bid request
	// (typically, starts with 1 and increments).
	//
	// Required.
	ID string `json:"id" bson:"id"`

	// An array of Metric object.
	Metric []Metric `json:"metric,omitempty" bson:"metric"`

	// A Banner object.
	//
	// Required if this impression is offered as a banner ad opportunity.
	Banner *Banner `json:"banner,omitempty" bson:"banner"`

	// A Video object.
	//
	// Required if this impression is offered as a video ad opportunity.
	Video *Video `json:"video,omitempty" bson:"video"`

	// An Audio object.
	//
	// Required if this impression is offered as an audio ad opportunity.
	Audio *Audio `json:"audio,omitempty" bson:"audio"`

	// A Native object.
	//
	// Required if this impression is offered as a native ad opportunity.
	Native *Native `json:"native,omitempty" bson:"native"`

	// A PMP object containing any private marketplace deals in effect for this impression.
	PMP *PMP `json:"pmp,omitempty" bson:"pmp"`

	// Name of ad mediation partner, SDK technology, or player responsible for rendering
	// ad (typically video or mobile). Used by some ad servers to customize ad code by partner.
	//
	// Recommended for video and/or apps.
	DisplayManager string `json:"displaymanager,omitempty" bson:"displaymanager"`

	// Version of ad mediation partner, SDK technology, or player responsible for rendering
	// ad (typically video or mobile). Used by some ad servers to customize ad code by partner.
	//
	// Recommended for video and/or apps.
	DisplayManagerVersion string `json:"displaymanagerver,omitempty" bson:"displaymanagerver"`

	//    1 = the ad is interstitial or full screen;
	//    0 = not interstitial.
	//
	// Default 0.
	Interstitial int `json:"instl,omitempty" bson:"instl"`

	// Identifier for specific ad placement or ad tag that was used to initiate the auction.
	// This can be useful for debugging of any issues, or for optimization by the buyer.
	TagID string `json:"tagid,omitempty" bson:"tagid"`

	// Minimum bid for this impression expressed in CPM.
	//
	// Default 0.0.
	BidFloor decimal.Decimal `json:"bidfloor,omitempty" bson:"bidfloor"`

	// Currency specified using ISO-4217 alpha codes. This may be different from bid currency
	// returned by bidder if this is allowed by the exchange.
	//
	// Default USD.
	BidFloorCurrency string `json:"bidfloorcur,omitempty" bson:"bidfloorcur"`

	// Indicates the type of browser opened upon clicking the creative in an app, where:
	//    0 = embedded;
	//    1 = native.
	// Note that the Safari View Controller in iOS 9.x devices is considered a native browser
	// for purposes of this attribute.
	ClickBrowser int `json:"clickbrowser,omitempty" bson:"clickbrowser"`

	// Flag to indicate if the impression requires secure HTTPS URL creative assets and markup,
	// where:
	//    0 = non-secure;
	//    1 = secure.
	// If omitted, the secure state is unknown, but non-secure HTTP support can be assumed.
	Secure int `json:"secure,omitempty" bson:"secure"`

	// Array of exchange-specific names of supported iframe busters.
	IFrameBusters []string `json:"iframebuster,omitempty" bson:"iframebuster"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (imp *Impression) AssetCount() int {
	var n int
	if imp.Banner != nil {
		n++
	}
	if imp.Video != nil {
		n++
	}
	if imp.Audio != nil {
		n++
	}
	if imp.Native != nil {
		n++
	}
	return n
}

// Validate the impression object.
func (imp *Impression) Validate() error {
	if imp.ID == "" {
		return ErrInvalidImpNoID
	} else if count := imp.AssetCount(); count > 1 {
		return ErrInvalidMultiAssets
	}

	if imp.Video != nil {
		if err := imp.Video.Validate(); err != nil {
			return err
		}
	}

	if imp.Audio != nil {
		if err := imp.Audio.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type jsonImpression Impression

// UnmarshalJSON custom unmarshalling.
func (imp *Impression) UnmarshalJSON(data []byte) error {
	j := jsonImpression{BidFloorCurrency: "USD"}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	*imp = (Impression)(j)
	return nil
}

func (imp *Impression) Copy() Impression {

	var metric []Metric
	if len(imp.Metric) != 0 {
		metric = make([]Metric, len(imp.Metric))
		for i := range imp.Metric {
			metric[i] = imp.Metric[i].Copy()
		}
	}

	var banner *Banner
	if imp.Banner != nil {
		banner = pointer.Pointer(imp.Banner.Copy())
	}

	var video *Video
	if imp.Video != nil {
		video = pointer.Pointer(imp.Video.Copy())
	}

	var audio *Audio
	if imp.Audio != nil {
		audio = pointer.Pointer(imp.Audio.Copy())
	}

	var native *Native
	if imp.Native != nil {
		native = pointer.Pointer(imp.Native.Copy())
	}

	var pmp *PMP
	if imp.PMP != nil {
		pmp = pointer.Pointer(imp.PMP.Copy())
	}

	var iframeBusters []string
	if len(imp.IFrameBusters) != 0 {
		iframeBusters = make([]string, len(imp.IFrameBusters))
		copy(iframeBusters, imp.IFrameBusters)
	}

	var ext []byte
	if len(imp.Ext) != 0 {
		ext = make([]byte, len(imp.Ext))
		copy(ext, imp.Ext)
	}

	return Impression{
		ID:                    imp.ID,
		Metric:                metric,
		Banner:                banner,
		Video:                 video,
		Audio:                 audio,
		Native:                native,
		PMP:                   pmp,
		DisplayManager:        imp.DisplayManager,
		DisplayManagerVersion: imp.DisplayManagerVersion,
		Interstitial:          imp.Interstitial,
		TagID:                 imp.TagID,
		BidFloor:              imp.BidFloor,
		BidFloorCurrency:      imp.BidFloorCurrency,
		ClickBrowser:          imp.ClickBrowser,
		Secure:                imp.Secure,
		IFrameBusters:         iframeBusters,
		Ext:                   ext,
	}
}
