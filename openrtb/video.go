package openrtb

import (
	"encoding/json"
	"errors"
)

// Validation errors.
var (
	ErrInvalidVideoNoMIMEs     = errors.New("video has no mimes")
	ErrInvalidVideoNoProtocols = errors.New("video has no protocols")
)

// Details for a video impression.
type Video struct {
	// Content MIME types supported (e.g., “video/x-ms-wmv”, “video/mp4”).
	//
	// Required.
	MIMEs []string `json:"mimes" bson:"mimes"`

	// Minimum video ad duration in seconds.
	//
	// Recommended.
	MinDuration int `json:"minduration" bson:"minduration"`

	// Maximum video ad duration in seconds.
	//
	// Recommended.
	MaxDuration int `json:"maxduration" bson:"maxduration"`

	// Array of supported video protocols. At least one supported protocol must be specified
	// in either the protocol or protocols attribute.
	//
	// Recommended.
	Protocols []Protocol `json:"protocols" bson:"protocols"`

	// Supported video protocol. At least one supported protocol must be specified in either
	// the protocol or protocols attribute.
	//
	// Deprecated: deprecated in favor of protocols.
	Protocol Protocol `json:"protocol" bson:"protocol"`

	// Width of the video player in device independent pixels (DIPS).
	//
	// Recommended.
	Width int `json:"w" bson:"w"`

	// Height of the video player in device independent pixels (DIPS).
	//
	// Recommended.
	Height int `json:"h" bson:"h"`

	// Indicates the start delay in seconds for pre-roll, mid-roll, or post-roll
	// ad placements.
	//
	// Recommended.
	StartDelay StartDelay `json:"startdelay" bson:"startdelay"`

	// Placement type for the impression.
	Placement VideoPlacement `json:"placement" bson:"placement"`

	// Indicates if the impression must be linear, nonlinear, etc. If none specified,
	// assume all are allowed.
	//
	// Default 1.
	Linearity VideoLinearity `json:"linearity" bson:"linearity"`

	// Indicates if the player will allow the video to be skipped, where:
	//    0 = no;
	//    1 = yes.
	// If a bidder sends markup/creative that is itself skippable, the Bid object should
	// include the attr array with an element of 16 indicating skippable video.
	Skip CreativeAttribute `json:"skip" bson:"skip"`

	// Videos of total duration greater than this number of seconds can be skippable;
	// only applicable if the ad is skippable.
	//
	// Default 0.
	SkipMin int `json:"skipmin" bson:"skipmin"`

	// Number of seconds a video must play before skipping is enabled; only applicable
	// if the ad is skippable.
	//
	// Default 0.
	SkipAfter int `json:"skipafter" bson:"skipafter"`

	// If multiple ad impressions are offered in the same bid request, the sequence number
	// will allow for the coordinated delivery of multiple creatives.
	//
	// Default 1.
	Sequence int `json:"sequence" bson:"sequence"`

	// Blocked creative attributes.
	BlockedAttrs []CreativeAttribute `json:"battr" bson:"battr"`

	// Maximum extended ad duration if extension is allowed.
	//
	// If blank or 0, extension is not allowed.
	//
	// If -1, extension is allowed, and there is no time limit imposed.
	//
	// If greater than 0, then the value represents the number of seconds
	// of extended play supported beyond the maxduration value.
	MaxExtended int `json:"maxextended" bson:"maxextended"`

	// Minimum bit rate in Kbps.
	MinBitrate int `json:"minbitrate" bson:"minbitrate"`

	// Maximum bit rate in Kbps.
	MaxBitrate int `json:"maxbitrate" bson:"maxbitrate"`

	// Indicates if letter-boxing of 4:3 content into a 16:9 window is allowed,
	// where:
	//    0 = no;
	//    1 = yes.
	//
	// Default 1.
	BoxingAllowed int `json:"boxingallowed" bson:"boxingallowed"`

	// Playback methods that may be in use.
	//
	// If none are specified, any method may be used.
	//
	// Only one method is typically used in practice. As a result, this array may be
	// converted to an integer in a future version of the specification. It is strongly
	// advised to use only the first element of this array in preparation for this change.
	PlaybackMethods []VideoPlayback `json:"playbackmethod" bson:"playbackmethod"`

	// The event that causes playback to end.
	PlaybackEnd VideoCessation `json:"playbackend" bson:"playbackend"`

	// Supported delivery methods (e.g., streaming, progressive).
	//
	// If none specified, assume all are supported.
	Delivery []ContentDelivery `json:"delivery" bson:"delivery"`

	// Ad position on screen
	Position AdPosition `json:"pos" bson:"pos"`

	// Array of Banner objects if companion ads are available.
	CompanionAds []Banner `json:"companionad" bson:"companionad"`

	// List of supported API frameworks for this impression.
	//
	// If an API is not explicitly listed, it is assumed not to be supported.
	APIs []APIFramework `json:"api" bson:"api"`

	// Supported VAST companion ad types.
	//
	// Recommended if companion Banner objects are included via the companionad array.
	//
	// If one of these banners will be rendered as an end-card, this can be specified
	// using the vcm attribute with the particular banner.
	CompanionTypes []int `json:"companiontype" bson:"companiontype"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

// Validate the video object.
func (v *Video) Validate() error {
	if len(v.MIMEs) == 0 {
		return ErrInvalidVideoNoMIMEs
	} else if v.Protocol == 0 && len(v.Protocols) == 0 {
		return ErrInvalidVideoNoProtocols
	}
	return nil
}

type jsonVideo Video

// UnmarshalJSON custom unmarshaling.
func (v *Video) UnmarshalJSON(data []byte) error {
	j := jsonVideo{
		Linearity:     VideoLinearityLinear,
		Sequence:      1,
		BoxingAllowed: 1,
	}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	*v = (Video)(j)
	return nil
}

// GetBoxingAllowed returns the boxing-allowed indicator.
//
// Deprecated.
func (v *Video) GetBoxingAllowed() int {
	return v.BoxingAllowed
}

// GetLinearity returns the video linearity value.
//
// Deprecated.
func (v *Video) GetLinearity() VideoLinearity {
	return v.Linearity
}

// GetSequence returns the sequence value.
//
// Deprecated.
func (v *Video) GetSequence() int {
	return v.Sequence
}

func (v *Video) Copy() Video {

	var mimes []string
	if len(v.MIMEs) != 0 {
		mimes = make([]string, len(v.MIMEs))
		copy(mimes, v.MIMEs)
	}

	var protocols []Protocol
	if len(v.Protocols) != 0 {
		protocols = make([]Protocol, len(v.Protocols))
		copy(protocols, v.Protocols)
	}

	var battr []CreativeAttribute
	if len(v.BlockedAttrs) != 0 {
		battr = make([]CreativeAttribute, len(v.BlockedAttrs))
		copy(battr, v.BlockedAttrs)
	}

	var delivery []ContentDelivery
	if len(v.Delivery) != 0 {
		delivery = make([]ContentDelivery, len(v.Delivery))
		copy(delivery, v.Delivery)
	}

	var playbackMethods []VideoPlayback
	if len(v.PlaybackMethods) != 0 {
		playbackMethods = make([]VideoPlayback, len(v.PlaybackMethods))
		copy(playbackMethods, v.PlaybackMethods)
	}

	var companionAds []Banner
	if len(v.CompanionAds) != 0 {
		companionAds = make([]Banner, len(v.CompanionAds))
		for i := range v.CompanionAds {
			companionAds[i] = v.CompanionAds[i].Copy()
		}
	}

	var apis []APIFramework
	if len(v.APIs) != 0 {
		apis = make([]APIFramework, len(v.APIs))
		copy(apis, v.APIs)
	}

	var companionTypes []int
	if len(v.CompanionTypes) != 0 {
		companionTypes = make([]int, len(v.CompanionTypes))
		copy(companionTypes, v.CompanionTypes)
	}

	var ext []byte
	if len(v.Ext) != 0 {
		ext = make([]byte, len(v.Ext))
		copy(ext, v.Ext)
	}

	return Video{
		MIMEs:           mimes,
		MinDuration:     v.MinDuration,
		MaxDuration:     v.MaxDuration,
		Protocols:       protocols,
		Protocol:        v.Protocol,
		Width:           v.Width,
		Height:          v.Height,
		StartDelay:      v.StartDelay,
		Placement:       v.Placement,
		Linearity:       v.Linearity,
		Skip:            v.Skip,
		SkipMin:         v.SkipMin,
		SkipAfter:       v.SkipAfter,
		Sequence:        v.Sequence,
		BlockedAttrs:    battr,
		MaxExtended:     v.MaxExtended,
		MinBitrate:      v.MinBitrate,
		MaxBitrate:      v.MaxBitrate,
		BoxingAllowed:   v.BoxingAllowed,
		PlaybackMethods: playbackMethods,
		PlaybackEnd:     v.PlaybackEnd,
		Delivery:        delivery,
		Position:        v.Position,
		CompanionAds:    companionAds,
		APIs:            apis,
		CompanionTypes:  companionTypes,
		Ext:             ext,
	}
}
