package openrtb

import "encoding/json"

// Container for a native impression conforming to the Dynamic Native Ads API.
type Native struct {
	// Request payload complying with the Native Ad Specification.
	//
	// Required.
	Request string `json:"request" bson:"request"`

	// Version of the Dynamic Native Ads API to which request complies.
	//
	// Highly recommended for efficient parsing.
	Version string `json:"ver,omitempty" bson:"ver"`

	// List of supported API frameworks for this impression.
	//
	// If an API is not explicitly listed, it is assumed not to be supported.
	APIs []APIFramework `json:"api,omitempty" bson:"api"`

	// Blocked creative attributes.
	BlockedAttrs []CreativeAttribute `json:"battr,omitempty" bson:"battr"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (n *Native) Copy() Native {

	var apis []APIFramework
	if len(n.APIs) != 0 {
		apis = make([]APIFramework, len(n.APIs))
		copy(apis, n.APIs)
	}

	var battrs []CreativeAttribute
	if len(n.BlockedAttrs) != 0 {
		battrs = make([]CreativeAttribute, len(n.BlockedAttrs))
		copy(battrs, n.BlockedAttrs)
	}

	var ext []byte
	if len(n.Ext) != 0 {
		ext = make([]byte, len(n.Ext))
		copy(ext, n.Ext)
	}

	return Native{
		Request:      n.Request,
		Version:      n.Version,
		APIs:         apis,
		BlockedAttrs: battrs,
		Ext:          ext,
	}
}
