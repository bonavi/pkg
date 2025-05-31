package native

import (
	"encoding/json"
	"errors"
)

// Validation error.
var ErrInvalidRequestNoAssets = errors.New("request has no assets")

// NativeRequest is a Native Markup request object.
type NativeRequest struct {
	// Version of the Native Markup version in use.
	//
	// Default: 1.2.
	Ver string `json:"ver,omitempty"`

	// The Layout ID of the Native Ad unit.
	//
	// Deprecated: since version 1.2.
	Layout Layout `json:"layout,omitempty"`

	// The Ad unit ID of the Native Ad unit.
	//
	// Deprecated: since version 1.2.
	AdUnit AdUnit `json:"adunit,omitempty"`

	// The context in which the ad appears.
	//
	// Recommended.
	ContextType ContextType `json:"context,omitempty"`

	// A more detailed context in which the ad appears.
	ContextSubType ContextSubType `json:"contextsubtype,omitempty"`

	// The design/format/layout of the ad unit being offered.
	//
	// Recommended.
	PlacementType PlacementType `json:"plcmttype,omitempty"`

	// The number of identical placements in this Layout.
	//
	// Default: 1.
	PlacementCount int `json:"plcmtcnt,omitempty"`

	// 0 for the first ad, 1 for the second ad, and so on.
	//
	// Default: 0.
	Sequence int `json:"seq,omitempty"`

	// An array of Asset Objects.
	//
	// Required.
	Assets []AssetRequest `json:"assets"`

	// Whether the supply source / impression supports returning an
	// assetsurl instead of an asset object.
	//
	// 0 or the absence of the field indicates no such support.
	//
	// Default: 0.
	AURLSupport int `json:"aurlsupport,omitempty"`

	// Whether the supply source / impression supports returning a
	// dco url instead of an asset object.
	//
	// 0 or the absence of the field indicates no such support.
	//
	// Beta feature.
	//
	// Default: 0.
	DURLSupport int `json:"durlsupport,omitempty"`

	// Specifies what type of event tracking is supported.
	EventTrackers []EventTrackerRequest `json:"eventtrackers,omitempty"`

	// Set to 1 when the Native Ad supports buyer-specific privacy
	// notice.
	//
	// Set to 0 (or field absent) when the Native Ad doesn’t support
	// custom privacy links or if support is unknown.
	//
	// Default: 0.
	Privacy int `json:"privacy,omitempty"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty"`
}

type jsonNativeRequest NativeRequest

// UnmarshalJSON custom unmarshaling.
func (r *NativeRequest) UnmarshalJSON(data []byte) error {
	j := jsonNativeRequest{Ver: "1.2", PlacementCount: 1} //nolint:exhaustruct
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	*r = (NativeRequest)(j)
	return nil
}

// Valdidate the Native Request object.
func (r *NativeRequest) Validate() error {
	if len(r.Assets) == 0 {
		return ErrInvalidRequestNoAssets
	}
	for i := range r.Assets {
		asset := r.Assets[i]
		if err := asset.Validate(); err != nil {
			return err
		}
	}
	return nil
}
