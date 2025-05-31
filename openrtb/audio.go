package openrtb

import (
	"encoding/json"
	"errors"
)

// Validation errors.
var (
	ErrInvalidAudioNoMIMEs = errors.New("audio has no MIMEs")
)

// Container for an audio impression.
// Audio object must be included directly in the impression object.
type Audio struct {
	// Content MIME types supported (e.g., “audio/mp4”).
	//
	// Required.
	MIMEs []string `json:"mimes" bson:"mimes"`

	// Minimum audio ad duration in seconds.
	//
	// Recommended.
	MinDuration int `json:"minduration,omitempty" bson:"minduration"`

	// Maximum audio ad duration in seconds.
	//
	// Recommended.
	MaxDuration int `json:"maxduration,omitempty" bson:"maxduration"`

	// Array of supported audio protocols.
	Protocols []Protocol `json:"protocols,omitempty" bson:"protocols"`

	// Indicates the start delay in seconds for pre-roll, mid-roll, or post-roll ad placements.
	//
	// Recommended.
	StartDelay StartDelay `json:"startdelay,omitempty" bson:"startdelay"`

	// If multiple ad impressions are offered in the same bid request, the sequence number
	// will allow for the coordinated delivery of multiple creatives.
	//
	// Default 1.
	Sequence int `json:"sequence,omitempty" bson:"sequence"`

	// Blocked creative attributes
	BlockedAttrs []CreativeAttribute `json:"battr,omitempty" bson:"battr"`

	// Maximum extended ad duration if extension is allowed.
	//
	// If blank or 0, extension is not allowed.
	//
	// If -1, extension is allowed, and there is no time limit imposed.
	//
	// If greater than 0, then the value represents the number of seconds of extended
	// play supported beyond the maxduration value.
	MaxExtended int `json:"maxextended,omitempty" bson:"maxextended"`

	// Minimum bit rate in Kbps.
	MinBitrate int `json:"minbitrate,omitempty" bson:"minbitrate"`

	// Maximum bit rate in Kbps.
	MaxBitrate int `json:"maxbitrate,omitempty" bson:"maxbitrate"`

	// Supported delivery methods (e.g., streaming, progressive). If none specified,
	// assume all are supported.
	Delivery []ContentDelivery `json:"delivery,omitempty" bson:"delivery"`

	// Array of Banner objects if companion ads are available.
	CompanionAds []Banner `json:"companionad,omitempty" bson:"companionad"`

	// List of supported API frameworks for this impression.
	//
	// If an API is not explicitly listed, it is assumed not to be supported.
	APIs []APIFramework `json:"api,omitempty" bson:"api"`

	// Supported DAAST companion ad types.
	//
	// Recommended if companion Banner objects are included via the companionad array.
	CompanionTypes []CompanionType `json:"companiontype,omitempty" bson:"companiontype"`

	// The maximum number of ads that can be played in an ad pod.
	MaxSequence int `json:"maxseq,omitempty" bson:"maxseq"`

	// Type of audio feed.
	Feed int `json:"feed,omitempty" bson:"feed"`

	// Indicates if the ad is stitched with audio content or delivered independently, where:
	//    0 = no;
	//    1 = yes.
	Stitched int `json:"stitched,omitempty" bson:"stitched"`

	// Volume normalization mode.
	VolumeNorm int `json:"nvol,omitempty" bson:"nvol"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

type jsonAudio Audio

// UnmarshalJSON custom unmarshaling.
func (a *Audio) UnmarshalJSON(data []byte) error {
	j := jsonAudio{Sequence: 1}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	*a = (Audio)(j)
	return nil
}

// Validate the Audio object.
func (a *Audio) Validate() error {
	if len(a.MIMEs) == 0 {
		return ErrInvalidAudioNoMIMEs
	}
	return nil
}

// GetSequence returns the sequence value.
//
// Deprecated.
func (a *Audio) GetSequence() int {
	return a.Sequence
}

func (a *Audio) Copy() Audio {

	var mimes []string
	if len(a.MIMEs) != 0 {
		mimes = make([]string, len(a.MIMEs))
		copy(mimes, a.MIMEs)
	}

	var protocols []Protocol
	if len(a.Protocols) != 0 {
		protocols = make([]Protocol, len(a.Protocols))
		copy(protocols, a.Protocols)
	}

	var battr []CreativeAttribute
	if len(a.BlockedAttrs) != 0 {
		battr = make([]CreativeAttribute, len(a.BlockedAttrs))
		copy(battr, a.BlockedAttrs)
	}

	var delivery []ContentDelivery
	if len(a.Delivery) != 0 {
		delivery = make([]ContentDelivery, len(a.Delivery))
		copy(delivery, a.Delivery)
	}

	var companionAds []Banner
	if len(a.CompanionAds) != 0 {
		companionAds = make([]Banner, len(a.CompanionAds))
		for i := range a.CompanionAds {
			companionAds[i] = a.CompanionAds[i].Copy()
		}
	}

	var apis []APIFramework
	if len(a.APIs) != 0 {
		apis = make([]APIFramework, len(a.APIs))
		copy(apis, a.APIs)
	}

	var companionTypes []CompanionType
	if len(a.CompanionTypes) != 0 {
		companionTypes = make([]CompanionType, len(a.CompanionTypes))
		copy(companionTypes, a.CompanionTypes)
	}

	var ext []byte
	if len(a.Ext) != 0 {
		ext = make([]byte, len(a.Ext))
		copy(ext, a.Ext)
	}

	return Audio{
		MIMEs:          mimes,
		MinDuration:    a.MinDuration,
		MaxDuration:    a.MaxDuration,
		Protocols:      protocols,
		StartDelay:     a.StartDelay,
		Sequence:       a.Sequence,
		BlockedAttrs:   battr,
		MaxExtended:    a.MaxExtended,
		MinBitrate:     a.MinBitrate,
		MaxBitrate:     a.MaxBitrate,
		Delivery:       delivery,
		CompanionAds:   companionAds,
		APIs:           apis,
		CompanionTypes: companionTypes,
		MaxSequence:    a.MaxSequence,
		Feed:           a.Feed,
		Stitched:       a.Stitched,
		VolumeNorm:     a.VolumeNorm,
		Ext:            ext,
	}
}
