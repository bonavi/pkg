// Package vast implements IAB VAST 4.2 https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf
package vast

import (
	"encoding/xml"
	"fmt"
)

const (
	VASTVersion1dot0 = "1.0"
	VASTVersion2dot0 = "2.0"
	VASTVersion3dot0 = "3.0"
	VASTVersion4dot0 = "4.0"
)

const EmptyVAST = `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="3.0"></VAST>`

// VAST is the root <VAST> tag
type VAST struct {
	// The version of the VAST spec (should be either "2.0" or "3.0")
	Version string `xml:"version,attr" json:",omitempty"`
	// XML namespace. Most likely 'http://www.iab.com/VAST'
	XMLNS string `xml:"xmlns,attr,omitempty" json:"xmlns,omitempty"`
	// One or more Ad elements. Advertisers and video content publishers may
	// associate an <Ad> element with a line item video ad defined in contract
	// documentation, usually an insertion order. These line item ads typically
	// specify the creative to display, price, delivery schedule, targeting,
	// and so on.
	Ads []Ad `xml:"Ad,omitempty" json:"Ad,omitempty"`
	// Contains a URI to a tracking resource that the video player should request
	// upon receiving a “no ad” response
	Errors []CDATAString `xml:"Error,omitempty" json:",omitempty"`

	Mute bool `xml:"mute,attr,omitempty" json:",omitempty"`
}

// Validate validates VAST object
func (v *VAST) Validate() error {
	if v.Version == "" {
		return fmt.Errorf("version is required")
	}
	if len(v.Ads) == 0 {
		return fmt.Errorf("at least one Ad is required")
	}
	for _, ad := range v.Ads {
		if err := ad.Validate(); err != nil {
			return fmt.Errorf("invalid Ad: %w", err)
		}
	}
	return nil
}

// Ad represent an <Ad> child tag in a VAST document
//
// Each <Ad> contains a single <InLine> element or <Wrapper> element (but never both).
type Ad struct {
	InLine  *InLine  `xml:",omitempty" json:",omitempty"`
	Wrapper *Wrapper `xml:",omitempty" json:",omitempty"`
	// An ad server-defined identifier string for the ad
	ID string `xml:"id,attr,omitempty" json:",omitempty"`

	// An optional string that identifies the type of ad. This allows VAST to support audio ad scenarios.
	// Possible values – video, audio, hybrid.
	// Default value – video (assumed to be video if attribute is not present)
	AdType string `xml:"adType,attr,omitempty" json:",omitempty"`

	// Custom attr
	Type string `xml:"type,attr,omitempty" json:",omitempty"`

	// A number greater than zero (0) that identifies the sequence in which
	// an ad should play; all <Ad> elements with sequence values are part of
	// a pod and are intended to be played in sequence
	Sequence int `xml:"sequence,attr,omitempty" json:",omitempty"`

	// [Deprecated in VAST 4.1, along with apiFramework]
	// A Boolean value that identifies a conditional ad.
	// In the case of programmatic ad serving, a VPAID ad unit or other mechanism might be used to decide whether there is an ad that matches the placement.
	// When there is no match, an ad may not be served.
	// Use of the conditionalAd attribute enables publishers to avoid accepting these ads in placements where an ad must be served.
	// A value of true indicates that the ad is conditional and should be used in all cases where the InLine executable unit (such as VPAID) is not an ad but is instead a framework for finding an ad;
	// a value of false is the default value and indicates that an ad is available.
	ConditionalAd bool `xml:"conditionalAd,attr,omitempty" json:",omitempty"`
}

// Validate validates Ad object
func (a *Ad) Validate() error {
	if a.InLine == nil && a.Wrapper == nil {
		return fmt.Errorf("either InLine or Wrapper must be present")
	}
	if a.InLine != nil && a.Wrapper != nil {
		return fmt.Errorf("both InLine and Wrapper cannot be present")
	}
	if a.InLine != nil {
		if err := a.InLine.Validate(); err != nil {
			return fmt.Errorf("invalid InLine: %w", err)
		}
	}
	if a.Wrapper != nil {
		if err := a.Wrapper.Validate(); err != nil {
			return fmt.Errorf("invalid Wrapper: %w", err)
		}
	}
	return nil
}

type CDATAString struct {
	CDATA string `xml:",cdata" json:"Data"`
}

type PlainString struct {
	CDATA string `xml:",chardata" json:"Data"`
}

// <InLine> is a vast <InLine> ad element containing actual ad definition
//
// The last ad server in the ad supply chain serves an <InLine> element.
// Within the nested elements of an <InLine> element are all the files and
// URIs necessary to display the ad.
type InLine struct {
	// The name of the ad server that returned the ad
	AdSystem *AdSystem
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Errors []CDATAString `xml:"Error,omitempty" json:"Error,omitempty"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions *[]Extension `xml:"Extensions>Extension,omitempty" json:",omitempty"`
	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// Provides a value that represents a price that can be used by real-time bidding
	// (RTB) systems. VAST is not designed to handle RTB since other methods exist,
	// but this element is offered for custom solutions if needed.
	Pricing *Pricing `xml:",omitempty" json:",omitempty"`
	// Any ad server that returns a VAST containing an <InLine> ad must generate a pseudo- unique identifier
	// that is appropriate for all involved parties to track the lifecycle of that ad.
	// Example: ServerName-47ed3bac-1768-4b9a-9d0e-0b92422ab066
	AdServingId string `xml:",omitempty" json:",omitempty"`
	// The common name of the ad
	AdTitle PlainString
	// The name of the advertiser as defined by the ad serving party.
	// This element can be used to prevent displaying ads with advertiser
	// competitors. Ad serving parties and publishers should identify how
	// to interpret values provided within this element. As with any optional
	// elements, the video player is not required to support it.
	Advertiser *Advertiser `xml:",omitempty" json:",omitempty"`

	Category *[]Category `xml:",omitempty" json:",omitempty"`
	// The container for one or more <Creative> elements
	Creatives []Creative `xml:"Creatives>Creative,omitempty"`
	// A string value that provides a longer description of the ad.
	Description *CDATAString `xml:",omitempty" json:",omitempty"`
	// A URI to a survey vendor that could be the survey, a tracking pixel,
	// or anything to do with the survey. Multiple survey elements can be provided.
	// A type attribute is available to specify the MIME type being served.
	// For example, the attribute might be set to type=”text/javascript”.
	// Surveys can be dynamically inserted into the VAST response as long as
	// cross-domain issues are avoided.
	Survey *Survey `xml:",omitempty" json:",omitempty"`
	// The number of seconds in which the ad is valid for execution.
	//In cases where the ad is requested ahead of time, this timing indicates how many seconds after the request that the ad expires and cannot be played.
	//This element is useful for preventing an ad from playing after a timeout has occurred.
	Expires *int `xml:"Expires,omitempty" json:"Expires,omitempty"`
	// The ad server may provide URIs for tracking publisher-determined viewability,
	// for both the InLine ad and any Wrappers, using the <ViewableImpression> element.
	// Tracking URIs may be provided in three containers: <Viewable>, <NotViewable>, and <ViewUndetermined>.
	ViewableImpression *ViewableImpression `xml:",omitempty" json:",omitempty"`
	// The <AdVerifications> element contains one or more <Verification> elements,
	// which list the resources and metadata required to execute third-party measurement code in order to verify creative playback.
	// The <AdVerifications> element is used to contain one or more <Verification> elements,
	// which are used to initiate a controlled container where code can be executed for collecting data to verify ad playback details.
	AdVerifications *AdVerifications `xml:"AdVerifications,omitempty" json:",omitempty"`
}

// Impression is a URI that directs the video player to a tracking resource file that
// the video player should request when the first frame of the ad is displayed
type Impression struct {
	ID  string `xml:"id,attr,omitempty" json:",omitempty"`
	URI string `xml:",cdata" `
}

// Pricing provides a value that represents a price that can be used by real-time
// bidding (RTB) systems. VAST is not designed to handle RTB since other methods
// exist,  but this element is offered for custom solutions if needed.
type Pricing struct {
	// Identifies the pricing model as one of "cpm", "cpc", "cpe" or "cpv".
	Model string `xml:"model,attr"`
	// The 3 letter ISO-4217 currency symbol that identifies the currency of
	// the value provided
	Currency string `xml:"currency,attr"`
	// If the value provided is to be obfuscated/encoded, publishers and advertisers
	// must negotiate the appropriate mechanism to do so. When included as part of
	// a VAST Wrapper in a chain of Wrappers, only the value offered in the first
	// Wrapper need be considered.
	Value string `xml:",cdata"`
}

func (p *Pricing) Validate() error {
	if p.Model == "" {
		return fmt.Errorf("model is required")
	}
	if p.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if p.Value == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

// <Wrapper> element contains a URI reference to a vendor ad server (often called
// a third party ad server). The destination ad server either provides the ad
// files within a VAST <InLine> ad element or may provide a secondary Wrapper
// ad, pointing to yet another ad server. Eventually, the final ad server in
// the ad supply chain must contain all the necessary files needed to display
// the ad.
type Wrapper struct {
	// The name of the ad server that returned the ad
	AdSystem *AdSystem
	// URL of ad tag of downstream Secondary Ad Server
	VASTAdTagURI CDATAString
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Errors []CDATAString `xml:"Error,omitempty" json:"Error,omitempty"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions *[]Extension `xml:"Extensions>Extension,omitempty" json:",omitempty"`
	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// The container for one or more <Creative> elements
	Creatives *[]CreativeWrapper `xml:"Creatives>Creative,omitempty" json:",omitempty"`
	// Provides a value that represents a price that can be used by real-time bidding
	// (RTB) systems. VAST is not designed to handle RTB since other methods exist,
	// but this element is offered for custom solutions if needed.
	Pricing *Pricing `xml:",omitempty" json:",omitempty"`
	// The ad server may provide URIs for tracking publisher-determined viewability,
	// for both the InLine ad and any Wrappers, using the <ViewableImpression> element.
	// Tracking URIs may be provided in three containers: <Viewable>, <NotViewable>, and <ViewUndetermined>.
	ViewableImpression *ViewableImpression `xml:",omitempty" json:",omitempty"`
	// The <AdVerifications> element contains one or more <Verification> elements,
	// which list the resources and metadata required to execute third-party measurement code in order to verify creative playback.
	// The <AdVerifications> element is used to contain one or more <Verification> elements,
	// which are used to initiate a controlled container where code can be executed for collecting data to verify ad playback details.
	AdVerifications *AdVerifications `xml:"AdVerifications,omitempty" json:",omitempty"`

	FallbackOnNoAd           *bool `xml:"fallbackOnNoAd,attr,omitempty" json:",omitempty"`
	AllowMultipleAds         *bool `xml:"allowMultipleAds,attr,omitempty" json:",omitempty"`
	FollowAdditionalWrappers *bool `xml:"followAdditionalWrappers,attr,omitempty" json:",omitempty"`
}

// AdSystem contains information about the system that returned the ad
type AdSystem struct {
	Version string `xml:"version,attr,omitempty" json:"Version,omitempty"`
	Name    string `xml:",chardata" json:"Data"`
}

// Creative is a file that is part of a VAST ad.
type Creative struct {
	// An ad server-defined identifier for the creative
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// The preferred order in which multiple Creatives should be displayed
	Sequence int `xml:"sequence,attr,omitempty" json:",omitempty"`
	// Identifies the ad with which the creative is served
	AdID string `xml:"adId,attr,omitempty" json:",omitempty"`
	// The technology used for any included API
	APIFramework string `xml:"apiFramework,attr,omitempty" json:",omitempty"`
	// If present, provides a VAST 4.x universal ad id
	UniversalAdID *[]UniversalAdID `xml:"UniversalAdId,omitempty" json:",omitempty"`
	// If present, defines a linear creative
	Linear *Linear `xml:",omitempty" json:",omitempty"`
	// If defined, defins companions creatives
	CompanionAds *CompanionAds `xml:",omitempty" json:",omitempty"`
	// If defined, defines non linear creatives
	NonLinearAds *NonLinearAds `xml:",omitempty" json:",omitempty"`
	// When an API framework is needed to execute creative, a
	// <CreativeExtensions> element can be added under the <Creative>. This
	// extension can be used to load an executable creative with or without using
	// a media file.
	// A <CreativeExtension> element is nested under the <CreativeExtensions>
	// (plural) element so that any xml extensions are separated from VAST xml.
	// Additionally, any xml used in this extension should identify an xml name
	// space (xmlns) to avoid confusing any of the xml element names with those
	// of VAST.
	// The nested <CreativeExtension> includes an attribute for type, which
	// specifies the MIME type needed to execute the extension.
	CreativeExtensions *[]Extension `xml:"CreativeExtensions>CreativeExtension,omitempty" json:",omitempty"`
}

// <CompanionAds> contains companions creatives
type CompanionAds struct {
	// Provides information about which companion creative to display.
	// All means that the player must attempt to display all. Any means the player
	// must attempt to play at least one. None means all companions are optional
	Required   string      `xml:"required,attr,omitempty" json:",omitempty"`
	Companions []Companion `xml:"Companion,omitempty" json:",omitempty"`
}

// NonLinearAds contains non linear creatives
type NonLinearAds struct {
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty" json:",omitempty"`
	// Non linear creatives
	NonLinears []NonLinear `xml:"NonLinear,omitempty" json:",omitempty"`
}

// <CreativeWrapper> defines wrapped creative's parent trackers
type CreativeWrapper struct {
	// An ad server-defined identifier for the creative
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// The preferred order in which multiple Creatives should be displayed
	Sequence int `xml:"sequence,attr,omitempty" json:",omitempty"`
	// Identifies the ad with which the creative is served
	AdID string `xml:"adId,attr,omitempty" json:",omitempty"`
	// If present, defines a linear creative
	Linear *LinearWrapper `xml:",omitempty" json:",omitempty"`
	// If defined, defines companions creatives
	CompanionAds *CompanionAds `xml:"CompanionAds,omitempty" json:",omitempty"`
	// If defined, defines non linear creatives
	NonLinearAds *NonLinearAdsWrapper `xml:"NonLinearAds,omitempty" json:",omitempty"`
}

// <NonLinearAdsWrapper> contains non linear creatives in a wrapper
type NonLinearAdsWrapper struct {
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty" json:",omitempty"`
	// Non linear creatives
	NonLinears []NonLinearWrapper `xml:"NonLinear,omitempty" json:",omitempty"`
}

// <Linear> is the most common type of video advertisement trafficked in the
// industry is a “linear ad”, which is an ad that displays in the same area
// as the content but not at the same time as the content. In fact, the video
// player must interrupt the content before displaying a linear ad.
// <Linear> ads are often displayed right before the video content plays.
// This ad position is called a “pre-roll” position. For this reason, a linear
// ad is often called a “pre-roll.”
type Linear struct {
	// To specify that a Linear creative can be skipped, the ad server must
	// include the skipoffset attribute in the <Linear> element. The value
	// for skipoffset is a time value in the format HH:MM:SS or HH:MM:SS.mmm
	// or a percentage in the format n%. The .mmm value in the time offset
	// represents milliseconds and is optional. This skipoffset value
	// indicates when the skip control should be provided after the creative
	// begins playing.
	SkipOffset *Offset `xml:"skipoffset,attr,omitempty" json:",omitempty"`
	// Duration in standard time format, hh:mm:ss
	Duration       Duration        `xml:"Duration,omitempty" json:",omitempty"`
	Icons          *Icons          `json:",omitempty"`
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty" json:",omitempty"`
	AdParameters   *AdParameters   `xml:",omitempty" json:",omitempty"`
	VideoClicks    *VideoClicks    `xml:",omitempty" json:",omitempty"`
	MediaFiles     *MediaFiles     `xml:"MediaFiles" json:",omitempty"`
}

// <LinearWrapper> defines a wrapped linear creative
type LinearWrapper struct {
	Icons          *Icons          `json:",omitempty"`
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty" json:",omitempty"`
	VideoClicks    *VideoClicks    `xml:",omitempty" json:",omitempty"`
}

// <Companion> defines a companion ad
type Companion struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// Pixel dimensions of companion slot.
	Width int `xml:"width,attr,omitempty"`
	// Pixel dimensions of companion slot.
	Height int `xml:"height,attr,omitempty"`
	// Pixel dimensions of the companion asset.
	AssetWidth int `xml:"assetWidth,attr,omitempty"`
	// Pixel dimensions of the companion asset.
	AssetHeight int `xml:"assetHeight,attr,omitempty"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr,omitempty"`
	// Pixel dimensions of expanding companion ad when in expanded state.
	ExpandedHeight int `xml:"expandedHeight,attr,omitempty"`
	// The apiFramework defines the method to use for communication with the companion.
	APIFramework string `xml:"apiFramework,attr,omitempty" json:",omitempty"`
	// Used to match companion creative to publisher placement areas on the page.
	AdSlotID string `xml:"adSlotId,attr,omitempty" json:",omitempty"`
	// URL to open as destination page when user clicks on the the companion banner ad.
	CompanionClickThrough *CDATAString `xml:",omitempty" json:",omitempty"`
	// URLs to ping when user clicks on the the companion banner ad.
	CompanionClickTrackings []CompanionClickTracking `xml:"CompanionClickTracking,omitempty" json:",omitempty"`
	// Alt text to be displayed when companion is rendered in HTML environment.
	AltText string `xml:",omitempty" json:",omitempty"`
	// The creativeView should always be requested when present. For Companions
	// creativeView is the only supported event.
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty" json:",omitempty"`
	// Data to be passed into the companion ads. The apiFramework defines the method
	// to use for communication (e.g. “FlashVar”)
	AdParameters *AdParameters `xml:",omitempty" json:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty" json:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource *CDATAString `xml:",omitempty" json:",omitempty"`
	// HTML to display the companion element
	HTMLResource  *HTMLResource `xml:",omitempty" json:",omitempty"`
	Pxratio       string        `xml:"pxratio,attr,omitempty" json:",omitempty"`
	RenderingMode string        `xml:"renderingMode,attr,omitempty" json:",omitempty"`
}

// <NonLinear> defines a non linear ad
type NonLinear struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// Pixel dimensions of companion.
	Width int `xml:"width,attr,omitempty"`
	// Pixel dimensions of companion.
	Height int `xml:"height,attr,omitempty"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr,omitempty"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandedHeight int `xml:"expandedHeight,attr,omitempty"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty" json:",omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty" json:",omitempty"`
	// Suggested duration to display non-linear ad, typically for animation to complete.
	// Expressed in standard time format hh:mm:ss.
	MinSuggestedDuration *Duration `xml:"minSuggestedDuration,attr,omitempty" json:",omitempty"`
	// The apiFramework defines the method to use for communication with the nonlinear element.
	APIFramework string `xml:"apiFramework,attr,omitempty" json:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty" json:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource *CDATAString `xml:",omitempty" json:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty" json:",omitempty"`
	// Data to be passed into the video ad.
	AdParameters *AdParameters `xml:",omitempty" json:",omitempty"`
	// URL to open as destination page when user clicks on the non-linear ad unit.
	NonLinearClickThrough *CDATAString `xml:",omitempty" json:",omitempty"`
	// URLs to ping when user clicks on the the non-linear ad.
	NonLinearClickTrackings []NonLinearClickTracking `xml:"NonLinearClickTracking,omitempty" json:",omitempty"`
}

// <NonLinearWrapper> defines a non linear ad in a wrapper
type NonLinearWrapper struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// Pixel dimensions of companion.
	Width int `xml:"width,attr,omitempty"`
	// Pixel dimensions of companion.
	Height int `xml:"height,attr,omitempty"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandedWidth int `xml:"expandedWidth,attr,omitempty"`
	// Pixel dimensions of expanding nonlinear ad when in expanded state.
	ExpandedHeight int `xml:"expandedHeight,attr,omitempty"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty" json:",omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty" json:",omitempty"`
	// Suggested duration to display non-linear ad, typically for animation to complete.
	// Expressed in standard time format hh:mm:ss.
	MinSuggestedDuration *Duration `xml:"minSuggestedDuration,attr,omitempty" json:",omitempty"`
	// The apiFramework defines the method to use for communication with the nonlinear element.
	APIFramework string `xml:"apiFramework,attr,omitempty" json:",omitempty"`
	// The creativeView should always be requested when present.
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty" json:",omitempty"`
	// URLs to ping when user clicks on the the non-linear ad.
	NonLinearClickTracking []CDATAString `xml:",omitempty" json:",omitempty"`
}

type Icons struct {
	XMLName xml.Name `xml:"Icons,omitempty" json:",omitempty"`
	Icon    *[]Icon  `xml:"Icon,omitempty" json:",omitempty"`
}

// <Icon> represents advertising industry initiatives like AdChoices.
type Icon struct {
	// Identifies the industry initiative that the icon supports.
	Program string `xml:"program,attr,omitempty"`
	// Pixel dimensions of icon.
	Width int `xml:"width,attr,omitempty"`
	// Pixel dimensions of icon.
	Height int `xml:"height,attr,omitempty"`
	// The horizontal alignment location (in pixels) or a specific alignment.
	// Must match ([0-9]*|left|right)
	XPosition string `xml:"xPosition,attr,omitempty"`
	// The vertical alignment location (in pixels) or a specific alignment.
	// Must match ([0-9]*|top|bottom)
	YPosition string `xml:"yPosition,attr,omitempty"`
	// Start time at which the player should display the icon. Expressed in standard time format hh:mm:ss.
	Offset Offset `xml:"offset,attr,omitempty"`
	// duration for which the player must display the icon. Expressed in standard time format hh:mm:ss.
	Duration Duration `xml:"duration,attr,omitempty"`
	// The apiFramework defines the method to use for communication with the icon element
	APIFramework string `xml:"apiFramework,attr,omitempty" json:",omitempty"`
	// The pixel ratio for which the icon creative is intended.
	// The pixel ratio is the ratio of physical pixels on the device to the device-independent pixels.
	// An ad intended for display on a device with a pixel ratio that is twice that of a standard 1:1 pixel ratio would use the value "2."
	// Default value is "1."
	Pxratio string `xml:"pxratio,attr,omitempty" json:",omitempty"`
	// Alternative text for the image.
	// In an html5 image tag this should be the text for the alt attribute.
	// This should enable screen readers to properly read back a description of the icon for visually impaired users.
	AltText string `xml:"altText,attr,omitempty" json:",omitempty"`
	// Hover text for the image.
	// In an html5 image tag this should be the text for the title attribute.
	HoverText string `xml:"hoverText,attr,omitempty" json:",omitempty"`
	// The view tracking for icons is used to track when the icon creative is displayed.
	// The player uses the included URI to notify the icon server when the icon has been displayed.
	IconViewTracking []CDATAString `xml:",omitempty" json:",omitempty"`
	// URL to open as destination page when user clicks on the icon.
	IconClickThrough *CDATAString `xml:"IconClicks>IconClickThrough,omitempty" json:",omitempty"`
	// URLs to ping when user clicks on the the icon.
	IconClickTracking       []CDATAString            `xml:"IconClicks>IconClickTracking,omitempty" json:",omitempty"`
	IconClickFallbackImages *IconClickFallbackImages `xml:"IconClicks>IconClickFallbackImages,omitempty" json:",omitempty"`
	// URL to a static file, such as an image or SWF file
	StaticResource *StaticResource `xml:",omitempty" json:",omitempty"`
	// URL source for an IFrame to display the companion element
	IFrameResource *CDATAString `xml:",omitempty" json:",omitempty"`
	// HTML to display the companion element
	HTMLResource *HTMLResource `xml:",omitempty" json:",omitempty"`
}

const (
	TrackingEventStart            = "start"
	TrackingEventFirstQuartile    = "firstQuartile"
	TrackingEventMidpoint         = "midpoint"
	TrackingEventThirdQuartile    = "thirdQuartile"
	TrackingEventComplete         = "complete"
	TrackingEventMute             = "mute"
	TrackingEventUnmute           = "unmute"
	TrackingEventPause            = "pause"
	TrackingEventRewind           = "rewind"
	TrackingEventResume           = "resume"
	TrackingEventFullscreen       = "fullscreen"
	TrackingEventExitFullscreen   = "exitFullscreen"
	TrackingEventExpand           = "expand"
	TrackingEventCollapse         = "collapse"
	TrackingEventAcceptInvitation = "acceptInvitation"
	TrackingEventClose            = "close"
	TrackingEventSkip             = "skip"
	TrackingEventProgress         = "progress"
)

// <Tracking> defines an event tracking URL
type Tracking struct {
	// The name of the event to track for the element. The creativeView should
	// always be requested when present.
	//
	// Possible values are creativeView, start, firstQuartile, midpoint, thirdQuartile,
	// complete, mute, unmute, pause, rewind, resume, fullscreen, exitFullscreen, expand,
	// collapse, acceptInvitation, close, skip, progress.
	Event string `xml:"event,attr"`
	// The time during the video at which this url should be pinged. Must be present for
	// progress event. Must match (\d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?|1?\d?\d(\.?\d)*%)
	Offset *Offset `xml:"offset,attr,omitempty" json:",omitempty"`
	URI    string  `xml:",cdata"`

	// custom attr
	UA string `xml:"ua,attr,omitempty" json:",omitempty"`
}

type TrackingEvents struct {
	Tracking []Tracking `xml:"Tracking,omitempty"`
}

// <StaticResource> is the URL to a static file, such as an image or SWF file
type StaticResource struct {
	// Mime type of static resource
	CreativeType string `xml:"creativeType,attr,omitempty" json:",omitempty"`
	// URL to a static file, such as an image or SWF file
	URI string `xml:",cdata"`
}

// <HTMLResource> is a container for HTML data
type HTMLResource struct {
	// Specifies whether the HTML is XML-encoded
	XMLEncoded bool   `xml:"xmlEncoded,attr,omitempty" json:",omitempty"`
	HTML       string `xml:",cdata"`
}

// <AdParameters> defines arbitrary ad parameters
type AdParameters struct {
	// Specifies whether the parameters are XML-encoded
	XMLEncoded bool   `xml:"xmlEncoded,attr,omitempty" json:",omitempty"`
	Parameters string `xml:",cdata"`
}

// <VideoClicks> contains types of video clicks
type VideoClicks struct {
	ClickTrackings []VideoClick `xml:"ClickTracking,omitempty" json:",omitempty"`
	CustomClicks   []VideoClick `xml:"CustomClick,omitempty" json:",omitempty"`
	ClickThroughs  []VideoClick `xml:"ClickThrough,omitempty" json:",omitempty"`
}

// <VideoClick> defines a click URL for a linear creative
type VideoClick struct {
	ID  string `xml:"id,attr,omitempty" json:",omitempty"`
	URI string `xml:",cdata"`
}

// <MediaFile> defines a reference to a linear creative asset
type MediaFile struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// Method of delivery of ad (either "streaming" or "progressive")
	Delivery string `xml:"delivery,attr"`
	// MIME type. Popular MIME types include, but are not limited to
	// “video/x-ms-wmv” for Windows Media, and “video/x-flv” for Flash
	// Video. Image ads or interactive ads can be included in the
	// MediaFiles section with appropriate Mime types
	Type string `xml:"type,attr"`
	// Bitrate of encoded video in Kbps. If bitrate is supplied, MinBitrate
	// and MaxBitrate should not be supplied.
	Bitrate int `xml:"bitrate,attr,omitempty" json:",omitempty"`
	// Pixel dimensions of video.
	Width int `xml:"width,attr"`
	// Pixel dimensions of video.
	Height int `xml:"height,attr"`
	// Minimum bitrate of an adaptive stream in Kbps. If MinBitrate is supplied,
	// MaxBitrate must be supplied and Bitrate should not be supplied.
	MinBitrate int `xml:"minBitrate,attr,omitempty" json:",omitempty"`
	// Maximum bitrate of an adaptive stream in Kbps. If MaxBitrate is supplied,
	// MinBitrate must be supplied and Bitrate should not be supplied.
	MaxBitrate int `xml:"maxBitrate,attr,omitempty" json:",omitempty"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty" json:",omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty" json:",omitempty"`
	// The codec used to produce the media file.
	Codec string `xml:"codec,attr,omitempty" json:",omitempty"`
	// The APIFramework defines the method to use for communication if the MediaFile
	// is interactive. Suggested values for this element are “VPAID”, “FlashVars”
	// (for Flash/Flex), “initParams” (for Silverlight) and “GetVariables” (variables
	// placed in key/value pairs on the asset request).
	APIFramework string `xml:"apiFramework,attr,omitempty" json:",omitempty"`
	URI          string `xml:",cdata"`
	// Label
	Label string `xml:"label,attr,omitempty" json:",omitempty"`
	// Optional field that helps eliminate the need to calculate the size based on bitrate and duration.
	// Units - Bytes
	FileSize int `xml:"fileSize,attr,omitempty" json:",omitempty"`
	// Type of media file (2D / 3D / 360 / etc). Optional.
	// Default value = 2D
	MediaType string `xml:"mediaType,attr,omitempty" json:",omitempty"`
}

// <UniversalAdID> describes a VAST 4.x universal ad id.
type UniversalAdID struct {
	IDRegistry string `xml:"idRegistry,attr"`
	ID         string `xml:",chardata"`
}

// CompanionClickTracking element is used to track the click
type CompanionClickTracking struct {
	// An id provided by the ad server to track the click in reports.
	ID  string `xml:"id,attr,omitempty" json:",omitempty"`
	URI string `xml:",cdata"`
}

// NonLinearClickTracking element is used to track the click
type NonLinearClickTracking struct {
	// An id provided by the ad server to track the click in reports
	ID  string `xml:"id,attr,omitempty" json:",omitempty"`
	URI string `xml:",cdata"`
}

type Category struct {
	Authority string `xml:"authority,attr"`
	Category  string `xml:",chardata"`
}

type Survey struct {
	// A type attribute is available to specify the MIME type being served. For example,
	// the attribute might be set to type="text/javascript".
	// Surveys can be dynamically inserted into the VAST response as long as cross-domain issues are avoided.
	Type string `xml:"type,attr"`
	// A URI to any resource relating to an integrated survey.
	URI string `xml:",cdata"`
}

type MediaFiles struct {
	MediaFile               []MediaFile               `xml:",omitempty" json:",omitempty"`
	Mezzanine               []Mezzanine               `xml:",omitempty" json:",omitempty"`
	InteractiveCreativeFile []InteractiveCreativeFile `xml:",omitempty" json:",omitempty"`
	ClosedCaptionFiles      *[]ClosedCaptionFile      `xml:"ClosedCaptionFiles>ClosedCaptionFile,omitempty" json:",omitempty"`
}

// Validate validates InLine object
func (i *InLine) Validate() error {
	if i.AdSystem == nil {
		return fmt.Errorf("AdSystem is required")
	}
	if err := i.AdSystem.Validate(); err != nil {
		return fmt.Errorf("invalid AdSystem: %w", err)
	}
	if len(i.Impressions) == 0 {
		return fmt.Errorf("at least one Impression is required")
	}
	if i.AdTitle.CDATA == "" {
		return fmt.Errorf("AdTitle is required")
	}
	if len(i.Creatives) == 0 {
		return fmt.Errorf("at least one Creative is required")
	}
	for _, creative := range i.Creatives {
		if err := creative.Validate(); err != nil {
			return fmt.Errorf("invalid Creative: %w", err)
		}
	}
	if i.Pricing != nil {
		if err := i.Pricing.Validate(); err != nil {
			return fmt.Errorf("invalid Pricing: %w", err)
		}
	}
	return nil
}

// Validate validates Wrapper object
func (w *Wrapper) Validate() error {
	if w.AdSystem == nil {
		return fmt.Errorf("AdSystem is required")
	}
	if err := w.AdSystem.Validate(); err != nil {
		return fmt.Errorf("invalid AdSystem: %w", err)
	}
	if len(w.Impressions) == 0 {
		return fmt.Errorf("at least one Impression is required")
	}
	if w.VASTAdTagURI.CDATA == "" {
		return fmt.Errorf("VASTAdTagURI is required")
	}
	return nil
}

// Validate validates AdSystem object
func (a *AdSystem) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// Validate validates Creative object
func (c *Creative) Validate() error {
	if c.Linear == nil && c.CompanionAds == nil && c.NonLinearAds == nil {
		return fmt.Errorf("at least one of Linear, CompanionAds, or NonLinearAds must be present")
	}
	if c.Linear != nil {
		if err := c.Linear.Validate(); err != nil {
			return fmt.Errorf("invalid Linear: %w", err)
		}
	}
	if c.CompanionAds != nil {
		if err := c.CompanionAds.Validate(); err != nil {
			return fmt.Errorf("invalid CompanionAds: %w", err)
		}
	}
	if c.NonLinearAds != nil {
		if err := c.NonLinearAds.Validate(); err != nil {
			return fmt.Errorf("invalid NonLinearAds: %w", err)
		}
	}
	return nil
}

// Validate validates Linear object
func (l *Linear) Validate() error {
	if l.Duration.IsEmpty() {
		return fmt.Errorf("Duration is required")
	}
	if l.MediaFiles == nil || len(l.MediaFiles.MediaFile) == 0 {
		return fmt.Errorf("at least one MediaFile is required")
	}
	for _, mf := range l.MediaFiles.MediaFile {
		if err := mf.Validate(); err != nil {
			return fmt.Errorf("invalid MediaFile: %w", err)
		}
	}
	return nil
}

// Validate validates MediaFile object
func (m *MediaFile) Validate() error {
	if m.Delivery == "" {
		return fmt.Errorf("delivery is required")
	}
	if m.Type == "" {
		return fmt.Errorf("type is required")
	}
	if m.Width == 0 {
		return fmt.Errorf("width is required")
	}
	if m.Height == 0 {
		return fmt.Errorf("height is required")
	}
	if m.URI == "" {
		return fmt.Errorf("URI is required")
	}
	return nil
}

// Validate validates CompanionAds object
func (c *CompanionAds) Validate() error {
	if len(c.Companions) == 0 {
		return fmt.Errorf("at least one Companion is required")
	}
	for _, companion := range c.Companions {
		if err := companion.Validate(); err != nil {
			return fmt.Errorf("invalid Companion: %w", err)
		}
	}
	return nil
}

// Validate validates Companion object
func (c *Companion) Validate() error {
	if c.Width == 0 {
		return fmt.Errorf("width is required")
	}
	if c.Height == 0 {
		return fmt.Errorf("height is required")
	}
	if c.StaticResource == nil && c.IFrameResource == nil && c.HTMLResource == nil {
		return fmt.Errorf("at least one of StaticResource, IFrameResource, or HTMLResource must be present")
	}
	if c.StaticResource != nil {
		if err := c.StaticResource.Validate(); err != nil {
			return fmt.Errorf("invalid StaticResource: %w", err)
		}
	}
	return nil
}

// Validate validates StaticResource object
func (s *StaticResource) Validate() error {
	if s.CreativeType == "" {
		return fmt.Errorf("creativeType is required")
	}
	if s.URI == "" {
		return fmt.Errorf("URI is required")
	}
	return nil
}

// Validate validates NonLinearAds object
func (n *NonLinearAds) Validate() error {
	if len(n.NonLinears) == 0 {
		return fmt.Errorf("at least one NonLinear is required")
	}
	for _, nl := range n.NonLinears {
		if err := nl.Validate(); err != nil {
			return fmt.Errorf("invalid NonLinear: %w", err)
		}
	}
	return nil
}

// Validate validates NonLinear object
func (n *NonLinear) Validate() error {
	if n.Width == 0 {
		return fmt.Errorf("width is required")
	}
	if n.Height == 0 {
		return fmt.Errorf("height is required")
	}
	if n.StaticResource == nil && n.IFrameResource == nil && n.HTMLResource == nil {
		return fmt.Errorf("at least one of StaticResource, IFrameResource, or HTMLResource must be present")
	}
	if n.StaticResource != nil {
		if err := n.StaticResource.Validate(); err != nil {
			return fmt.Errorf("invalid StaticResource: %w", err)
		}
	}
	return nil
}

func (v *VAST) String() (string, error) {

	// Преобразуем структуру VAST в XML
	output, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}

	// Возвращаем строку XML
	return xml.Header + string(output), nil
}
