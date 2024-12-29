package openrtb

import "encoding/json"

// Request source details on post-auction decisioning (e.g., header bidding).
type Source struct {
	// Entity responsible for the final impression sale decision, where:
	//    0 = exchange;
	//    1 = upstream source.
	//
	// Recommended.
	FinalSaleDecision int `json:"fd" bson:"fd"`

	// Transaction ID that must be common across all participants in this bid request
	// (e.g., potentially multiple exchanges).
	//
	// Recommended.
	TransactionID string `json:"tid" bson:"tid"`

	// Payment ID chain string containing embedded syntax described in the TAG Payment
	// ID Protocol v1.0.
	//
	// Recommended.
	PaymentChain string `json:"pchain" bson:"pchain"`

	// Placeholder for exchange-specific extensions to OpenRTB.
	Ext json.RawMessage `json:"ext,omitempty" bson:"ext"`
}

func (s *Source) Copy() Source {

	var ext []byte
	if len(s.Ext) != 0 {
		ext = make([]byte, len(s.Ext))
		copy(ext, s.Ext)
	}

	return Source{
		FinalSaleDecision: s.FinalSaleDecision,
		TransactionID:     s.TransactionID,
		PaymentChain:      s.PaymentChain,
		Ext:               ext,
	}
}
