package openrtb

import (
	"errors"

	"pkg/decimal"
	"pkg/openrtb/ord"
	"pkg/pointer"
)

// Validation errors.
var (
	ErrInvalidBidNoID    = errors.New("bid has no ID")
	ErrInvalidBidNoImpID = errors.New("bid has no impression ID")
)

// Bid object contains bid information.
// ID, ImpID and Price are required; all other optional.
// If the bidder wins the impression, the exchange calls notice URL (nurl):
//
//	a) to inform the bidder of the win;
//	b) to convey certain information using substitution macros.
//
// Adomain can be used to check advertiser block list compliance.
// Cid can be used to block ads that were previously identified as inappropriate.
// Substitution macros may allow a bidder to use a static notice URL for all of its bids.
type Bid struct {
	// Bidder generated bid ID to assist with logging/tracking.
	//
	// Required.
	ID string `json:"id" bson:"id"`

	// ID of the Imp object in the related bid request.
	//
	// Required.
	ImpID string `json:"impid" bson:"impid"`

	// Bid price expressed as CPM although the actual transaction is for a unit
	// impression only. Note that while the type indicates float, integer math is highly
	// recommended when handling currencies (e.g., BigDecimal in Java).
	Price decimal.Decimal `json:"price" bson:"price"`

	// Win notice URL called by the exchange if the bid wins (not necessarily indicative
	// of a delivered, viewed, or billable ad); optional means of serving ad markup.
	// Substitution macros may be included in both the URL and optionally returned markup.
	NoticeURL string `json:"nurl,omitempty" bson:"nurl"`

	// Billing notice URL called by the exchange when a winning bid becomes billable based
	// on exchange-specific business policy (e.g., typically delivered, viewed, etc.).
	// Substitution macros may be included.
	BillingURL string `json:"burl,omitempty" bson:"burl"`

	// Loss notice URL called by the exchange when a bid is known to have been lost.
	// Substitution macros may be included. Exchange-specific policy may preclude support for
	// loss notices or the disclosure of winning clearing prices resulting in ${AUCTION_PRICE}
	// macros being removed (i.e., replaced with a zero-length string).
	LossURL string `json:"lurl,omitempty" bson:"lurl"`

	// Optional means of conveying ad markup in case the bid wins; supersedes the win notice
	// if markup is included in both. Substitution macros may be included.
	AdMarkup string `json:"adm,omitempty" bson:"adm"`

	// ID of a preloaded ad to be served if the bid wins.
	AdID string `json:"adid,omitempty" bson:"adid"`

	// Advertiser domain for block list checking (e.g., “ford.com”). This can be an array of
	// for the case of rotating creatives. Exchanges can mandate that only one domain is allowed.
	AdvDomains []string `json:"adomain,omitempty" bson:"adomain"`

	// A platform-specific application identifier intended to be unique to the app and independent
	// of the exchange. On Android, this should be a bundle or package name (e.g., com.foo.mygame).
	// On iOS, it is a numeric ID.
	Bundle string `json:"bundle,omitempty" bson:"bundle"`

	// URL without cache-busting to an image that is representative of the content of the campaign
	// for ad quality/safety checking.
	ImageURL string `json:"iurl,omitempty" bson:"iurl"`

	// Campaign ID to assist with ad quality checking; the collection of creatives for which iurl
	// should be representative.
	CampaingID string `json:"cid,omitempty" bson:"cid"`

	// Creative ID to assist with ad quality checking.
	CreativeID string `json:"crid,omitempty" bson:"crid"`

	// Tactic ID to enable buyers to label bids for reporting to the exchange the tactic through
	// which their bid was submitted. The specific usage and meaning of the tactic ID should be
	// communicated between buyer and exchanges a priori.
	Tactic string `json:"tactic,omitempty" bson:"tactic"`

	// IAB content categories of the creative.
	Categories []ContentCategory `json:"cat,omitempty" bson:"cat"`

	// Set of attributes describing the creative.
	Attrs []CreativeAttribute `json:"attr,omitempty" bson:"attr"`

	// API required by the markup if applicable.
	API APIFramework `json:"api,omitempty" bson:"api"`

	// Video response protocol of the markup if applicable.
	Protocol Protocol `json:"protocol" bson:"protocol"`

	// Creative media rating per IQG guidelines.
	MediaRating IQGRating `json:"qagmediarating,omitempty" bson:"qagmediarating"`

	// Language of the creative using ISO-639-1-alpha-2. The nonstandard code “xx” may also be
	// used if the creative has no linguistic content (e.g., a banner with just a company logo).
	Language string `json:"language,omitempty" bson:"language"`

	// Reference to the deal.id from the bid request if this bid pertains to a private
	// marketplace direct deal.
	DealID string `json:"dealid,omitempty" bson:"dealid"`

	// Width of the creative in device independent pixels (DIPS).
	Width int `json:"w,omitempty" bson:"w"`

	// Height of the creative in device independent pixels (DIPS).
	Height int `json:"h,omitempty" bson:"h"`

	// Relative width of the creative when expressing size as a ratio.
	//
	// Required for Flex Ads.
	WidthRatio int `json:"wratio,omitempty" bson:"wratio"`

	// Relative height of the creative when expressing size as a ratio.
	//
	// Required for Flex Ads.
	HeightRatio int `json:"hratio,omitempty" bson:"hratio"`

	// Advisory as to the number of seconds the bidder is willing to wait between the auction
	// and the actual impression.
	Exp int `json:"exp,omitempty" bson:"exp"`

	// Placeholder for bidder-specific extensions to OpenRTB.
	Ext *BidExt `json:"ext,omitempty" bson:"ext"`
}

type BidExt struct {
	Nroa *ord.Nroa `json:"nroa,omitempty" bson:"nroa"`
}

func (b *BidExt) copy() BidExt {

	var nroa *ord.Nroa
	if b.Nroa != nil {
		nroa = pointer.Pointer(b.Nroa.Copy())
	}

	return BidExt{
		Nroa: nroa,
	}
}

// Validate required attributes.
func (b *Bid) Validate() error {
	if b.ID == "" {
		return ErrInvalidBidNoID
	} else if b.ImpID == "" {
		return ErrInvalidBidNoImpID
	}
	return nil
}

func (b *Bid) Copy() Bid {

	var advDomains []string
	if len(b.AdvDomains) != 0 {
		advDomains = make([]string, len(b.AdvDomains))
		copy(advDomains, b.AdvDomains)
	}

	var categories []ContentCategory
	if len(b.Categories) != 0 {
		categories = make([]ContentCategory, len(b.Categories))
		copy(categories, b.Categories)
	}

	var attrs []CreativeAttribute
	if len(b.Attrs) != 0 {
		attrs = make([]CreativeAttribute, len(b.Attrs))
		copy(attrs, b.Attrs)
	}

	var ext *BidExt
	if b.Ext != nil {
		ext = pointer.Pointer(b.Ext.copy())
	}

	return Bid{
		ID:          b.ID,
		ImpID:       b.ImpID,
		Price:       b.Price,
		NoticeURL:   b.NoticeURL,
		BillingURL:  b.BillingURL,
		LossURL:     b.LossURL,
		AdMarkup:    b.AdMarkup,
		AdID:        b.AdID,
		AdvDomains:  advDomains,
		Bundle:      b.Bundle,
		ImageURL:    b.ImageURL,
		CampaingID:  b.CampaingID,
		CreativeID:  b.CreativeID,
		Tactic:      b.Tactic,
		Categories:  categories,
		Attrs:       attrs,
		API:         b.API,
		Protocol:    b.Protocol,
		MediaRating: b.MediaRating,
		Language:    b.Language,
		DealID:      b.DealID,
		Width:       b.Width,
		Height:      b.Height,
		WidthRatio:  b.WidthRatio,
		HeightRatio: b.HeightRatio,
		Exp:         b.Exp,
		Ext:         ext,
	}
}
