package openrtb

import (
	"encoding/json"
	"errors"
)

// Validation errors.
var (
	ErrInvalidBidResponseNoID = errors.New("bidresponse has no ID")
)

// BidResponse is the bid response wrapper object.
// ID and at least one "seatbid” object is required, which contains a bid on at least one impression.
// Other attributes are optional since an exchange may establish default values.
// No-Bids on all impressions should be indicated as a HTTP 204 response.
// For no-bids on specific impressions, the bidder should omit these from the bid response.
type BidResponse struct {
	// ID of the bid request to which this is a response.
	//
	// Required.
	ID string `json:"id" bson:"id"`

	// Array of seatbid objects; 1+ required if a bid is to be made.
	//
	// Required.
	SeatBids []SeatBid `json:"seatbid" bson:"seatbid"`

	// Bidder generated response ID to assist with logging/tracking.
	BidID string `json:"bidid,omitempty" bson:"bidid"`

	// Bid currency using ISO-4217 alpha codes.
	//
	// Default USD.
	Currency string `json:"cur,omitempty" bson:"cur"`

	// Optional feature to allow a bidder to set data in the exchange’s cookie.
	// The string must be in base85 cookie safe characters and be in any format.
	// Proper JSON encoding must be used to include “escaped” quotation marks.
	CustomData string `json:"customdata,omitempty" bson:"customdata"`

	// Reason for not bidding.
	NBR NBR `json:"nbr,omitempty" bson:"nbr"`

	// Placeholder for bidder-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

// Validate required attributes.
func (r *BidResponse) Validate() error {
	if r.ID == "" {
		return ErrInvalidBidResponseNoID
	} else if len(r.SeatBids) == 0 {
		return ErrInvalidSeatBidNoBids
	}

	for i := range r.SeatBids {
		seatBids := r.SeatBids[i]
		if err := seatBids.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r *BidResponse) Copy() BidResponse {

	var seatBids []SeatBid
	if len(r.SeatBids) != 0 {
		seatBids = make([]SeatBid, len(r.SeatBids))
		for i := range r.SeatBids {
			seatBids[i] = r.SeatBids[i].Copy()
		}
	}

	var ext []byte
	if len(r.Ext) != 0 {
		ext = make([]byte, len(r.Ext))
		copy(ext, r.Ext)
	}

	return BidResponse{
		ID:         r.ID,
		SeatBids:   seatBids,
		BidID:      r.BidID,
		Currency:   r.Currency,
		NBR:        r.NBR,
		CustomData: r.CustomData,
		Ext:        ext,
	}
}
